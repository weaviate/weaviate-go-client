package query_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

// As a general rule, client never returns nil maps, as they can
// cause a nil pointer dereference if not handled correctly.
// Tests that don't expect a field to be present in the response
// should use of the variables below to make the intent explicit.
var (
	noReferences = make(map[string][]types.Object[map[string]any])
	noProperties = make(map[string]any)
	noVectors    = make(types.Vectors)
)

func TestNearVector(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Songs",
		Tenant:           "john_doe",
		ConsistencyLevel: api.ConsistencyLevelQuorum,
	}

	for _, tt := range []struct {
		name  string
		nv    query.NearVector // Object to be inserted.
		stubs []testkit.Stub[api.SearchRequest, api.SearchResponse]
		want  *query.Result // Expected return value.
		err   testkit.Error
	}{
		{
			name: "request ok",
			nv: query.NearVector{
				Limit:      1,
				Offset:     2,
				AutoLimit:  3,
				After:      testkit.UUID,
				Similarity: query.Distance(.456),
				ReturnMetadata: query.ReturnMetadata{
					CreatedAt:    true,
					LastUpdateAt: true,
					Distance:     true,
					Certainty:    true,
					Score:        true,
					ExplainScore: true,
				},
				ReturnVectors:    []string{"title_vec", "lyrics_vec"},
				ReturnProperties: []string{"title", "duration_sec", "release_date"},
				ReturnNestedProperties: []query.NestedProperty{
					{
						Name: "label",
						ReturnNestedProperties: []query.NestedProperty{
							{Name: "name"},
							{Name: "logo"},
							{
								Name: "address",
								ReturnNestedProperties: []query.NestedProperty{
									{Name: "street"},
									{Name: "building_nr"},
								},
							},
						},
					},
				},
				ReturnReferences: []query.Reference{
					{
						PropertyName:     "hasAwards",
						TargetCollection: "GrammyAward",
						ReturnProperties: []string{"categories"},
					},
					{
						PropertyName:     "hasAwards",
						TargetCollection: "TonyAward",
						ReturnVectors:    []string{"recording_vec"},
					},
				},
			},
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{
					Request: &api.SearchRequest{
						RequestDefaults: rd,
						Limit:           1,
						Offset:          2,
						AutoLimit:       3,
						After:           testkit.UUID,
						ReturnMetadata: api.ReturnMetadata{
							CreatedAt:    true,
							LastUpdateAt: true,
							Distance:     true,
							Certainty:    true,
							Score:        true,
							ExplainScore: true,
						},
						ReturnVectors: []string{"title_vec", "lyrics_vec"},
						ReturnProperties: []api.ReturnProperty{
							{Name: "title"},
							{Name: "duration_sec"},
							{Name: "release_date"},
							{
								Name: "label",
								NestedProperties: []api.ReturnProperty{
									{Name: "name"},
									{Name: "logo"},
									{
										Name: "address",
										NestedProperties: []api.ReturnProperty{
											{Name: "street"},
											{Name: "building_nr"},
										},
									},
								},
							},
						},
						ReturnReferences: []api.ReturnReference{
							{
								PropertyName:     "hasAwards",
								TargetCollection: "GrammyAward",
								ReturnProperties: []api.ReturnProperty{
									{Name: "categories"},
								},
							},
							{
								PropertyName:     "hasAwards",
								TargetCollection: "TonyAward",
								ReturnVectors:    []string{"recording_vec"},
							},
						},
						NearVector: &api.NearVector{
							Distance: testkit.Ptr(.456),
						},
					},
					Response: api.SearchResponse{
						Took: 92 * time.Second,
						Results: []api.Object{
							{
								Collection: "Songs",
								Metadata: api.ObjectMetadata{
									UUID:          testkit.UUID,
									CreatedAt:     &testkit.Now,
									LastUpdatedAt: &testkit.Now,
									Distance:      testkit.Ptr[float32](.123),
									Certainty:     testkit.Ptr[float32](.4),
									Score:         testkit.Ptr[float32](10),
									ExplainScore:  testkit.Ptr("10/10"),
									Vectors: api.Vectors{
										"title_vec":  {Name: "title_vec", Single: []float32{1, 2, 3}},
										"lyrics_vec": {Name: "lyrics_vec", Multi: [][]float32{{1, 2, 3}, {1, 2, 3}}},
									},
								},
								Properties: map[string]any{
									"title":        "High Speed Dirt",
									"duration_sec": 252,
									"release_date": testkit.Now,
									"label": map[string]any{
										"name": "Capitol Records",
										"logo": "logo.png",
										"address": map[string]any{
											"street":      "Vine St",
											"building_nr": 1750,
										},
									},
								},
								References: map[string][]api.Object{
									"hasAwards": {
										{
											Collection: "GrammyAward",
											Properties: map[string]any{
												"categories": []string{"thrash_metal", "heavy_metal"},
											},
											References: make(map[string][]api.Object),
										},
										{
											Collection: "TonyAward",
											Properties: noProperties,
											References: make(map[string][]api.Object),
											Metadata: api.ObjectMetadata{
												UUID: testkit.UUID,
												Vectors: api.Vectors{
													"recording_vec": {
														Name:   "recording_vec",
														Single: []float32{4, 5, 6},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: &query.Result{
				Objects: []query.Object[map[string]any]{
					{
						Object: types.Object[map[string]any]{
							Collection: "Songs",
							UUID:       testkit.UUID,
							Properties: map[string]any{
								"title":        "High Speed Dirt",
								"duration_sec": 252,
								"release_date": testkit.Now,
								"label": map[string]any{
									"name": "Capitol Records",
									"logo": "logo.png",
									"address": map[string]any{
										"street":      "Vine St",
										"building_nr": 1750,
									},
								},
							},
							References: map[string][]types.Object[map[string]any]{
								"hasAwards": {
									{
										Collection: "GrammyAward",
										Properties: map[string]any{
											"categories": []string{"thrash_metal", "heavy_metal"},
										},
										Vectors:    noVectors,
										References: noReferences,
									},
									{
										Collection: "TonyAward",
										UUID:       testkit.UUID,
										Vectors: types.Vectors{
											"recording_vec": {
												Name:   "recording_vec",
												Single: []float32{4, 5, 6},
											},
										},
										Properties: noProperties,
										References: noReferences,
									},
								},
							},
							Vectors: types.Vectors{
								"title_vec":  {Name: "title_vec", Single: []float32{1, 2, 3}},
								"lyrics_vec": {Name: "lyrics_vec", Multi: [][]float32{{1, 2, 3}, {1, 2, 3}}},
							},
							CreatedAt:     &testkit.Now,
							LastUpdatedAt: &testkit.Now,
						},
						Metadata: query.Metadata{
							Distance:     testkit.Ptr[float32](.123),
							Certainty:    testkit.Ptr[float32](.4),
							Score:        testkit.Ptr[float32](10),
							ExplainScore: testkit.Ptr("10/10"),
						},
					},
				},
			},
		},
		{
			name: "request error",
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)

			c := query.NewClient(transport, rd)
			require.NotNil(t, c, "client")

			got, err := c.NearVector(t.Context(), tt.nv)
			tt.err.Require(t, err, "near vector query")
			require.Equal(t, tt.want, got, "query result")
		})
	}
}

func TestNearVector_GroupBy(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Songs",
		Tenant:           "john_doe",
		ConsistencyLevel: api.ConsistencyLevelQuorum,
	}

	for _, tt := range []struct {
		name    string
		nv      query.NearVector // Object to be inserted.
		groupBy query.GroupBy    // GroupBy clause.
		stubs   []testkit.Stub[api.SearchRequest, api.SearchResponse]
		want    *query.GroupByResult // Expected return value.
		err     testkit.Error
	}{
		{
			name: "request ok",
			nv: query.NearVector{
				Similarity: query.Certainty(.123),
			},
			groupBy: query.GroupBy{
				Property:       "album",
				ObjectLimit:    2,
				NumberOfGroups: 2,
			},
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{
					Request: &api.SearchRequest{
						RequestDefaults: rd,
						NearVector: &api.NearVector{
							Certainty: testkit.Ptr(.123),
						},
						GroupBy: &api.GroupBy{
							Property:       "album",
							ObjectLimit:    2,
							NumberOfGroups: 2,
						},
					},
					Response: api.SearchResponse{
						Took: 92 * time.Second,
						GroupByResults: []api.Group{
							{
								Name:        "Countdown To Extinction",
								MinDistance: .123,
								MaxDistance: .456,
								Size:        2,
								Objects: []api.GroupObject{
									{
										BelongsToGroup: "Countdown To Extinction",
										Object: api.Object{
											Properties: map[string]any{
												"title": "High Speed Dirt",
											},
										},
									},
									{
										BelongsToGroup: "Countdown To Extinction",
										Object: api.Object{
											Properties: map[string]any{
												"title": "Architechture Of Aggression",
											},
										},
									},
								},
							},
							{
								Name:        "Youthanasia",
								MinDistance: .321,
								MaxDistance: .654,
								Size:        1,
								Objects: []api.GroupObject{
									{
										BelongsToGroup: "Youthanasia",
										Object: api.Object{
											Properties: map[string]any{
												"title": "New World Order",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: &query.GroupByResult{
				Objects: []query.GroupObject[map[string]any]{
					{
						BelongsToGroup: "Countdown To Extinction",
						Object: query.Object[map[string]any]{
							Object: types.Object[map[string]any]{
								Properties: map[string]any{
									"title": "High Speed Dirt",
								},
								References: noReferences,
								Vectors:    noVectors,
							},
						},
					},
					{
						BelongsToGroup: "Countdown To Extinction",
						Object: query.Object[map[string]any]{
							Object: types.Object[map[string]any]{
								Properties: map[string]any{
									"title": "Architechture Of Aggression",
								},
								References: noReferences,
								Vectors:    noVectors,
							},
						},
					},
					{
						BelongsToGroup: "Youthanasia",
						Object: query.Object[map[string]any]{
							Object: types.Object[map[string]any]{
								Properties: map[string]any{
									"title": "New World Order",
								},
								References: noReferences,
								Vectors:    noVectors,
							},
						},
					},
				},
				Groups: map[string]query.Group[map[string]any]{
					"Countdown To Extinction": {
						Name:        "Countdown To Extinction",
						MinDistance: .123,
						MaxDistance: .456,
						Size:        2,
						Objects: []query.GroupObject[map[string]any]{
							{
								BelongsToGroup: "Countdown To Extinction",
								Object: query.Object[map[string]any]{
									Object: types.Object[map[string]any]{
										Properties: map[string]any{
											"title": "High Speed Dirt",
										},
										References: noReferences,
										Vectors:    noVectors,
									},
								},
							},
							{
								BelongsToGroup: "Countdown To Extinction",
								Object: query.Object[map[string]any]{
									Object: types.Object[map[string]any]{
										Properties: map[string]any{
											"title": "Architechture Of Aggression",
										},
										References: noReferences,
										Vectors:    noVectors,
									},
								},
							},
						},
					},
					"Youthanasia": {
						Name:        "Youthanasia",
						MinDistance: .321,
						MaxDistance: .654,
						Size:        1,
						Objects: []query.GroupObject[map[string]any]{
							{
								BelongsToGroup: "Youthanasia",
								Object: query.Object[map[string]any]{
									Object: types.Object[map[string]any]{
										Properties: map[string]any{
											"title": "New World Order",
										},
										References: noReferences,
										Vectors:    noVectors,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "request error",
			stubs: []testkit.Stub[api.SearchRequest, api.SearchResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)

			c := query.NewClient(transport, rd)
			require.NotNil(t, c, "client")

			got, err := c.NearVector.GroupBy(t.Context(), tt.nv, tt.groupBy)
			tt.err.Require(t, err, "near vector query")
			require.Equal(t, tt.want, got, "query result")
		})
	}
}

func TestSimilarity(t *testing.T) {
	t.Run("not set", func(t *testing.T) {
		var s query.Similarity
		assert.Nil(t, s.Distance(), "distance")
		assert.Nil(t, s.Certainty(), "certainty")
	})

	t.Run("distance", func(t *testing.T) {
		s := query.Distance(.1)
		assert.Equal(t, s.Distance(), testkit.Ptr(.1), "distance")
		assert.Nil(t, s.Certainty(), "certainty")
	})

	t.Run("certainty", func(t *testing.T) {
		s := query.Certainty(.1)
		assert.Equal(t, s.Certainty(), testkit.Ptr(.1), "certainty")
		assert.Nil(t, s.Distance(), "distance")
	})
}
