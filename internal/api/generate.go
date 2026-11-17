package api

import (
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type GenerateRequest struct {
	Single  *SingleTask  // Per-object generative task.
	Grouped *GroupedTask // Per-group generative task.
}

type (
	GenerativeProvider proto.GenerativeProvider
	SingleTask         struct {
		Prompt   string
		Provider *GenerativeProvider
		Debug    bool
	}
	GroupedTask struct {
		Prompt            string
		IncludeProperties []string
		Provider          *GenerativeProvider
		Debug             bool
	}
)

type TaskOutput struct {
	Text      string
	DebugInfo DebugInfo
}

type ProviderMetadata struct {
	PromptTokens     int64
	CompletionTokens int64
	TotalTokens      int64
}

type DebugInfo struct {
	FullPrompt string
}

func marshalGenerateRequest(r GenerateRequest) *proto.GenerativeSearch {
	if r.Single == nil && r.Grouped == nil {
		return nil
	}
	var s proto.GenerativeSearch
	if r.Single != nil {
		s.Single = &proto.GenerativeSearch_Single{
			Prompt: r.Single.Prompt,
			Debug:  r.Single.Debug,
		}
		if r.Single.Provider != nil {
			s.Single.Queries = []*proto.GenerativeProvider{(*proto.GenerativeProvider)(r.Single.Provider)}
		}
	}
	if r.Grouped != nil {
		s.Grouped = &proto.GenerativeSearch_Grouped{
			Task:  r.Grouped.Prompt,
			Debug: r.Grouped.Debug,
			Properties: &proto.TextArray{
				Values: r.Grouped.IncludeProperties,
			},
		}
		if r.Grouped.Provider != nil {
			s.Grouped.Queries = []*proto.GenerativeProvider{(*proto.GenerativeProvider)(r.Grouped.Provider)}
		}
	}
	return &s
}
