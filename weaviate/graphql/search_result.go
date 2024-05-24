package graphql

import (
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
	"github.com/weaviate/weaviate/usecases/byteops"
)

type SearchResult struct {
	ID         string
	Collection string
	Properties map[string]any
	References []Reference
	Metadata   MetadataResult
	Vectors    map[string][]float32
	Vector     []float32
}

type Reference struct {
	Name                string
	ReferenceProperties []ReferencePropertiesResult
}

type ReferencePropertiesResult struct {
	Properties map[string]any
	Metadata   ReferenceMetadataResult
}

type MetadataResult struct {
	CreationTimeUnix           int64
	LastUpdateTimeUnix         int64
	Certainty, Distance, Score float32
	ExplainScore               string
	RerankScore                float64
	IsConsistent               bool
}

type ReferenceMetadataResult struct {
	MetadataResult
	ID      string
	Vectors map[string][]float32
	Vector  []float32
}

func toResults(results []*pb.SearchResult) []SearchResult {
	searchResults := make([]SearchResult, len(results))
	for i, r := range results {
		searchResults[i] = SearchResult{
			ID:         extractID(r.GetMetadata()),
			Collection: extractCollection(r.GetProperties()),
			Properties: extractProperties(r.GetProperties()),
			References: extractReferences(r.GetProperties()),
			Metadata:   extractMetadata(r.GetMetadata()),
			Vector:     extractVector(r.GetMetadata()),
			Vectors:    extractVectors(r.GetMetadata()),
		}
	}
	return searchResults
}

func extractProperties(p *pb.PropertiesResult) map[string]any {
	if p != nil {
		properties := make(map[string]any)
		if nonRefProps := p.GetNonRefProperties(); nonRefProps != nil {
			properties = nonRefProps.AsMap()
		} else if nonRefProps := p.GetNonRefProps(); nonRefProps != nil {
			for name, val := range nonRefProps.GetFields() {
				properties[name] = getValue(val)
			}
		}
		if props := p.GetTextArrayProperties(); len(props) > 0 {
			for i := range props {
				properties[props[i].GetPropName()] = props[i].GetValues()
			}
		}
		if props := p.GetBooleanArrayProperties(); len(props) > 0 {
			for i := range props {
				properties[props[i].GetPropName()] = props[i].GetValues()
			}
		}
		if props := p.GetIntArrayProperties(); len(props) > 0 {
			for i := range props {
				properties[props[i].GetPropName()] = props[i].GetValues()
			}
		}
		if props := p.GetNumberArrayProperties(); len(props) > 0 {
			for i := range props {
				properties[props[i].GetPropName()] = props[i].GetValues()
			}
		}
		return properties
	}
	return nil
}

func extractReferences(p *pb.PropertiesResult) []Reference {
	if p != nil {
		if refProps := p.GetRefProps(); len(refProps) > 0 {
			references := make([]Reference, len(refProps))
			for i := range refProps {
				references[i] = Reference{
					Name:                refProps[i].GetPropName(),
					ReferenceProperties: extractReferenceProperties(refProps[i].GetProperties()),
				}
			}
			return references
		}
	}
	return nil
}

func extractReferenceProperties(p []*pb.PropertiesResult) []ReferencePropertiesResult {
	if len(p) > 0 {
		properties := make([]ReferencePropertiesResult, len(p))
		for i := range p {
			properties[i] = ReferencePropertiesResult{
				Properties: extractProperties(p[i]),
				Metadata:   extractReferenceMetadata(p[i].GetMetadata()),
			}
		}
		return properties
	}
	return nil
}

func extractReferenceMetadata(m *pb.MetadataResult) ReferenceMetadataResult {
	var metadata ReferenceMetadataResult
	if m != nil {
		metadata = ReferenceMetadataResult{
			MetadataResult: extractMetadata(m),
			ID:             extractID(m),
			Vector:         extractVector(m),
			Vectors:        extractVectors(m),
		}
	}
	return metadata
}

func extractCollection(p *pb.PropertiesResult) string {
	if p != nil {
		return p.GetTargetCollection()
	}
	return ""
}

func extractID(m *pb.MetadataResult) string {
	if m != nil {
		return m.GetId()
	}
	return ""
}

func extractMetadata(m *pb.MetadataResult) MetadataResult {
	var metadata MetadataResult
	if m != nil {
		if m.GetCreationTimeUnixPresent() {
			metadata.CreationTimeUnix = m.GetCreationTimeUnix()
		}
		if m.GetLastUpdateTimeUnixPresent() {
			metadata.LastUpdateTimeUnix = m.GetLastUpdateTimeUnix()
		}
		if m.GetCertaintyPresent() {
			metadata.Certainty = m.GetCertainty()
		}
		if m.GetDistancePresent() {
			metadata.Distance = m.GetDistance()
		}
		if m.GetScorePresent() {
			metadata.Score = m.GetScore()
		}
		if m.GetExplainScorePresent() {
			metadata.ExplainScore = m.GetExplainScore()
		}
		if m.GetRerankScorePresent() {
			metadata.RerankScore = m.GetRerankScore()
		}
	}
	return metadata
}

func extractVector(m *pb.MetadataResult) []float32 {
	if m != nil && len(m.GetVectorBytes()) > 0 {
		return byteops.Float32FromByteVector(m.GetVectorBytes())
	}
	return nil
}

func extractVectors(m *pb.MetadataResult) map[string][]float32 {
	if m != nil && len(m.GetVectors()) > 0 {
		vectors := make(map[string][]float32, len(m.GetVectors()))
		for _, v := range m.GetVectors() {
			vectors[v.GetName()] = byteops.Float32FromByteVector(v.GetVectorBytes())
		}
		return vectors
	}
	return nil
}

func getValue(val *pb.Value) any {
	switch val.GetKind().(type) {
	case *pb.Value_TextValue:
		return val.GetTextValue()
	case *pb.Value_NumberValue:
		return val.GetNumberValue()
	case *pb.Value_BlobValue:
		return val.GetBlobValue()
	case *pb.Value_BoolValue:
		return val.GetBoolValue()
	case *pb.Value_DateValue:
		return val.GetDateValue()
	case *pb.Value_IntValue:
		return val.GetIntValue()
	case *pb.Value_UuidValue:
		return val.GetUuidValue()
	case *pb.Value_GeoValue:
		return val.GetGeoValue()
	case *pb.Value_ListValue:
		return val.GetListValue()
	case *pb.Value_ObjectValue:
		return val.GetObjectValue()
	case *pb.Value_PhoneValue:
		return val.GetPhoneValue()
	case *pb.Value_NullValue:
		return val.GetNullValue()
	default:
		return nil
	}
}
