package batch

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/grpc/common"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema/crossref"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
	"github.com/weaviate/weaviate/usecases/byteops"
	"google.golang.org/protobuf/types/known/structpb"
)

type Batch struct {
	gRPCVersionSupport *db.GRPCVersionSupport
}

func New(gRPCVersionSupport *db.GRPCVersionSupport) Batch {
	return Batch{gRPCVersionSupport}
}

func (b Batch) GetBatchObjects(objects []*models.Object) ([]*pb.BatchObject, error) {
	result := make([]*pb.BatchObject, len(objects))
	for i, obj := range objects {
		if obj == nil {
			return nil, fmt.Errorf("object at index %d is nil", i)
		}
		properties, err := b.getProperties(obj.Properties)
		if err != nil {
			return nil, err
		}
		uid := obj.ID.String()
		if obj.ID == "" {
			uid = uuid.New().String()
		}
		batchObject := &pb.BatchObject{
			Uuid:       uid,
			Collection: obj.Class,
			Vector:     obj.Vector,
			Tenant:     obj.Tenant,
			Properties: properties,
		}
		if obj.Vector != nil {
			if b.gRPCVersionSupport.SupportsVectorBytesField() {
				batchObject.VectorBytes = byteops.Fp32SliceToBytes(obj.Vector)
			} else {
				// We fall back to vector field for backward compatibility reasons
				batchObject.Vector = obj.Vector
			}
		}
		if len(obj.Vectors) > 0 {
			vectors := []*pb.Vectors{}
			for targetVector, vector := range obj.Vectors {
				switch v := vector.(type) {
				case []float32:
					vectors = append(vectors, &pb.Vectors{
						Name:        targetVector,
						VectorBytes: byteops.Fp32SliceToBytes(v),
						Type:        pb.Vectors_VECTOR_TYPE_SINGLE_FP32,
					})
				case [][]float32:
					vectors = append(vectors, &pb.Vectors{
						Name:        targetVector,
						VectorBytes: byteops.Fp32SliceOfSlicesToBytes(v),
						Type:        pb.Vectors_VECTOR_TYPE_MULTI_FP32,
					})
				default:
					// do nothing
				}
			}
			batchObject.Vectors = vectors
		}
		result[i] = batchObject
	}
	return result, nil
}

func (b Batch) getProperties(properties models.PropertySchema) (*pb.BatchObject_Properties, error) {
	props, ok := properties.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("object properties type expected: map[string]interface{} got: %T", properties)
	}
	var result *pb.BatchObject_Properties
	if len(props) > 0 {
		nonRefPropsStruct, numberArrayProperties, intArrayProperties, textArrayProperties,
			booleanArrayProperties, objectProperties, objectArrayProperties,
			singleTargetRefProps, multiTargetRefProps, err := b.extractProperties(props, true)
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

func (b Batch) extractProperties(properties map[string]interface{}, rootLevel bool) (nonRefProperties *structpb.Struct,
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
				values = toInt64Array(vv)
			case []int32:
				values = toInt64Array(vv)
			case []int64:
				values = vv
			case []uint:
				values = toInt64Array(vv)
			case []uint32:
				values = toInt64Array(vv)
			case []uint64:
				values = toInt64Array(vv)
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
				textArrayProps, booleanArrayProps, objectProps, objectArrayProps, _, _, objPropErr := b.extractProperties(v, false)
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
						crossRefs, crossRefsErr = b.extractCrossRefs(crossRefs, valueString)
						if crossRefsErr != nil {
							err = fmt.Errorf("cross-reference property: %w", crossRefsErr)
							return
						}
					}
				}
				singleRefProps, multiRefProps := b.getCrossRefs(name, crossRefs)
				singleTargetRefProps = append(singleTargetRefProps, singleRefProps...)
				multiTargetRefProps = append(multiTargetRefProps, multiRefProps...)
			}
		case []map[string]string:
			if rootLevel {
				crossRefs := map[string][]string{}
				var crossRefsErr error
				for i := range v {
					if len(v[i]) == 1 {
						crossRefs, crossRefsErr = b.extractCrossRefs(crossRefs, v[i]["beacon"])
						if crossRefsErr != nil {
							err = fmt.Errorf("cross-reference property: %w", crossRefsErr)
							return
						}
					}
				}
				singleRefProps, multiRefProps := b.getCrossRefs(name, crossRefs)
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
						textArrayProps, booleanArrayProps, objectProps, objectArrayProps, _, _, objPropErr := b.extractProperties(objArrValTyped, false)
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

func (b Batch) extractCrossRefs(crossRefs map[string][]string, beacon string) (map[string][]string, error) {
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

func (b Batch) getCrossRefs(propName string, crossRefs map[string][]string) ([]*pb.BatchObject_SingleTargetRefProps, []*pb.BatchObject_MultiTargetRefProps) {
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

func (b Batch) GetConsistencyLevel(consistencyLevel string) *pb.ConsistencyLevel {
	return common.GetConsistencyLevel(consistencyLevel)
}

func (b Batch) ParseReply(reply *pb.BatchObjectsReply, objects []*models.Object) []models.ObjectsGetResponse {
	errors := map[int]*models.ObjectsGetResponseAO2Result{}
	// parse errors
	if reply != nil && len(reply.Errors) > 0 {
		for _, res := range reply.Errors {
			errors[int(res.Index)] = b.getErrorGetResponse(res)
		}
	}
	// prepare response
	success := models.ObjectsGetResponseAO2ResultStatusSUCCESS
	result := make([]models.ObjectsGetResponse, len(objects))
	for i := range objects {
		if err, ok := errors[i]; ok {
			// error
			result[i] = models.ObjectsGetResponse{
				Object: *objects[i],
				Result: err,
			}
		} else {
			// success
			result[i] = models.ObjectsGetResponse{
				Object: *objects[i],
				Result: &models.ObjectsGetResponseAO2Result{
					Status: &success,
				},
			}
		}
	}
	return result
}

func (b Batch) getErrorGetResponse(res *pb.BatchObjectsReply_BatchError) *models.ObjectsGetResponseAO2Result {
	var errors *models.ErrorResponse
	if res.Error != "" {
		errors = &models.ErrorResponse{
			Error: []*models.ErrorResponseErrorItems0{
				{Message: res.Error},
			},
		}
	}
	failed := models.ObjectsGetResponseAO2ResultStatusFAILED
	return &models.ObjectsGetResponseAO2Result{
		Errors: errors,
		Status: &failed,
	}
}

func toInt64Array[T int | int32 | int64 | uint | uint32 | uint64](arr []T) []int64 {
	result := make([]int64, len(arr))
	for i, val := range arr {
		result[i] = int64(val)
	}
	return result
}
