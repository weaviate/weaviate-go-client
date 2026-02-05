package api_test

import (
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"google.golang.org/protobuf/types/known/structpb"
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
	err  testkit.Error
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

// ----------------------------------------------------------------------------

type MessageUnmarshalerTest[Out api.ReplyMessage] struct {
	testkit.Only

	name  string
	reply *Out                        // Protobuf message that needs to be unmarshaled.
	dest  api.MessageUnmarshaler[Out] // Set dest to a pointer to response struct.
	want  any                         // Expected response value (pointer).
	err   testkit.Error
}

// testMessageMarshaler runs test cases for [api.MessageUnmarshaler] implementations.
func testMessageUnmarshaler[Out api.ReplyMessage](t *testing.T, tests []MessageUnmarshalerTest[Out]) {
	t.Helper()
	for _, tt := range testkit.WithOnly(t, tests) {
		t.Run(tt.name, func(t *testing.T) {
			testkit.IsPointer(t, tt.want, "want")

			err := tt.dest.UnmarshalMessage(tt.reply)
			tt.err.Require(t, err, "unmarshal")
			require.EqualExportedValues(t, tt.want, tt.dest)
		})
	}
}

func TestSearchResponse_UnmarshalMessage(t *testing.T) {
	idAsBytes, err := testkit.UUID.MarshalBinary()
	require.NoError(t, err, "marshal uuid bytes")

	testMessageUnmarshaler(t, []MessageUnmarshalerTest[proto.SearchReply]{
		{
			name: "metadata",
			reply: &proto.SearchReply{
				Took: 92,
				Results: []*proto.SearchResult{
					{
						Metadata: &proto.MetadataResult{
							IdAsBytes: idAsBytes,
							Distance:  .123, DistancePresent: true,
							Certainty: .123, CertaintyPresent: false,
							Score: .456, ScorePresent: true,
							ExplainScore: "very good", ExplainScorePresent: false,
							CreationTimeUnix: testkit.Now.UnixMilli(), CreationTimeUnixPresent: true,
							LastUpdateTimeUnix: testkit.Now.UnixMilli(), LastUpdateTimeUnixPresent: false,
							Vectors: []*proto.Vectors{
								{
									Name:        "title_vec",
									VectorBytes: singleVectorBytes,
									Type:        proto.Vectors_VECTOR_TYPE_SINGLE_FP32,
								},
								{
									Name:        "lyrics_vec",
									VectorBytes: multiVectorBytes,
									Type:        proto.Vectors_VECTOR_TYPE_MULTI_FP32,
								},
							},
						},
					},
				},
			},
			dest: new(api.SearchResponse),
			want: &api.SearchResponse{
				Took: 92 * time.Second,
				Results: []api.Object{
					{
						Metadata: api.ObjectMetadata{
							UUID:          testkit.UUID,
							Distance:      testkit.Ptr[float32](.123),
							Score:         testkit.Ptr[float32](.456),
							CreatedAt:     testkit.Ptr(testkit.Now),
							Certainty:     nil, // present == false
							ExplainScore:  nil, // present == false
							LastUpdatedAt: nil, // present == false
							NamedVectors: api.Vectors{
								"title_vec": api.Vector{
									Name:   "title_vec",
									Single: singleVector,
								},
								"lyrics_vec": api.Vector{
									Name:  "lyrics_vec",
									Multi: multiVector,
								},
							},
						},
						Properties: make(map[string]any),
						References: make(map[string][]api.Object),
					},
				},
				GroupByResults: make([]api.Group, 0),
			},
		},
		{
			name: "unnamed vector",
			reply: &proto.SearchReply{
				Results: []*proto.SearchResult{
					{
						Metadata: &proto.MetadataResult{
							VectorBytes: singleVectorBytes,
						},
					},
				},
			},
			dest: new(api.SearchResponse),
			want: &api.SearchResponse{
				Results: []api.Object{
					{
						Metadata: api.ObjectMetadata{
							UnnamedVector: &api.Vector{
								Name:   api.DefaultVectorName,
								Single: singleVector,
							},
							NamedVectors: make(api.Vectors),
						},
						Properties: make(map[string]any),
						References: make(map[string][]api.Object),
					},
				},
				GroupByResults: make([]api.Group, 0),
			},
		},
		{
			name: "missing vector type",
			reply: &proto.SearchReply{
				Results: []*proto.SearchResult{
					{
						Metadata: &proto.MetadataResult{
							Vectors: []*proto.Vectors{
								{Name: "no_type", VectorBytes: singleVectorBytes},
							},
						},
					},
				},
			},
			dest: new(api.SearchResponse),
			want: new(api.SearchResponse),
			err:  testkit.ExpectError,
		},
		{
			name: "bad uuid",
			reply: &proto.SearchReply{
				Results: []*proto.SearchResult{
					{
						Metadata: &proto.MetadataResult{
							IdAsBytes: []byte("00-00-00"),
						},
					},
				},
			},
			dest: new(api.SearchResponse),
			want: new(api.SearchResponse),
			err:  testkit.ExpectError,
		},
		{
			name: "properties",
			reply: &proto.SearchReply{
				Results: []*proto.SearchResult{
					{
						Properties: &proto.PropertiesResult{
							TargetCollection: "Songs",
							NonRefProps: &proto.Properties{
								Fields: map[string]*proto.Value{
									"title":        text("High Speed Dirt"),
									"is_single":    boolean(false),
									"release_date": date(testkit.Now),
									"duration_sec": integer(252),
									"retail_price": number(53.99),
									"album_cover":  blob("cover.png"),
									"uuid":         UUID(testkit.UUID),
									"extra": object(map[string]*proto.Value{
										"key": text("D"),
									}),
									"kpop_version": null(),
								},
							},
						},
					},
				},
			},
			dest: new(api.SearchResponse),
			want: &api.SearchResponse{
				Results: []api.Object{
					{
						Collection: "Songs",
						Properties: map[string]any{
							"title":        "High Speed Dirt",
							"is_single":    false,
							"release_date": testkit.Now,
							"duration_sec": int64(252),
							"retail_price": float64(53.99),
							"album_cover":  "cover.png",
							"uuid":         testkit.UUID,
							"extra": map[string]any{
								"key": "D",
							},
							"kpop_version": nil,
						},
						References: make(map[string][]api.Object),
					},
				},
				GroupByResults: make([]api.Group, 0),
			},
		},
		{
			name: "references",
			reply: &proto.SearchReply{
				Results: []*proto.SearchResult{
					{
						Properties: &proto.PropertiesResult{
							RefProps: []*proto.RefPropertiesResult{
								{
									PropName: "hasAwards",
									Properties: []*proto.PropertiesResult{
										{
											TargetCollection: "GrammyAward",
											NonRefProps: &proto.Properties{
												Fields: map[string]*proto.Value{
													"category": text("metal"),
												},
											},
										},
										{
											TargetCollection: "TonyAward",
											Metadata: &proto.MetadataResult{
												IdAsBytes: idAsBytes,
											},
										},
									},
								},
								{
									PropName: "writtenBy",
									Properties: []*proto.PropertiesResult{
										{
											TargetCollection: "Artists",
											RefProps: []*proto.RefPropertiesResult{
												{
													PropName: "belongsToBand",
													Properties: []*proto.PropertiesResult{
														{
															TargetCollection: "MetalBands",
															NonRefProps: &proto.Properties{
																Fields: map[string]*proto.Value{
																	"name": text("Megadeth"),
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			dest: new(api.SearchResponse),
			want: &api.SearchResponse{
				Results: []api.Object{
					{
						Properties: make(map[string]any),
						References: map[string][]api.Object{
							"hasAwards": {
								{
									Collection: "GrammyAward",
									Properties: map[string]any{
										"category": "metal",
									},
									References: make(map[string][]api.Object),
								},
								{
									Collection: "TonyAward",
									Metadata: api.ObjectMetadata{
										UUID:         testkit.UUID,
										NamedVectors: make(api.Vectors),
									},
									Properties: make(map[string]any),
									References: make(map[string][]api.Object),
								},
							},
							"writtenBy": {
								{
									Collection: "Artists",
									Properties: make(map[string]any),
									References: map[string][]api.Object{
										"belongsToBand": {
											{
												Collection: "MetalBands",
												Properties: map[string]any{
													"name": "Megadeth",
												},
												References: make(map[string][]api.Object),
											},
										},
									},
								},
							},
						},
					},
				},
				GroupByResults: make([]api.Group, 0),
			},
		},
		{
			name: "grouped result",
			reply: &proto.SearchReply{
				GroupByResults: []*proto.GroupByResult{
					{
						Name:            "metadata and properties",
						MinDistance:     .05,
						MaxDistance:     .1,
						NumberOfObjects: 2,
						Objects: []*proto.SearchResult{
							{
								Metadata: &proto.MetadataResult{
									IdAsBytes: idAsBytes,
									Distance:  .123, DistancePresent: true,
									Certainty: .123, CertaintyPresent: false,
									Score: .456, ScorePresent: true,
									ExplainScore: "very good", ExplainScorePresent: false,
									CreationTimeUnix: testkit.Now.UnixMilli(), CreationTimeUnixPresent: true,
									LastUpdateTimeUnix: testkit.Now.UnixMilli(), LastUpdateTimeUnixPresent: false,
									Vectors: []*proto.Vectors{
										{
											Name:        "title_vec",
											VectorBytes: singleVectorBytes,
											Type:        proto.Vectors_VECTOR_TYPE_SINGLE_FP32,
										},
										{
											Name:        "lyrics_vec",
											VectorBytes: multiVectorBytes,
											Type:        proto.Vectors_VECTOR_TYPE_MULTI_FP32,
										},
									},
								},
							},
							{
								Properties: &proto.PropertiesResult{
									TargetCollection: "Songs",
									NonRefProps: &proto.Properties{
										Fields: map[string]*proto.Value{
											"title":        text("High Speed Dirt"),
											"is_single":    boolean(false),
											"release_date": date(testkit.Now),
											"duration_sec": integer(252),
											"retail_price": number(53.99),
											"album_cover":  blob("cover.png"),
											"uuid":         UUID(testkit.UUID),
											"extra": object(map[string]*proto.Value{
												"key": text("D"),
											}),
											"kpop_version": null(),
										},
									},
								},
							},
						},
					},
					{
						Name:            "references",
						MinDistance:     .6,
						MaxDistance:     .7,
						NumberOfObjects: 1,
						Objects: []*proto.SearchResult{
							{
								Properties: &proto.PropertiesResult{
									RefProps: []*proto.RefPropertiesResult{
										{
											PropName: "hasAwards",
											Properties: []*proto.PropertiesResult{
												{
													TargetCollection: "GrammyAward",
													NonRefProps: &proto.Properties{
														Fields: map[string]*proto.Value{
															"category": text("metal"),
														},
													},
												},
												{
													TargetCollection: "TonyAward",
													Metadata: &proto.MetadataResult{
														IdAsBytes: idAsBytes,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			dest: new(api.SearchResponse),
			want: &api.SearchResponse{
				Results: make([]api.Object, 0),
				GroupByResults: []api.Group{
					{
						Name:        "metadata and properties",
						MinDistance: .05,
						MaxDistance: .1,
						Size:        2,
						Objects: []api.GroupObject{
							{
								BelongsToGroup: "metadata and properties",
								Object: api.Object{
									Metadata: api.ObjectMetadata{
										UUID:          testkit.UUID,
										Distance:      testkit.Ptr[float32](.123),
										Score:         testkit.Ptr[float32](.456),
										CreatedAt:     testkit.Ptr(testkit.Now),
										Certainty:     nil, // present == false
										ExplainScore:  nil, // present == false
										LastUpdatedAt: nil, // present == false
										NamedVectors: api.Vectors{
											"title_vec": api.Vector{
												Name:   "title_vec",
												Single: singleVector,
											},
											"lyrics_vec": api.Vector{
												Name:  "lyrics_vec",
												Multi: multiVector,
											},
										},
									},
									Properties: make(map[string]any),
									References: make(map[string][]api.Object),
								},
							},
							{
								BelongsToGroup: "metadata and properties",
								Object: api.Object{
									Collection: "Songs",
									Properties: map[string]any{
										"title":        "High Speed Dirt",
										"is_single":    false,
										"release_date": testkit.Now,
										"duration_sec": int64(252),
										"retail_price": float64(53.99),
										"album_cover":  "cover.png",
										"uuid":         testkit.UUID,
										"extra": map[string]any{
											"key": "D",
										},
										"kpop_version": nil,
									},
									References: make(map[string][]api.Object),
								},
							},
						},
					},
					{
						Name:        "references",
						MinDistance: .6,
						MaxDistance: .7,
						Size:        1,
						Objects: []api.GroupObject{
							{
								BelongsToGroup: "references",
								Object: api.Object{
									Properties: make(map[string]any),
									References: map[string][]api.Object{
										"hasAwards": {
											{
												Collection: "GrammyAward",
												Properties: map[string]any{
													"category": "metal",
												},
												References: make(map[string][]api.Object),
											},
											{
												Collection: "TonyAward",
												Metadata: api.ObjectMetadata{
													UUID:         testkit.UUID,
													NamedVectors: make(api.Vectors),
												},
												Properties: make(map[string]any),
												References: make(map[string][]api.Object),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
}

func text(s string) *proto.Value {
	return &proto.Value{Kind: &proto.Value_TextValue{TextValue: s}}
}

func boolean(b bool) *proto.Value {
	return &proto.Value{Kind: &proto.Value_BoolValue{BoolValue: b}}
}

func date(t time.Time) *proto.Value {
	return &proto.Value{Kind: &proto.Value_DateValue{DateValue: t.Format(api.TimeLayout)}}
}

func integer(i int64) *proto.Value {
	return &proto.Value{Kind: &proto.Value_IntValue{IntValue: i}}
}

func number(f float64) *proto.Value {
	return &proto.Value{Kind: &proto.Value_NumberValue{NumberValue: f}}
}

func blob(s string) *proto.Value {
	return &proto.Value{Kind: &proto.Value_BlobValue{BlobValue: s}}
}

func null() *proto.Value {
	return &proto.Value{Kind: &proto.Value_NullValue{NullValue: structpb.NullValue_NULL_VALUE}}
}

func UUID(u uuid.UUID) *proto.Value {
	return &proto.Value{Kind: &proto.Value_UuidValue{UuidValue: u.String()}}
}

func object(m map[string]*proto.Value) *proto.Value {
	return &proto.Value{Kind: &proto.Value_ObjectValue{ObjectValue: &proto.Properties{
		Fields: m,
	}}}
}
