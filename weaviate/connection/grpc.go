package connection

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/db"
	grpcbatch "github.com/weaviate/weaviate-go-client/v4/weaviate/grpc/batch"
	"github.com/weaviate/weaviate/entities/models"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GrpcClient struct {
	client  pb.WeaviateClient
	headers map[string]string
	batch   grpcbatch.Batch
}

func NewGrpcClient(host string, secured bool, headers map[string]string,
	gRPCVersionSupport *db.GRPCVersionSupport,
) (*GrpcClient, error) {
	client, err := createClient(host, secured)
	if err != nil {
		return nil, fmt.Errorf("create grpc client: %w", err)
	}
	return &GrpcClient{client, headers, grpcbatch.New(gRPCVersionSupport)}, nil
}

func (c *GrpcClient) BatchObjects(ctx context.Context, objects []*models.Object,
	consistencyLevel string,
) ([]models.ObjectsGetResponse, error) {
	batchRequest, err := c.getBatchRequest(objects, consistencyLevel)
	if err != nil {
		return nil, err
	}
	reply, err := c.client.BatchObjects(c.ctxWithHeaders(ctx), batchRequest, c.getOptions()...)
	return c.batch.ParseReply(reply, objects), err
}

func (c *GrpcClient) getBatchRequest(objects []*models.Object, consistencyLevel string) (*pb.BatchObjectsRequest, error) {
	batchObjects, err := c.batch.GetBatchObjects(objects)
	if err != nil {
		return nil, err
	}
	return &pb.BatchObjectsRequest{
		Objects:          batchObjects,
		ConsistencyLevel: c.batch.GetConsistencyLevel(consistencyLevel),
	}, nil
}

func (c *GrpcClient) ctxWithHeaders(ctx context.Context) context.Context {
	if len(c.headers) > 0 {
		return metadata.NewOutgoingContext(ctx, metadata.New(c.headers))
	}
	return ctx
}

func (c *GrpcClient) getOptions() []grpc.CallOption {
	return []grpc.CallOption{}
}

func createClient(host string, secured bool) (pb.WeaviateClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	if secured || strings.HasSuffix(host, ":443") {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.Dial(getAddress(host, secured), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	return pb.NewWeaviateClient(conn), nil
}

func getAddress(host string, secured bool) string {
	if strings.Contains(host, ":") {
		return host
	}
	if secured {
		return fmt.Sprintf("%s:443", host)
	}
	return fmt.Sprintf("%s:80", host)
}
