package common

import (
	"github.com/weaviate/weaviate-go-client/v4/weaviate/data/replication"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
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
