package api_test

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

// MessageMarshalerTest tests that [api.Message] produces a correct request message.
//
// We do not verify the [Message.Method] part of the request, because the specific
// [api.MethodFunc] returned is an implementation detail and there's probably not
// a lot of room for error.
type MessageMarshalerTest[In api.RequestMessage, Out api.ReplyMessage] struct {
	testkit.Only

	name string
	req  api.Message[In, Out] // Request struct.
	want *In                  // Expected protobuf request message.
	err  testkit.Error        // Set to testkit.ExpectError to expect an error.
}

// testMessageMarshaler runs [MessageMarshalerTest] test cases.
func testMessageMarshaler[In api.RequestMessage, Out api.ReplyMessage](t *testing.T, tests []MessageMarshalerTest[In, Out]) {
	t.Helper()
	for _, tt := range testkit.WithOnly(t, tests) {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.req, "invalid test: nil req")

			body := tt.req.Body()
			require.NotNil(t, body, "request body")

			got, err := body.MarshalMessage()
			tt.err.Require(t, err)
			require.EqualExportedValues(t, tt.want, got)
		})
	}
}

var (
	singleVector      = []float32{1, 2, 3}
	singleVectorBytes = []byte{
		0x0, 0x0, 0x80, 0x3f,
		0x0, 0x0, 0x0, 0x40,
		0x0, 0x0, 0x40, 0x40,
	}
	multiVector      = [][]float32{{1, 2, 3}, {1, 2, 3}}
	multiVectorBytes = []byte{
		0x3, 0x0, // inner array size, uint16(3)
		0x0, 0x0, 0x80, 0x3f, // first vector
		0x0, 0x0, 0x0, 0x40,
		0x0, 0x0, 0x40, 0x40,
		0x0, 0x0, 0x80, 0x3f, // second vector
		0x0, 0x0, 0x0, 0x40,
		0x0, 0x0, 0x40, 0x40,
	}
)

