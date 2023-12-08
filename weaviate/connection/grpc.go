package connection

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/data/replication"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/db"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema/crossref"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
	"github.com/weaviate/weaviate/usecases/byteops"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

type GrpcClient struct {
	client             pb.WeaviateClient
	headers            map[string]string
	gRPCVersionSupport *db.GRPCVersionSupport
}

func NewGrpcClient(host string, secured bool, headers map[string]string,
	gRPCVersionSupport *db.GRPCVersionSupport,
) (*GrpcClient, error) {
	client, err := createClient(host, secured)
	if err != nil {
		return nil, fmt.Errorf("create grpc client: %w", err)
	}
	return &GrpcClient{client, headers, gRPCVersionSupport}, nil
}

func (c *GrpcClient) BatchObjects(ctx context.Context, objects []*models.Object,
	consistencyLevel string,
) ([]models.ObjectsGetResponse, error) {
	batchRequest, err := c.getBatchRequest(objects, consistencyLevel)
	if err != nil {
		return nil, err
	}
	reply, err := c.client.BatchObjects(c.ctxWithHeaders(ctx), batchRequest, c.getOptions()...)
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

func (c *GrpcClient) ctxWithHeaders(ctx context.Context) context.Context {
	if len(c.headers) > 0 {
		return metadata.NewOutgoingContext(ctx, metadata.New(c.headers))
	}
	return ctx
}

func (c *GrpcClient) getOptions() []grpc.CallOption {
	return []grpc.CallOption{}
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
		if obj.Vector != nil {
			if c.gRPCVersionSupport.SupportsVectorBytesField() {
				batchObject.VectorBytes = byteops.Float32ToByteVector(obj.Vector)
			} else {
				batchObject.Vector = obj.Vector
			}
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
		nonRefPropsStruct, numberArrayProperties, intArrayProperties, textArrayProperties,
			booleanArrayProperties, objectProperties, objectArrayProperties,
			singleTargetRefProps, multiTargetRefProps, err := c.extractProperties(props, true)
		if err != nil {
			return nil, err
		}
		result = &pb.BatchObject_Properties{
			NonRefProperties:       nonRefPropsStruct,
			TextArrayProperties:    textArrayProperties,
			IntArrayProperties:     intArrayProperties,
			NumberArrayProperties:  numberArrayProperties,
			BooleanArrayProperties: booleanArrayProperties,
			ObjectProperties:       objectProperties,
			ObjectArrayProperties:  objectArrayProperties,
			SingleTargetRefProps:   singleTargetRefProps,
			MultiTargetRefProps:    multiTargetRefProps,
		}
	}
	return result, nil
}

func (c *GrpcClient) extractProperties(properties map[string]interface{}, rootLevel bool) (nonRefProperties *structpb.Struct,
	numberArrayProperties []*pb.NumberArrayProperties,
	intArrayProperties []*pb.IntArrayProperties,
	textArrayProperties []*pb.TextArrayProperties,
	booleanArrayProperties []*pb.BooleanArrayProperties,
	objectProperties []*pb.ObjectProperties,
	objectArrayProperties []*pb.ObjectArrayProperties,
	singleTargetRefProps []*pb.BatchObject_SingleTargetRefProps,
	multiTargetRefProps []*pb.BatchObject_MultiTargetRefProps,
	err error,
) {
	nonRefPropertiesMap := map[string]interface{}{}
	numberArrayProperties = []*pb.NumberArrayProperties{}
	intArrayProperties = []*pb.IntArrayProperties{}
	textArrayProperties = []*pb.TextArrayProperties{}
	booleanArrayProperties = []*pb.BooleanArrayProperties{}
	objectProperties = []*pb.ObjectProperties{}
	objectArrayProperties = []*pb.ObjectArrayProperties{}
	singleTargetRefProps = []*pb.BatchObject_SingleTargetRefProps{}
	multiTargetRefProps = []*pb.BatchObject_MultiTargetRefProps{}
	for name, value := range properties {
		switch v := value.(type) {
		case bool, int, int32, int64, uint, uint32, uint64, float32, float64, string:
			nonRefPropertiesMap[name] = v
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
		case map[string]interface{}:
			// Object Property
			nonRefProps, numberArrayProps, intArrayProps,
				textArrayProps, booleanArrayProps, objectProps, objectArrayProps, _, _, objPropErr := c.extractProperties(v, false)
			if objPropErr != nil {
				err = fmt.Errorf("object properties: object property: %w", objPropErr)
				return
			}
			objectPropertiesValue := &pb.ObjectPropertiesValue{
				NonRefProperties:       nonRefProps,
				NumberArrayProperties:  numberArrayProps,
				IntArrayProperties:     intArrayProps,
				TextArrayProperties:    textArrayProps,
				BooleanArrayProperties: booleanArrayProps,
				ObjectProperties:       objectProps,
				ObjectArrayProperties:  objectArrayProps,
			}
			objectProperties = append(objectProperties, &pb.ObjectProperties{
				PropName: name,
				Value:    objectPropertiesValue,
			})
		case []map[string]interface{}:
			// it's a cross ref
			if rootLevel {
				crossRefs := map[string][]string{}
				var crossRefsErr error
				for i := range v {
					if len(v[i]) == 1 {
						valueString, ok := v[i]["beacon"].(string)
						if !ok {
							err = fmt.Errorf("cross-reference property has no beacon field")
							return
						}
						crossRefs, crossRefsErr = c.extractCrossRefs(crossRefs, valueString)
						if crossRefsErr != nil {
							err = fmt.Errorf("cross-reference property: %w", crossRefsErr)
							return
						}
					}
				}
				singleRefProps, multiRefProps := c.getCrossRefs(name, crossRefs)
				singleTargetRefProps = append(singleTargetRefProps, singleRefProps...)
				multiTargetRefProps = append(multiTargetRefProps, multiRefProps...)
			}
		case []map[string]string:
			if rootLevel {
				crossRefs := map[string][]string{}
				var crossRefsErr error
				for i := range v {
					if len(v[i]) == 1 {
						crossRefs, crossRefsErr = c.extractCrossRefs(crossRefs, v[i]["beacon"])
						if crossRefsErr != nil {
							err = fmt.Errorf("cross-reference property: %w", crossRefsErr)
							return
						}
					}
				}
				singleRefProps, multiRefProps := c.getCrossRefs(name, crossRefs)
				singleTargetRefProps = append(singleTargetRefProps, singleRefProps...)
				multiTargetRefProps = append(multiTargetRefProps, multiRefProps...)
			}
		case []interface{}:
			// Object Array Property
			objectArrayPropertiesValues := []*pb.ObjectPropertiesValue{}
			for _, objArrVal := range v {
				switch objArrValTyped := objArrVal.(type) {
				case map[string]interface{}:
					nonRefProps, numberArrayProps, intArrayProps,
						textArrayProps, booleanArrayProps, objectProps, objectArrayProps, _, _, objPropErr := c.extractProperties(objArrValTyped, false)
					if objPropErr != nil {
						err = fmt.Errorf("object properties: object array property: %w", objPropErr)
						return
					}
					objectPropertiesValue := &pb.ObjectPropertiesValue{
						NonRefProperties:       nonRefProps,
						NumberArrayProperties:  numberArrayProps,
						IntArrayProperties:     intArrayProps,
						TextArrayProperties:    textArrayProps,
						BooleanArrayProperties: booleanArrayProps,
						ObjectProperties:       objectProps,
						ObjectArrayProperties:  objectArrayProps,
					}
					objectArrayPropertiesValues = append(objectArrayPropertiesValues, objectPropertiesValue)
				default:
					err = fmt.Errorf("object properties: object array property: unsupported type: %T", objArrVal)
					return
				}
			}
			objectArrayProperties = append(objectArrayProperties, &pb.ObjectArrayProperties{
				PropName: name,
				Values:   objectArrayPropertiesValues,
			})
		}
		nonRefProperties, err = structpb.NewStruct(nonRefPropertiesMap)
		if err != nil {
			err = fmt.Errorf("object properties: %w", err)
			return
		}
	}
	return
}

func (c *GrpcClient) extractCrossRefs(crossRefs map[string][]string, beacon string) (map[string][]string, error) {
	cref, err := crossref.Parse(beacon)
	if err != nil {
		return nil, err
	}
	_, ok := crossRefs[cref.Class]
	if !ok {
		crossRefs[cref.Class] = []string{cref.TargetID.String()}
	} else {
		crossRefs[cref.Class] = append(crossRefs[cref.Class], cref.TargetID.String())
	}
	return crossRefs, nil
}

func (c *GrpcClient) getCrossRefs(propName string, crossRefs map[string][]string) ([]*pb.BatchObject_SingleTargetRefProps, []*pb.BatchObject_MultiTargetRefProps) {
	singleTargetRefProps := []*pb.BatchObject_SingleTargetRefProps{}
	multiTargetRefProps := []*pb.BatchObject_MultiTargetRefProps{}
	if len(crossRefs) == 1 {
		for key := range crossRefs {
			singleTargetRefProps = append(singleTargetRefProps, &pb.BatchObject_SingleTargetRefProps{
				PropName: propName,
				Uuids:    crossRefs[key],
			})
		}
	} else {
		for key := range crossRefs {
			multiTargetRefProps = append(multiTargetRefProps, &pb.BatchObject_MultiTargetRefProps{
				PropName:         propName,
				Uuids:            crossRefs[key],
				TargetCollection: key,
			})
		}
	}
	return singleTargetRefProps, multiTargetRefProps
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
