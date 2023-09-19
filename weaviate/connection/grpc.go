package connection

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/data/replication"
	"github.com/weaviate/weaviate/entities/models"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

type GrpcClient struct {
	client  pb.WeaviateClient
	headers map[string]string
}

func NewGrpcClient(scheme, host string, headers map[string]string) (*GrpcClient, error) {
	client, err := createClient(scheme, host)
	if err != nil {
		return nil, fmt.Errorf("create grpc client: %w", err)
	}
	return &GrpcClient{client, headers}, nil
}

func (c *GrpcClient) BatchObjects(ctx context.Context, objects []*models.Object,
	consistencyLevel string,
) ([]models.ObjectsGetResponse, error) {
	batchRequest, err := c.getBatchRequest(objects, consistencyLevel)
	if err != nil {
		return nil, err
	}
	reply, err := c.client.BatchObjects(ctx, batchRequest, c.getOptions()...)
	return c.parseReply(reply, objects), err
}

func (c *GrpcClient) getBatchRequest(objects []*models.Object, consistencyLevel string) (*pb.BatchObjectsRequest, error) {
	batchObjects, err := c.getBatchObjects(objects)
	if err != nil {
		return nil, err
	}
	return &pb.BatchObjectsRequest{
		Objects:          batchObjects,
		ConsistencyLevel: c.getConsistencyLevel(consistencyLevel),
	}, nil
}

func (c *GrpcClient) getOptions() []grpc.CallOption {
	var opts []grpc.CallOption
	if len(c.headers) > 0 {
		md := metadata.New(c.headers)
		opts = append(opts, grpc.Header(&md))
	}
	return opts
}

func (c *GrpcClient) getBatchObjects(objects []*models.Object) ([]*pb.BatchObject, error) {
	result := make([]*pb.BatchObject, len(objects))
	for i, obj := range objects {
		properties, err := c.getProperties(obj.Properties)
		if err != nil {
			return nil, err
		}
		batchObject := &pb.BatchObject{
			Uuid:       obj.ID.String(),
			Collection: obj.Class,
			Vector:     obj.Vector,
			Tenant:     obj.Tenant,
			Properties: properties,
		}
		result[i] = batchObject
	}
	return result, nil
}

func (c *GrpcClient) getProperties(properties models.PropertySchema) (*pb.BatchObject_Properties, error) {
	props, ok := properties.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("object properties type expected: map[string]interface{} got: %T", properties)
	}
	var result *pb.BatchObject_Properties
	if len(props) > 0 {
		nonRefProperties := map[string]interface{}{}
		numberArrayProperties := []*pb.NumberArrayProperties{}
		intArrayProperties := []*pb.IntArrayProperties{}
		textArrayProperties := []*pb.TextArrayProperties{}
		booleanArrayProperties := []*pb.BooleanArrayProperties{}
		for name, value := range props {
			switch v := value.(type) {
			case bool, int, int32, int64, uint, uint32, uint64, float32, float64, string:
				nonRefProperties[name] = v
			case []string:
				textArrayProperties = append(textArrayProperties, &pb.TextArrayProperties{
					PropName: name, Values: v,
				})
			case []bool:
				booleanArrayProperties = append(booleanArrayProperties, &pb.BooleanArrayProperties{
					PropName: name, Values: v,
				})
			case []int, []int32, []int64, []uint, []uint32, []uint64:
				var values []int64
				switch vv := v.(type) {
				case []int:
					values = toInt64Array[int](vv)
				case []int32:
					values = toInt64Array[int32](vv)
				case []int64:
					values = vv
				case []uint:
					values = toInt64Array[uint](vv)
				case []uint32:
					values = toInt64Array[uint32](vv)
				case []uint64:
					values = toInt64Array[uint64](vv)
				}
				intArrayProperties = append(intArrayProperties, &pb.IntArrayProperties{
					PropName: name, Values: values,
				})
			case []float32, []float64:
				var values []float64
				switch vv := v.(type) {
				case []float32:
					for i := range vv {
						values = append(values, float64(vv[i]))
					}
				case []float64:
					values = vv
				}
				numberArrayProperties = append(numberArrayProperties, &pb.NumberArrayProperties{
					PropName: name, Values: values,
				})
			}
			nonRefPropsStruct, err := structpb.NewStruct(nonRefProperties)
			if err != nil {
				return nil, fmt.Errorf("object properties: %w", err)
			}
			result = &pb.BatchObject_Properties{
				NonRefProperties:       nonRefPropsStruct,
				TextArrayProperties:    textArrayProperties,
				IntArrayProperties:     intArrayProperties,
				NumberArrayProperties:  numberArrayProperties,
				BooleanArrayProperties: booleanArrayProperties,
			}
		}
	}
	return result, nil
}

func (c *GrpcClient) getConsistencyLevel(consistencyLevel string) *pb.ConsistencyLevel {
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

func (c *GrpcClient) parseReply(reply *pb.BatchObjectsReply, objects []*models.Object) []models.ObjectsGetResponse {
	if reply != nil && len(reply.Errors) > 0 {
		result := make([]models.ObjectsGetResponse, len(reply.Errors))
		for i, res := range reply.Errors {
			var errors *models.ErrorResponse
			if res.Error != "" {
				errors = &models.ErrorResponse{
					Error: []*models.ErrorResponseErrorItems0{
						{Message: res.Error},
					},
				}
			}
			failed := models.ObjectsGetResponseAO2ResultStatusFAILED
			result[i] = models.ObjectsGetResponse{
				Result: &models.ObjectsGetResponseAO2Result{
					Errors: errors,
					Status: &failed,
				},
			}
		}
		return result
	}
	// all is OK
	success := models.ObjectsGetResponseAO2ResultStatusSUCCESS
	result := make([]models.ObjectsGetResponse, len(objects))
	for i := range objects {
		result[i] = models.ObjectsGetResponse{
			Result: &models.ObjectsGetResponseAO2Result{
				Status: &success,
			},
		}
	}
	return result
}

func toInt64Array[T int | int32 | int64 | uint | uint32 | uint64](arr []T) []int64 {
	result := make([]int64, len(arr))
	for i, val := range arr {
		result[i] = int64(val)
	}
	return result
}

func createClient(scheme, host string) (pb.WeaviateClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	if scheme == "https" || strings.HasSuffix(host, ":443") {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.Dial(getAddress(scheme, host), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	return pb.NewWeaviateClient(conn), nil
}

func getAddress(scheme, host string) string {
	if strings.Contains(host, ":") {
		return host
	}
	if scheme == "https" {
		return fmt.Sprintf("%s:443", host)
	}
	return fmt.Sprintf("%s:80", host)
}
