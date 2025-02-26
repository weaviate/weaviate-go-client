package connection

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/grpc"
	"github.com/weaviate/weaviate/entities/models"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
	"go.nhat.io/grpcmock"
)

func TestTimeoutWeaviateRESTBatch(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/.well-known/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	})
	mux.HandleFunc("/v1/batch/objects", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second + time.Millisecond) // longer than timeout
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	cfg := weaviate.Config{Host: strings.TrimPrefix(s.URL, "http://"), Scheme: "http", Timeout: time.Second}
	client := weaviate.New(cfg)

	className := "TimeoutWeaviate"

	objects := []*models.Object{{Class: className}}
	_, batchErrSlice := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	assert.NotNil(t, batchErrSlice)
	assert.Contains(t, batchErrSlice.Error(), "context deadline exceeded")
}

func TestTimeoutWeaviateGRPCBatch(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/.well-known/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	ss := grpcmock.MockServer(
		grpcmock.RegisterService(pb.RegisterWeaviateServer),
		func(s *grpcmock.Server) {
			s.ExpectUnary("weaviate.v1.Weaviate/BatchObjects").After(time.Second + time.Millisecond).
				Return(&pb.BatchObjectsReply{})
		},
	)(t)

	cfg := weaviate.Config{
		Host:   strings.TrimPrefix(s.URL, "http://"),
		Scheme: "http",
		GrpcConfig: &grpc.Config{
			Host: strings.TrimPrefix(ss.Address(), "http://"),
		},
		Timeout: time.Second,
	}
	client := weaviate.New(cfg)

	className := "TimeoutWeaviateGRPC"

	objects := []*models.Object{{Class: className, Properties: map[string]interface{}{}}}
	_, batchErrSlice := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	assert.NotNil(t, batchErrSlice)
	assert.Contains(t, batchErrSlice.Error(), "context deadline exceeded")
}
