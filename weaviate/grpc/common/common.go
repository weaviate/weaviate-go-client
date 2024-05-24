package common

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/data/replication"
	"github.com/weaviate/weaviate/entities/models"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
	"github.com/weaviate/weaviate/usecases/byteops"
)

func GetConsistencyLevel(consistencyLevel string) *pb.ConsistencyLevel {
	switch consistencyLevel {
	case replication.ConsistencyLevel.ALL:
		return pb.ConsistencyLevel_CONSISTENCY_LEVEL_ALL.Enum()
	case replication.ConsistencyLevel.ONE:
		return pb.ConsistencyLevel_CONSISTENCY_LEVEL_ONE.Enum()
	case replication.ConsistencyLevel.QUORUM:
		return pb.ConsistencyLevel_CONSISTENCY_LEVEL_QUORUM.Enum()
	default:
		return nil
	}
}

func GetVector(name string, vector models.Vector) *pb.Vectors {
	switch v := vector.(type) {
	case []float32:
		return &pb.Vectors{
			Name:        name,
			VectorBytes: byteops.Fp32SliceToBytes(v),
			Type:        pb.Vectors_VECTOR_TYPE_SINGLE_FP32,
		}
	case [][]float32:
		return &pb.Vectors{
			Name:        name,
			VectorBytes: byteops.Fp32SliceOfSlicesToBytes(v),
			Type:        pb.Vectors_VECTOR_TYPE_MULTI_FP32,
		}
	default:
		return nil
	}
}
