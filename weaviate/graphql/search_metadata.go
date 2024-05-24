package graphql

import (
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type Metadata struct {
	ID                 bool
	CreationTimeUnix   bool
	LastUpdateTimeUnix bool
	Distance           bool
	Certainty          bool
	Score              bool
	ExplainScore       bool
	IsConsistent       bool
	Vector             bool
	Vectors            []string
}

func (m *Metadata) togrpc() *pb.MetadataRequest {
	metadata := &pb.MetadataRequest{
		Uuid:               m.ID,
		CreationTimeUnix:   m.CreationTimeUnix,
		LastUpdateTimeUnix: m.LastUpdateTimeUnix,
		Distance:           m.Distance,
		Certainty:          m.Certainty,
		Score:              m.Score,
		ExplainScore:       m.ExplainScore,
		IsConsistent:       m.IsConsistent,
		Vector:             m.Vector,
		Vectors:            m.Vectors,
	}
	return metadata
}
