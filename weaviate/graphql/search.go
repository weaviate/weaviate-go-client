package graphql

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/grpc/common"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type Search struct {
	grpcClient *connection.GrpcClient

	collection string
	tenant     string

	limit            uint32
	offset           uint32
	autocut          uint32
	after            string
	consistencyLevel string

	withNearText    *NearTextArgumentBuilder
	withNearVector  *NearVectorArgumentBuilder
	withNearObject  *NearObjectArgumentBuilder
	withNearImage   *NearImageArgumentBuilder
	withNearAudio   *NearAudioArgumentBuilder
	withNearVideo   *NearVideoArgumentBuilder
	withNearDepth   *NearDepthArgumentBuilder
	withNearThermal *NearThermalArgumentBuilder
	withNearImu     *NearImuArgumentBuilder
	withHybrid      *HybridArgumentBuilder
	withBM25        *BM25ArgumentBuilder

	withSortBy *SortBuilder

	withProperties []string
	withReferences []*Reference
	withMetadata   *Metadata
}

func NewSearch(grpcClient *connection.GrpcClient) *Search {
	return &Search{grpcClient: grpcClient}
}

func (s *Search) WithCollection(collection string) *Search {
	s.collection = collection
	return s
}

func (s *Search) WithTenant(tenant string) *Search {
	s.tenant = tenant
	return s
}

func (s *Search) WithLimit(limit int) *Search {
	s.limit = uint32(limit)
	return s
}

func (s *Search) WithOffset(offset int) *Search {
	s.offset = uint32(offset)
	return s
}

func (s *Search) WithAfter(after string) *Search {
	s.after = after
	return s
}

func (s *Search) WithAutocut(autocut int) *Search {
	s.autocut = uint32(autocut)
	return s
}

func (s *Search) WithConsistencyLevel(consistencyLevel string) *Search {
	s.consistencyLevel = consistencyLevel
	return s
}

func (s *Search) WithNearText(nearText *NearTextArgumentBuilder) *Search {
	s.withNearText = nearText
	return s
}

func (s *Search) WithNearVector(nearVector *NearVectorArgumentBuilder) *Search {
	s.withNearVector = nearVector
	return s
}

func (s *Search) WithNearObject(nearObject *NearObjectArgumentBuilder) *Search {
	s.withNearObject = nearObject
	return s
}

func (s *Search) WithNearImage(nearImage *NearImageArgumentBuilder) *Search {
	s.withNearImage = nearImage
	return s
}

func (s *Search) WithNearAudio(nearAudio *NearAudioArgumentBuilder) *Search {
	s.withNearAudio = nearAudio
	return s
}

func (s *Search) WithNearVideo(nearVideo *NearVideoArgumentBuilder) *Search {
	s.withNearVideo = nearVideo
	return s
}

func (s *Search) WithNearDepth(nearDepth *NearDepthArgumentBuilder) *Search {
	s.withNearDepth = nearDepth
	return s
}

func (s *Search) WithNearImu(nearImu *NearImuArgumentBuilder) *Search {
	s.withNearImu = nearImu
	return s
}

func (s *Search) WithNearThermal(nearThermal *NearThermalArgumentBuilder) *Search {
	s.withNearThermal = nearThermal
	return s
}

func (s *Search) WithHybrid(hybrid *HybridArgumentBuilder) *Search {
	s.withHybrid = hybrid
	return s
}

func (s *Search) WithBM25(bm25 *BM25ArgumentBuilder) *Search {
	s.withBM25 = bm25
	return s
}

func (s *Search) WithSort(sort ...Sort) *Search {
	if len(sort) > 0 {
		s.withSortBy = &SortBuilder{sort}
	}
	return s
}

func (s *Search) WithProperties(properties ...string) *Search {
	s.withProperties = properties
	return s
}

func (s *Search) WithReferences(references ...*Reference) *Search {
	s.withReferences = references
	return s
}

func (s *Search) WithMetadata(metadata *Metadata) *Search {
	s.withMetadata = metadata
	return s
}

func (s *Search) togrpc() *pb.SearchRequest {
	req := &pb.SearchRequest{
		Collection:       s.collection,
		Tenant:           s.tenant,
		Limit:            s.limit,
		Offset:           s.offset,
		Autocut:          s.autocut,
		After:            s.after,
		ConsistencyLevel: common.GetConsistencyLevel(s.consistencyLevel),
	}
	if s.withNearText != nil {
		req.NearText = s.withNearText.togrpc()
	}
	if s.withNearVector != nil {
		req.NearVector = s.withNearVector.togrpc()
	}
	if s.withNearObject != nil {
		req.NearObject = s.withNearObject.togrpc()
	}
	if s.withNearImage != nil {
		req.NearImage = s.withNearImage.togrpc()
	}
	if s.withNearAudio != nil {
		req.NearAudio = s.withNearAudio.togrpc()
	}
	if s.withNearVideo != nil {
		req.NearVideo = s.withNearVideo.togrpc()
	}
	if s.withNearDepth != nil {
		req.NearDepth = s.withNearDepth.togrpc()
	}
	if s.withNearImu != nil {
		req.NearImu = s.withNearImu.togrpc()
	}
	if s.withNearThermal != nil {
		req.NearThermal = s.withNearThermal.togrpc()
	}
	if s.withHybrid != nil {
		req.HybridSearch = s.withHybrid.togrpc()
	}
	if s.withBM25 != nil {
		req.Bm25Search = s.withBM25.togrpc()
	}
	if s.withSortBy != nil {
		req.SortBy = s.withSortBy.togrpc()
	}
	withProps := &Properties{}
	if len(s.withProperties) > 0 {
		withProps.WithProperties(s.withProperties...)
	}
	if len(s.withReferences) > 0 {
		withProps.WithReferences(s.withReferences...)
	}
	req.Properties = withProps.togrpc()
	if s.withMetadata != nil {
		req.Metadata = s.withMetadata.togrpc()
	} else {
		// by default always return ID
		req.Metadata = &pb.MetadataRequest{Uuid: true}
	}
	req.Uses_123Api = true
	req.Uses_125Api = true
	req.Uses_127Api = true
	return req
}

func (s *Search) Do(ctx context.Context) ([]SearchResult, error) {
	if s.grpcClient != nil {
		reply, err := s.grpcClient.Search(ctx, s.togrpc())
		if err != nil {
			return nil, err
		}
		return toResults(reply.Results), nil
	}
	return nil, fmt.Errorf("please provide gRPC config to the client in order to use search functionality")
}
