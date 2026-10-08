package generate

import proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"

type Provider interface {
	GenerativeProvider() proto.GenerativeProvider
}
