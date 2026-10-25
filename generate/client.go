package generate

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/query"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func NewClient(t internal.Transport, rd api.RequestDefaults) *Client {
	dev.AssertNotNil(t, "transport")
	return &Client{
		transport:  t,
		defaults:   rd,
		OverAll:    overAllFunc(t, rd),
		NearVector: nearVectorFunc(t, rd),
		NearObject: nearObjectFunc(t, rd),
		NearMedia:  nearMediaFunc(t, rd),
		NearText:   nearTextFunc(t, rd),
		Hybrid:     hybridFunc(t, rd),
		BM25:       bm25Func(t, rd),
	}
}

// Provider supplies configuration for runtime generative provider to the query.
type Provider interface {
	GenerativeProvider() *proto.GenerativeProvider
}

type Client struct {
	transport internal.Transport
	defaults  api.RequestDefaults

	OverAll    OverAllFunc
	NearVector NearVectorFunc
	NearObject NearObjectFunc
	NearMedia  NearMediaFunc
	NearText   NearTextFunc
	BM25       BM25Func
	Hybrid     HybridFunc
}

type Query interface {
	QueryRequest(api.RequestDefaults, *query.GroupBy) query.Request
}

// Request contains common generation parameters.
type Request[Q Query] struct {
	Query   Q
	Single  *SingleTask  // Per-object generative task.
	Grouped *GroupedTask // Per-group generative task.
	groupBy *query.GroupBy
}

type (
	SingleTask struct {
		Prompt   string
		Provider Provider
		Debug    bool
	}
	GroupedTask struct {
		Prompt            string
		IncludeProperties []string
		Provider          Provider
		Debug             bool
	}
)

type Result struct {
	Took       time.Duration
	Objects    []Object[map[string]any]
	Generative TaskOutput // Generative search summary.
}

type Object[T any] struct {
	query.Object[T]
	Generative TaskOutput // Per-object generative output.
}

type GroupByResult struct {
	Took       time.Duration
	Objects    []GroupObject[map[string]any]
	Groups     map[string]Group[map[string]any]
	Generative TaskOutput // Generative search summary.
}

type Group[T any] struct {
	Name                     string
	MinDistance, MaxDistance float32
	Size                     int64
	Objects                  []GroupObject[T]
	Generative               TaskOutput // Per-group generative output.
}

type GroupObject[T any] struct {
	Object[T]
	BelongsToGroup string
}

type TaskOutput struct {
	Text      string
	DebugInfo DebugInfo
}

type DebugInfo api.DebugInfo

func generate[Q Query](ctx context.Context, t internal.Transport, rd api.RequestDefaults, r *Request[Q]) (*Result, error) {
	q := r.Query.QueryRequest(rd, r.groupBy)
	req := q.Marshal()
	dev.AssertNotNil(req, "req")

	if task := r.Single; task != nil {
		req.Generate.Single = &api.SingleTask{
			Prompt: task.Prompt,
			Debug:  task.Debug,
		}
		if task.Provider != nil {
			req.Generate.Single.Provider = (*api.GenerativeProvider)(task.Provider.GenerativeProvider())
		}
	}

	if task := r.Grouped; task != nil {
		req.Generate.Grouped = &api.GroupedTask{
			Prompt:            task.Prompt,
			Debug:             task.Debug,
			IncludeProperties: task.IncludeProperties,
		}
		if task.Provider != nil {
			req.Generate.Grouped.Provider = (*api.GenerativeProvider)(task.Provider.GenerativeProvider())
		}
	}

	var resp api.SearchResponse
	if err := t.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Label, err)
	}

	// query was called from the GroupBy() method. This means we should put
	// GroupByResult in the context, as the first return value will be discarded.
	if q.GroupBy != nil {
		groups := internal.MakeMap[string, Group[map[string]any]](len(resp.GroupByResults))
		objects := make([]GroupObject[map[string]any], 0)
		for _, group := range resp.GroupByResults {

			objects = slices.Grow(objects, len(group.Objects))
			for _, obj := range group.Objects {
				objects = append(objects, GroupObject[map[string]any]{
					BelongsToGroup: group.Name,
					Object:         unmarshalObject(&obj.Object),
				})
			}

			// Create a view into the objects slice rather than allocating a separate one.
			from, to := len(objects)-len(group.Objects), len(objects)
			groups[group.Name] = Group[map[string]any]{
				Name:        group.Name,
				MinDistance: group.MinDistance,
				MaxDistance: group.MaxDistance,
				Size:        group.Size,
				Objects:     objects[from:to],
				Generative: TaskOutput{
					Text:      group.Generative.Text,
					DebugInfo: DebugInfo(group.Generative.DebugInfo),
				},
			}
		}
		internal.SetContextValue(ctx, groupByResultKey, &GroupByResult{
			Took:    resp.Took,
			Groups:  groups,
			Objects: objects,
		})
		return nil, nil
	}

	objects := slices.Grow([]Object[map[string]any](nil), len(resp.Results))
	for _, obj := range resp.Results {
		objects = append(objects, unmarshalObject(&obj))
	}
	return &Result{
		Took:    resp.Took,
		Objects: objects,
		Generative: TaskOutput{
			Text:      resp.Generative.Text,
			DebugInfo: DebugInfo(resp.Generative.DebugInfo),
		},
	}, nil
}

func unmarshalObject(o *api.Object) Object[map[string]any] {
	return Object[map[string]any]{
		Object: query.UnmarshalObject(o),
		Generative: TaskOutput{
			Text:      o.Generative.Text,
			DebugInfo: DebugInfo(o.Generative.DebugInfo),
		},
	}
}

// groupByResultKey is used to pass grouped query results to the GroupBy caller.
var groupByResultKey = internal.ContextKey{}

func generateGroupBy[In any](ctx context.Context, f func(context.Context, In) (*Result, error), in In) (*GroupByResult, error) {
	ctx = internal.ContextWithPlaceholder[GroupByResult](ctx, groupByResultKey)
	if _, err := f(ctx, in); err != nil {
		return nil, err
	}
	return internal.ValueFromContext[GroupByResult](ctx, groupByResultKey), nil
}
