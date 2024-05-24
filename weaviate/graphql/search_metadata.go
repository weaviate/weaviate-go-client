package graphql

import (
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type Metadata struct {
	uuid               bool
	creationTimeUnix   bool
	lastUpdateTimeUnix bool
	distance           bool
	certainty          bool
	score              bool
	explainScore       bool
	isConsistent       bool
	vector             bool
	vectors            []string
}

func NewMetadata() *Metadata {
	return &Metadata{}
}

func (m *Metadata) WithID() *Metadata {
	m.uuid = true
	return m
}

func (m *Metadata) WithCertainty() *Metadata {
	m.certainty = true
	return m
}

func (m *Metadata) WithDistance() *Metadata {
	m.distance = true
	return m
}

func (m *Metadata) WithScore() *Metadata {
	m.score = true
	return m
}

func (m *Metadata) WithExplainScore() *Metadata {
	m.explainScore = true
	return m
}

func (m *Metadata) WithCreationTimeUnix() *Metadata {
	m.creationTimeUnix = true
	return m
}

func (m *Metadata) WithLastUpdateTimeUnix() *Metadata {
	m.lastUpdateTimeUnix = true
	return m
}

func (m *Metadata) WithVector() *Metadata {
	m.vector = true
	return m
}

func (m *Metadata) WithVectors(vector ...string) *Metadata {
	m.vectors = vector
	return m
}

func (m *Metadata) WithIsConsistent() *Metadata {
	m.isConsistent = true
	return m
}

func (m *Metadata) togrpc() *pb.MetadataRequest {
	metadata := &pb.MetadataRequest{
		Uuid:               m.uuid,
		CreationTimeUnix:   m.creationTimeUnix,
		LastUpdateTimeUnix: m.lastUpdateTimeUnix,
		Distance:           m.distance,
		Certainty:          m.certainty,
		Score:              m.score,
		ExplainScore:       m.explainScore,
		IsConsistent:       m.isConsistent,
		Vector:             m.vector,
		Vectors:            m.vectors,
	}
	return metadata
}
