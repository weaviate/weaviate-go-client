package connection

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
	grpcbatch "github.com/weaviate/weaviate-go-client/v5/weaviate/grpc/batch"
	"github.com/weaviate/weaviate/entities/models"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

type GrpcClient struct {
	client  pb.WeaviateClient
	headers map[string]string
	timeout time.Duration
	batch   grpcbatch.Batch
}

func NewGrpcClient(host string, secured bool, headers map[string]string,
	gRPCVersionSupport *db.GRPCVersionSupport, timeout, startupTimeout time.Duration,
) (*GrpcClient, error) {
	client, err := createClient(host, secured, startupTimeout)
	if err != nil {
		return nil, fmt.Errorf("create grpc client: %w", err)
	}
	return &GrpcClient{client, headers, timeout, grpcbatch.New(gRPCVersionSupport)}, nil
}

func (c *GrpcClient) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchReply, error) {
	ctxWithTimeoutAndHeaders, cancel := c.ctxWithTimeoutWithHeaders(ctx)
	defer cancel()

	return c.client.Search(ctxWithTimeoutAndHeaders, req, c.getOptions()...)
}

func (c *GrpcClient) BatchObjects(ctx context.Context, objects []*models.Object,
	consistencyLevel string,
) ([]models.ObjectsGetResponse, error) {
	batchRequest, err := c.getBatchRequest(objects, consistencyLevel)
	if err != nil {
		return nil, err
	}
	reply, err := c.doBatchObjects(ctx, batchRequest)
	if err != nil {
		return nil, fmt.Errorf("batch objects: %w", err)
	}
	return c.batch.ParseReply(reply, objects), err
}

func (c *GrpcClient) doBatchObjects(ctx context.Context, batchRequest *pb.BatchObjectsRequest) (*pb.BatchObjectsReply, error) {
	ctxWithTimeoutAndHeaders, cancel := c.ctxWithTimeoutWithHeaders(ctx)
	defer cancel()

	return c.client.BatchObjects(ctxWithTimeoutAndHeaders, batchRequest, c.getOptions()...)
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

func (c *GrpcClient) ctxWithTimeoutWithHeaders(ctx context.Context) (context.Context, context.CancelFunc) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, c.timeout)
	if len(c.headers) > 0 {
		return metadata.NewOutgoingContext(ctxWithTimeout, metadata.New(c.headers)), cancel
	}
	return ctxWithTimeout, cancel
}

func (c *GrpcClient) getOptions() []grpc.CallOption {
	return []grpc.CallOption{}
}

func createClient(host string, secured bool, startupTimeout time.Duration) (pb.WeaviateClient, error) {
	var opts []grpc.DialOption
	if secured || strings.HasSuffix(host, ":443") {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.NewClient(getAddress(host, secured), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}
	if startupTimeout != 0 {
		// check if the gRPC connection is possible
		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), startupTimeout)
		client := grpc_health_v1.NewHealthClient(conn)
		_, err := client.Check(ctxWithTimeout, &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to connect to host: %s with secured set to: %v: %w", host, secured, err)
		}
		cancel()
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