func TestSearchRequest_MarshalMessage(t *testing.T) {
	// UUID is always included in the [proto.MetadataRequest].
	testMessageMarshaler(t, []MessageMarshalerTest[proto.SearchRequest, proto.SearchReply]{
		{
			name: "base options",
			req: &api.SearchRequest{
				RequestDefaults: api.RequestDefaults{
					CollectionName:   "Songs",
					Tenant:           "john_doe",
					ConsistencyLevel: api.ConsistencyLevelOne,
				},
				Limit:     1,
				Offset:    2,
				AutoLimit: 3,
				After:     uuid.Max,
			},
			want: &proto.SearchRequest{
				Metadata:         &proto.MetadataRequest{Uuid: true},
				Collection:       "Songs",
				Tenant:           "john_doe",
				ConsistencyLevel: testkit.Ptr(proto.ConsistencyLevel_CONSISTENCY_LEVEL_ONE),
				Limit:            1,
				Offset:           2,
				Autocut:          3,
				After:            uuid.Max.String(),
				Properties: &proto.PropertiesRequest{
					ReturnAllNonrefProperties: true,
				},
			},
		},
		{
			name: "return metadata",
			req: &api.SearchRequest{
				ReturnMetadata: api.ReturnMetadata{
					Distance:     true,
					Certainty:    true,
					CreatedAt:    true,
					LastUpdateAt: true,
					Score:        true,
					ExplainScore: true,
				},
				ReturnVectors: []string{"title_vec", "lyrics_vec"},
			},
			want: &proto.SearchRequest{
				Metadata: &proto.MetadataRequest{
					Uuid:               true,
					Distance:           true,
					Certainty:          true,
					CreationTimeUnix:   true,
					LastUpdateTimeUnix: true,
					Score:              true,
					ExplainScore:       true,
					Vectors:            []string{"title_vec", "lyrics_vec"},
				},
				Properties: &proto.PropertiesRequest{
					ReturnAllNonrefProperties: true,
				},
			},
		},
		{
			name: "return default vector",
			req: &api.SearchRequest{
				ReturnVectors: []string{},
			},
			want: &proto.SearchRequest{
				Metadata: &proto.MetadataRequest{
					Uuid:   true,
					Vector: true,
				},
				Properties: &proto.PropertiesRequest{
					ReturnAllNonrefProperties: true,
				},
			},
		},
		{
			name: "return properties",
			req: &api.SearchRequest{
				ReturnProperties: []api.ReturnProperty{
					{Name: "title"},
					{Name: "lyrics"},
				},
			},
			want: &proto.SearchRequest{
				Metadata: &proto.MetadataRequest{Uuid: true},
				Properties: &proto.PropertiesRequest{
					NonRefProperties: []string{"title", "lyrics"},
				},
			},
		},
		{
			name: "return all properties",
			req:  &api.SearchRequest{},
			want: &proto.SearchRequest{
				Metadata: &proto.MetadataRequest{Uuid: true},
				Properties: &proto.PropertiesRequest{
					ReturnAllNonrefProperties: true,
				},
			},
		},
		{
			name: "return object properties",
			req: &api.SearchRequest{
				ReturnProperties: []api.ReturnProperty{
					{
						Name: "label",
						NestedProperties: []api.ReturnProperty{
							{Name: "name"},
							{Name: "logo"},
							{
								Name: "address",
								NestedProperties: []api.ReturnProperty{
									{Name: "street"},
									{Name: "building_nr"},
								},
							},
							{
								Name: "equipment",
								NestedProperties: []api.ReturnProperty{
									{
										Name: "microphone",
										NestedProperties: []api.ReturnProperty{
											{Name: "price"},
										},
									},
									{
										Name: "headphones",
										NestedProperties: []api.ReturnProperty{
											{Name: "price"},
										},
									},
								},
							},
						},
					},
				},
			},
			want: &proto.SearchRequest{
				Metadata: &proto.MetadataRequest{Uuid: true},
				Properties: &proto.PropertiesRequest{
					ObjectProperties: []*proto.ObjectPropertiesRequest{
						{
							PropName:            "label",
							PrimitiveProperties: []string{"name", "logo"},
							ObjectProperties: []*proto.ObjectPropertiesRequest{
								{
									PropName:            "address",
									PrimitiveProperties: []string{"street", "building_nr"},
								},
								{
									PropName: "equipment",
									ObjectProperties: []*proto.ObjectPropertiesRequest{
										{
											PropName:            "microphone",
											PrimitiveProperties: []string{"price"},
										},
										{
											PropName:            "headphones",
											PrimitiveProperties: []string{"price"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "return references",
			req: &api.SearchRequest{
				ReturnProperties: []api.ReturnProperty{},
				ReturnReferences: []api.ReturnReference{
					{
						PropertyName:     "performedBy",
						ReturnProperties: []api.ReturnProperty{},
						ReturnMetadata: api.ReturnMetadata{
							CreatedAt: true,
						},
					},
					{
						PropertyName:     "hasAwards",
						TargetCollection: "GrammyAward",
						ReturnProperties: []api.ReturnProperty{
							{Name: "categories"},
						},
					},
					{
						PropertyName:     "hasAwards",
						TargetCollection: "TonyAward",
						ReturnVectors:    []string{"recoding_vec"},
					},
					{
						PropertyName:     "writtenBy",
						ReturnProperties: []api.ReturnProperty{},
						ReturnReferences: []api.ReturnReference{
							{
								PropertyName:     "belongsToBand",
								TargetCollection: "MetalBands",
								ReturnProperties: []api.ReturnProperty{
									{Name: "name"},
								},
								ReturnReferences: []api.ReturnReference{
									{PropertyName: "foundedBy"},
								},
							},
						},
					},
				},
			},
			want: &proto.SearchRequest{
				Metadata: &proto.MetadataRequest{Uuid: true},
				Properties: &proto.PropertiesRequest{
					RefProperties: []*proto.RefPropertiesRequest{
						{
							ReferenceProperty: "performedBy",
							Metadata: &proto.MetadataRequest{
								Uuid:             true,
								CreationTimeUnix: true,
							},
							Properties: new(proto.PropertiesRequest),
						},
						{
							ReferenceProperty: "hasAwards",
							TargetCollection:  "GrammyAward",
							Properties: &proto.PropertiesRequest{
								NonRefProperties: []string{"categories"},
							},
							Metadata: &proto.MetadataRequest{Uuid: true},
						},
						{
							ReferenceProperty: "hasAwards",
							TargetCollection:  "TonyAward",
							Properties: &proto.PropertiesRequest{
								ReturnAllNonrefProperties: true,
							},
							Metadata: &proto.MetadataRequest{
								Uuid:    true,
								Vectors: []string{"recoding_vec"},
							},
						},
						{
							ReferenceProperty: "writtenBy",
							Properties: &proto.PropertiesRequest{
								RefProperties: []*proto.RefPropertiesRequest{
									{
										ReferenceProperty: "belongsToBand",
										TargetCollection:  "MetalBands",
										Properties: &proto.PropertiesRequest{
											NonRefProperties: []string{"name"},
											RefProperties: []*proto.RefPropertiesRequest{
												{
													ReferenceProperty: "foundedBy",
													Properties: &proto.PropertiesRequest{
														ReturnAllNonrefProperties: true,
													},
													Metadata: &proto.MetadataRequest{Uuid: true},
												},
											},
										},
										Metadata: &proto.MetadataRequest{Uuid: true},
									},
								},
							},
							Metadata: &proto.MetadataRequest{Uuid: true},
						},
					},
				},
			},
		},
		{
			name: "near vector no targets",
			req: &api.SearchRequest{
				NearVector: &api.NearVector{
					Target: api.SearchTarget{
						Vectors: []api.TargetVector{},
					},
				},
			},
			want: &proto.SearchRequest{
				NearVector: nil,
				Metadata:   &proto.MetadataRequest{Uuid: true},
				Properties: &proto.PropertiesRequest{
					ReturnAllNonrefProperties: true,
				},
			},
		},
		{
			name: "near vector invalid target",
			req: &api.SearchRequest{
				NearVector: &api.NearVector{
					Target: api.SearchTarget{
						Vectors: []api.TargetVector{
							{Vector: api.Vector{Name: "empty"}},
						},
					},
				},
			},
			err: testkit.ExpectError,
		},
		{
			name: "near vector single target named",
			req: &api.SearchRequest{
				NearVector: &api.NearVector{
					Distance: testkit.Ptr(.123),
					Target: api.SearchTarget{
						Vectors: []api.TargetVector{
							{
								Vector: api.Vector{
									Name:   "title_vec",
									Single: singleVector,
								},
							},
						},
					},
				},
			},
			want: &proto.SearchRequest{
				NearVector: &proto.NearVector{
					Distance: testkit.Ptr(.123),
					VectorForTargets: []*proto.VectorForTarget{
						{
							Name: "title_vec",
							Vectors: []*proto.Vectors{
								{
									Name:        "title_vec",
									VectorBytes: singleVectorBytes,
									Type:        proto.Vectors_VECTOR_TYPE_SINGLE_FP32,
								},
							},
						},
					},
				},
				Metadata: &proto.MetadataRequest{Uuid: true},
				Properties: &proto.PropertiesRequest{
					ReturnAllNonrefProperties: true,
				},
			},
		},
		{
			name: "near vector single target anonymous",
			req: &api.SearchRequest{
				NearVector: &api.NearVector{
					Distance: testkit.Ptr(.123),
					Target: api.SearchTarget{
						Vectors: []api.TargetVector{
							{
								Vector: api.Vector{
									Single: singleVector,
								},
							},
						},
					},
				},
			},
			want: &proto.SearchRequest{
				NearVector: &proto.NearVector{
					Distance: testkit.Ptr(.123),
					Vectors: []*proto.Vectors{
						{
							VectorBytes: singleVectorBytes,
							Type:        proto.Vectors_VECTOR_TYPE_SINGLE_FP32,
						},
					},
				},
				Metadata: &proto.MetadataRequest{Uuid: true},
				Properties: &proto.PropertiesRequest{
					ReturnAllNonrefProperties: true,
				},
			},
		},
		{
			name: "near vector multi target average",
			req: &api.SearchRequest{
				NearVector: &api.NearVector{
					Target: api.SearchTarget{
						CombinationMethod: api.CombinationMethodAverage,
						Vectors: []api.TargetVector{
							{
								Vector: api.Vector{
									Name:   "title_vec",
									Single: singleVector,
								},
							},
							{
								Vector: api.Vector{
									Name:  "lyrics_vec",
									Multi: multiVector,
								},
							},
						},
					},
				},
			},
			want: &proto.SearchRequest{
				NearVector: &proto.NearVector{
					Targets: &proto.Targets{
						TargetVectors: []string{"title_vec", "lyrics_vec"},
						Combination:   proto.CombinationMethod_COMBINATION_METHOD_TYPE_AVERAGE,
					},
					VectorForTargets: []*proto.VectorForTarget{
						{
							Name: "title_vec",
							Vectors: []*proto.Vectors{
								{
									Name:        "title_vec",
									VectorBytes: singleVectorBytes,
									Type:        proto.Vectors_VECTOR_TYPE_SINGLE_FP32,
								},
							},
						},
						{
							Name: "lyrics_vec",
							Vectors: []*proto.Vectors{
								{
									Name:        "lyrics_vec",
									VectorBytes: multiVectorBytes,
									Type:        proto.Vectors_VECTOR_TYPE_MULTI_FP32,
								},
							},
						},
					},
				},
				Metadata: &proto.MetadataRequest{Uuid: true},
				Properties: &proto.PropertiesRequest{
					ReturnAllNonrefProperties: true,
				},
			},
		},
		{
			name: "near vector multi target manual weights",
			req: &api.SearchRequest{
				NearVector: &api.NearVector{
					Target: api.SearchTarget{
						CombinationMethod: api.CombinationMethodManualWeights,
						Vectors: []api.TargetVector{
							{
								Vector: api.Vector{
									Name:   "title_vec",
									Single: singleVector,
								},
								Weight: testkit.Ptr[float32](.4),
							},
							{
								Vector: api.Vector{
									Name:  "lyrics_vec",
									Multi: multiVector,
								},
								Weight: testkit.Ptr[float32](.6),
							},
						},
					},
				},
			},
			want: &proto.SearchRequest{
				NearVector: &proto.NearVector{
					Targets: &proto.Targets{
						TargetVectors: []string{"title_vec", "lyrics_vec"},
						Combination:   proto.CombinationMethod_COMBINATION_METHOD_TYPE_MANUAL,
						WeightsForTargets: []*proto.WeightsForTarget{
							{
								Target: "title_vec",
								Weight: .4,
							},
							{
								Target: "lyrics_vec",
								Weight: .6,
							},
						},
					},
					VectorForTargets: []*proto.VectorForTarget{
						{
							Name: "title_vec",
							Vectors: []*proto.Vectors{
								{
									Name:        "title_vec",
									VectorBytes: singleVectorBytes,
									Type:        proto.Vectors_VECTOR_TYPE_SINGLE_FP32,
								},
							},
						},
						{
							Name: "lyrics_vec",
							Vectors: []*proto.Vectors{
								{
									Name:        "lyrics_vec",
									VectorBytes: multiVectorBytes,
									Type:        proto.Vectors_VECTOR_TYPE_MULTI_FP32,
								},
							},
						},
					},
				},
				Metadata: &proto.MetadataRequest{Uuid: true},
				Properties: &proto.PropertiesRequest{
					ReturnAllNonrefProperties: true,
				},
			},
		},
	})
}
