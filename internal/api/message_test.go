package api_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"google.golang.org/protobuf/types/known/structpb"
)

// MessageMarshalerTest tests that [transport.Message] produces a correct request message.
//
// We do not verify the [Message.Method] part of the request, because the specific
// [api.MethodFunc] returned is an implementation detail and there's probably not
// a lot of room for error.
type MessageMarshalerTest[In transport.RequestMessage, Out transport.ReplyMessage] struct {
	testkit.Only

	name string
	req  transport.Message[In, Out] // Request struct.
	want *In                        // Expected protobuf request message.
	err  testkit.Error
}

// testMessageMarshaler runs [MessageMarshalerTest] test cases.
func testMessageMarshaler[In transport.RequestMessage, Out transport.ReplyMessage](t *testing.T, tests []MessageMarshalerTest[In, Out]) {
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
					Similarity: api.VectorSimilarity{Distance: testkit.Ptr(.123)},
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
					Targets: &proto.Targets{
						TargetVectors: []string{"title_vec"},
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
					Similarity: api.VectorSimilarity{Distance: testkit.Ptr(.123)},
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
					Targets: &proto.Targets{
						TargetVectors: []string{""},
					},
					VectorForTargets: []*proto.VectorForTarget{
						{
							Vectors: []*proto.Vectors{
								{
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
		{
			name: "near vector multiple vectors per target",
			req: &api.SearchRequest{
				NearVector: &api.NearVector{
					Target: api.SearchTarget{
						CombinationMethod: api.CombinationMethodRelativeScore,
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
									Name:   "title_vec",
									Single: singleVector,
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
						TargetVectors: []string{"title_vec"},
						Combination:   proto.CombinationMethod_COMBINATION_METHOD_TYPE_RELATIVE_SCORE,
						WeightsForTargets: []*proto.WeightsForTarget{
							{
								Target: "title_vec",
								Weight: .4,
							},
							{
								Target: "title_vec",
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
			name: "near text implicit target",
			req: &api.SearchRequest{
				NearText: &api.NearText{
					Concepts:   []string{"apples", "oranges"},
					Similarity: api.VectorSimilarity{Distance: testkit.Ptr(.92)},
					MoveTo: &api.Move{
						Force:    0.58,
						Concepts: []string{"computers"},
					},
					MoveAway: &api.Move{
						Force:   0.22,
						Objects: []uuid.UUID{uuid.Nil, uuid.Max},
					},
					Selection: api.Selection{
						MMR: &api.SelectionMMR{Limit: int32(3)},
					},
				},
			},
			want: &proto.SearchRequest{
				NearText: &proto.NearTextSearch{
					Query:    []string{"apples", "oranges"},
					Distance: testkit.Ptr(.92),
					MoveTo: &proto.NearTextSearch_Move{
						Force:    0.58,
						Concepts: []string{"computers"},
					},
					MoveAway: &proto.NearTextSearch_Move{
						Force: 0.22,
						Uuids: []string{uuid.Nil.String(), uuid.Max.String()},
					},
					Selection: &proto.Selection{
						Selection: &proto.Selection_Mmr{
							Mmr: &proto.Selection_MMR{
								Limit: testkit.Ptr[uint32](3),
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
			name: "near text explicit target",
			req: &api.SearchRequest{
				NearText: &api.NearText{
					Concepts: []string{"apples", "oranges"},
					Target: api.SearchTarget{
						Vectors: []api.TargetVector{
							{Vector: api.Vector{Name: "title_vec"}},
						},
					},
				},
			},
			want: &proto.SearchRequest{
				NearText: &proto.NearTextSearch{
					Query: []string{"apples", "oranges"},
					Targets: &proto.Targets{
						TargetVectors: []string{"title_vec"},
					},
				},
				Metadata: &proto.MetadataRequest{Uuid: true},
				Properties: &proto.PropertiesRequest{
					ReturnAllNonrefProperties: true,
				},
			},
		},
		{
			name: "near text multi target sum",
			req: &api.SearchRequest{
				NearText: &api.NearText{
					Concepts: []string{"apples", "oranges"},
					Target: api.SearchTarget{
						Vectors: []api.TargetVector{
							{Vector: api.Vector{Name: "title_vec"}},
							{Vector: api.Vector{Name: "lyrics_vec"}},
						},
						CombinationMethod: api.CombinationMethodSum,
					},
				},
			},
			want: &proto.SearchRequest{
				NearText: &proto.NearTextSearch{
					Query: []string{"apples", "oranges"},
					Targets: &proto.Targets{
						TargetVectors: []string{"title_vec", "lyrics_vec"},
						Combination:   proto.CombinationMethod_COMBINATION_METHOD_TYPE_SUM,
					},
				},
				Metadata: &proto.MetadataRequest{Uuid: true},
				Properties: &proto.PropertiesRequest{
					ReturnAllNonrefProperties: true,
				},
			},
		},
		{
			name: "near text multi target relative score",
			req: &api.SearchRequest{
				NearText: &api.NearText{
					Concepts: []string{"apples", "oranges"},
					Target: api.SearchTarget{
						Vectors: []api.TargetVector{
							{
								Vector: api.Vector{Name: "title_vec"},
								Weight: testkit.Ptr[float32](.11),
							},
							{
								Vector: api.Vector{Name: "lyrics_vec"},
								Weight: testkit.Ptr[float32](.22),
							},
						},
						CombinationMethod: api.CombinationMethodRelativeScore,
					},
				},
			},
			want: &proto.SearchRequest{
				NearText: &proto.NearTextSearch{
					Query: []string{"apples", "oranges"},
					Targets: &proto.Targets{
						TargetVectors: []string{"title_vec", "lyrics_vec"},
						Combination:   proto.CombinationMethod_COMBINATION_METHOD_TYPE_RELATIVE_SCORE,
						WeightsForTargets: []*proto.WeightsForTarget{
							{Target: "title_vec", Weight: .11},
							{Target: "lyrics_vec", Weight: .22},
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
			name: "hybrid with near vector",
			req: &api.SearchRequest{
				Hybrid: &api.Hybrid{
					Query:           "yellow submarine",
					QueryProperties: []string{"title", "lyrics"},
					Alpha:           testkit.Ptr[float32](.44),
					Fusion:          api.HybridFusionRanked,
					KeywordSimilarity: api.KeywordSimilarity{
						AllTokensMatch: true,
					},
					NearVector: &api.NearVector{
						Similarity: api.VectorSimilarity{Distance: testkit.Ptr(1.23)},
						Target: api.SearchTarget{
							Vectors: []api.TargetVector{
								{Vector: api.Vector{
									Name:   "lyrics_vec",
									Single: singleVector,
								}},
							},
						},
					},
				},
			},
			want: &proto.SearchRequest{
				HybridSearch: &proto.Hybrid{
					Query:      "yellow submarine",
					Properties: []string{"title", "lyrics"},
					AlphaParam: testkit.Ptr[float32](.44),
					FusionType: proto.Hybrid_FUSION_TYPE_RANKED,
					Bm25SearchOperator: &proto.SearchOperatorOptions{
						Operator: proto.SearchOperatorOptions_OPERATOR_AND,
					},
					NearVector: &proto.NearVector{
						Distance: testkit.Ptr(1.23),
						Targets: &proto.Targets{
							TargetVectors: []string{"lyrics_vec"},
						},
						VectorForTargets: []*proto.VectorForTarget{
							{
								Name: "lyrics_vec",
								Vectors: []*proto.Vectors{
									{
										Name:        "lyrics_vec",
										VectorBytes: singleVectorBytes,
										Type:        proto.Vectors_VECTOR_TYPE_SINGLE_FP32,
									},
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
			name: "hybrid with near text",
			req: &api.SearchRequest{
				Hybrid: &api.Hybrid{
					Query:           "yellow submarine",
					QueryProperties: []string{"title", "lyrics"},
					Alpha:           testkit.Ptr[float32](.44),
					Fusion:          api.HybridFusionRanked,
					KeywordSimilarity: api.KeywordSimilarity{
						MinimumTokensMatch: testkit.Ptr[int32](1),
					},
					NearText: &api.NearText{
						Concepts:   []string{"apples", "oranges"},
						Similarity: api.VectorSimilarity{Distance: testkit.Ptr(1.23)},
						Selection: api.Selection{
							MMR: &api.SelectionMMR{Limit: int32(3)},
						},
						Target: api.SearchTarget{
							Vectors: []api.TargetVector{
								{Vector: api.Vector{Name: "title_vec"}},
							},
						},
					},
				},
			},
			want: &proto.SearchRequest{
				HybridSearch: &proto.Hybrid{
					Query:      "yellow submarine",
					Properties: []string{"title", "lyrics"},
					AlphaParam: testkit.Ptr[float32](.44),
					FusionType: proto.Hybrid_FUSION_TYPE_RANKED,
					Bm25SearchOperator: &proto.SearchOperatorOptions{
						Operator:             proto.SearchOperatorOptions_OPERATOR_OR,
						MinimumOrTokensMatch: testkit.Ptr[int32](1),
					},
					NearText: &proto.NearTextSearch{
						Query:    []string{"apples", "oranges"},
						Distance: testkit.Ptr(1.23),
						Targets: &proto.Targets{
							TargetVectors: []string{"title_vec"},
						},
						Selection: &proto.Selection{
							Selection: &proto.Selection_Mmr{
								Mmr: &proto.Selection_MMR{
									Limit: testkit.Ptr[uint32](3),
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
			name: "group by",
			req: &api.SearchRequest{
				GroupBy: &api.GroupBy{
					Property:       "album",
					Limit:          6,
					NumberOfGroups: 7,
				},
			},
			want: &proto.SearchRequest{
				GroupBy: &proto.GroupBy{
					Path:            []string{"album"},
					ObjectsPerGroup: 6,
					NumberOfGroups:  7,
				},
				Metadata: &proto.MetadataRequest{Uuid: true},
				Properties: &proto.PropertiesRequest{
					ReturnAllNonrefProperties: true,
				},
			},
		},
	})
}

// TestAggregateRequest_MarshalMessage tests that api.AggregateRequest creates
// the expected proto.AggregateRequest when its MarshalMessage is called.
// Most of the test cases are split by property type to simplify error logs;
// there's nothing stopping a test case from having mixed property types though.
func TestAggregateRequest_MarshalMessage(t *testing.T) {
	testMessageMarshaler(t, []MessageMarshalerTest[proto.AggregateRequest, proto.AggregateReply]{
		{
			name: "with object limit",
			req:  &api.AggregateRequest{ObjectLimit: 10},
			want: &proto.AggregateRequest{ObjectLimit: testkit.Ptr[uint32](10)},
		},
		{
			name: "text properties",
			req: &api.AggregateRequest{
				Text: []api.AggregateTextRequest{
					{Property: "colour", Count: true, TopOccurrences: true},
					{Property: "tag", TopOccurrences: true, TopOccurencesCutoff: 10},
				},
			},
			want: &proto.AggregateRequest{
				Aggregations: []*proto.AggregateRequest_Aggregation{
					{Property: "colour", Aggregation: &proto.AggregateRequest_Aggregation_Text_{
						Text: &proto.AggregateRequest_Aggregation_Text{
							Count: true, TopOccurences: true,
						},
					}},
					{Property: "tag", Aggregation: &proto.AggregateRequest_Aggregation_Text_{
						Text: &proto.AggregateRequest_Aggregation_Text{
							TopOccurences: true, TopOccurencesLimit: testkit.Ptr[uint32](10),
						},
					}},
				},
			},
		},
		{
			name: "integer properties",
			req: &api.AggregateRequest{
				Integer: []api.AggregateIntegerRequest{
					{Property: "price", Sum: true, Min: true, Max: true},
					{Property: "size", Count: true, Mode: true, Median: true},
				},
			},
			want: &proto.AggregateRequest{
				Aggregations: []*proto.AggregateRequest_Aggregation{
					{Property: "price", Aggregation: &proto.AggregateRequest_Aggregation_Int{
						Int: &proto.AggregateRequest_Aggregation_Integer{
							Sum: true, Minimum: true, Maximum: true,
						},
					}},
					{Property: "size", Aggregation: &proto.AggregateRequest_Aggregation_Int{
						Int: &proto.AggregateRequest_Aggregation_Integer{
							Count: true, Mode: true, Median: true,
						},
					}},
				},
			},
		},
		{
			name: "number properties",
			req: &api.AggregateRequest{
				Number: []api.AggregateNumberRequest{
					{Property: "price", Sum: true, Min: true, Max: true},
					{Property: "size", Count: true, Mode: true, Median: true},
				},
			},
			want: &proto.AggregateRequest{
				Aggregations: []*proto.AggregateRequest_Aggregation{
					{Property: "price", Aggregation: &proto.AggregateRequest_Aggregation_Number_{
						Number: &proto.AggregateRequest_Aggregation_Number{
							Sum: true, Minimum: true, Maximum: true,
						},
					}},
					{Property: "size", Aggregation: &proto.AggregateRequest_Aggregation_Number_{
						Number: &proto.AggregateRequest_Aggregation_Number{
							Count: true, Mode: true, Median: true,
						},
					}},
				},
			},
		},
		{
			name: "boolean properties",
			req: &api.AggregateRequest{
				Boolean: []api.AggregateBooleanRequest{
					{Property: "onSale", Type: true, PercentageTrue: true, PercentageFalse: true},
					{Property: "newArrival", Count: true, TotalTrue: true, TotalFalse: true},
				},
			},
			want: &proto.AggregateRequest{
				Aggregations: []*proto.AggregateRequest_Aggregation{
					{Property: "onSale", Aggregation: &proto.AggregateRequest_Aggregation_Boolean_{
						Boolean: &proto.AggregateRequest_Aggregation_Boolean{
							Type: true, PercentageTrue: true, PercentageFalse: true,
						},
					}},
					{Property: "newArrival", Aggregation: &proto.AggregateRequest_Aggregation_Boolean_{
						Boolean: &proto.AggregateRequest_Aggregation_Boolean{
							Count: true, TotalTrue: true, TotalFalse: true,
						},
					}},
				},
			},
		},
		{
			name: "date properties",
			req: &api.AggregateRequest{
				Date: []api.AggregateDateRequest{
					{Property: "lastPurchase", Count: true, Min: true, Max: true},
					{Property: "lastReturn", Mode: true, Median: true},
				},
			},
			want: &proto.AggregateRequest{
				Aggregations: []*proto.AggregateRequest_Aggregation{
					{Property: "lastPurchase", Aggregation: &proto.AggregateRequest_Aggregation_Date_{
						Date: &proto.AggregateRequest_Aggregation_Date{
							Count: true, Minimum: true, Maximum: true,
						},
					}},
					{Property: "lastReturn", Aggregation: &proto.AggregateRequest_Aggregation_Date_{
						Date: &proto.AggregateRequest_Aggregation_Date{
							Mode: true, Median: true,
						},
					}},
				},
			},
		},
		{
			name: "group by",
			req: &api.AggregateRequest{
				RequestDefaults: api.RequestDefaults{
					CollectionName: "Songs",
				},
				GroupBy: &api.GroupBy{
					Property: "album",
					Limit:    92,
				},
			},
			want: &proto.AggregateRequest{
				Collection: "Songs",
				Limit:      testkit.Ptr[uint32](92),
				GroupBy: &proto.AggregateRequest_GroupBy{
					Collection: "Songs",
					Property:   "album",
				},
			},
		},
		{
			name: "count objects",
			req: &api.CountObjectsRequest{
				CollectionName: "Songs",
				Tenant:         "john_doe",
			},
			want: &proto.AggregateRequest{
				Collection:   "Songs",
				Tenant:       "john_doe",
				ObjectsCount: true,
			},
		},
	})

	t.Run("with query filter", func(t *testing.T) {
		for _, tt := range []struct {
			name string
			req  transport.Message[proto.AggregateRequest, proto.AggregateReply]
			get  func(*proto.AggregateRequest) any
			want any
		}{
			{
				name: "near vector",
				req: &api.AggregateRequest{
					NearVector: &api.NearVector{
						Similarity: api.VectorSimilarity{
							Distance: testkit.Ptr(0.5),
						},
						Target: api.SearchTarget{Vectors: []api.TargetVector{
							{Vector: api.Vector{Name: "1d", Single: singleVector}},
						}},
					},
				},
				get: returnAny((*proto.AggregateRequest).GetNearVector),
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				require.NotNil(t, tt.req, "invalid test: nil req")

				body := tt.req.Body()
				require.NotNil(t, body, "request body")

				message, err := body.MarshalMessage()
				require.Nil(t, err, "marshal error")

				require.NotNil(t, tt.get(message))
			})
		}
	})
}

// Wrap func(T) *U into a func(T) any.
func returnAny[T, U any](f func(*T) *U) func(*T) any {
	return func(req *T) any {
		return f(req)
	}
}

func TestObjectBatchRequest_MarshalMessage(t *testing.T) {
	// UUID is always included in the [proto.MetadataRequest].
	testMessageMarshaler(t, []MessageMarshalerTest[proto.BatchObjectsRequest, proto.BatchObjectsReply]{
		{
			name: "properties",
			req: &api.InsertObjectsRequest{
				RequestDefaults: api.RequestDefaults{
					CollectionName:   "Songs",
					Tenant:           "john_doe",
					ConsistencyLevel: api.ConsistencyLevelOne,
				},
				Objects: []api.BatchObject{
					{
						UUID: testkit.UUID,
						Properties: map[string]any{
							"artist": "Angine de Poitrine",
							"title":  "Mata Zyklek",
						},
					},
				},
			},
			want: &proto.BatchObjectsRequest{
				ConsistencyLevel: testkit.Ptr(proto.ConsistencyLevel_CONSISTENCY_LEVEL_ONE),
				Objects: []*proto.BatchObject{
					{
						Uuid:       testkit.UUID.String(),
						Collection: "Songs",
						Tenant:     "john_doe",
						Properties: &proto.BatchObject_Properties{
							NonRefProperties: mustNewStruct(map[string]any{
								"artist": "Angine de Poitrine",
								"title":  "Mata Zyklek",
							}),
						},
					},
				},
			},
		},
		{
			name: "references",
			req: &api.InsertObjectsRequest{
				RequestDefaults: api.RequestDefaults{
					CollectionName:   "Songs",
					Tenant:           "john_doe",
					ConsistencyLevel: api.ConsistencyLevelOne,
				},
				Objects: []api.BatchObject{
					{
						UUID: testkit.UUID,
						References: api.ObjectReferences{
							"performedBy": {
								{UUID: testkit.UUID, Collection: "Drummers"},
							},
							"onLabel": {{UUID: testkit.UUID}},
						},
					},
				},
			},
			want: &proto.BatchObjectsRequest{
				ConsistencyLevel: testkit.Ptr(proto.ConsistencyLevel_CONSISTENCY_LEVEL_ONE),
				Objects: []*proto.BatchObject{
					{
						Uuid:       testkit.UUID.String(),
						Collection: "Songs",
						Tenant:     "john_doe",
						Properties: &proto.BatchObject_Properties{
							MultiTargetRefProps: []*proto.BatchObject_MultiTargetRefProps{
								{
									PropName:         "performedBy",
									TargetCollection: "Drummers",
									Uuids:            []string{testkit.UUID.String()},
								},
							},
							SingleTargetRefProps: []*proto.BatchObject_SingleTargetRefProps{{
								PropName: "onLabel",
								Uuids:    []string{testkit.UUID.String()},
							}},
						},
					},
				},
			},
		},
		{
			name: "vectors",
			req: &api.InsertObjectsRequest{
				RequestDefaults: api.RequestDefaults{
					CollectionName:   "Songs",
					Tenant:           "john_doe",
					ConsistencyLevel: api.ConsistencyLevelOne,
				},
				Objects: []api.BatchObject{
					{
						UUID: testkit.UUID,
						Vectors: []api.Vector{
							{Name: "title_vec", Single: singleVector},
						},
					},
				},
			},
			want: &proto.BatchObjectsRequest{
				ConsistencyLevel: testkit.Ptr(proto.ConsistencyLevel_CONSISTENCY_LEVEL_ONE),
				Objects: []*proto.BatchObject{
					{
						Uuid:       testkit.UUID.String(),
						Collection: "Songs",
						Tenant:     "john_doe",
						Vectors: []*proto.Vectors{{
							Name:        "title_vec",
							Type:        proto.Vectors_VECTOR_TYPE_SINGLE_FP32,
							VectorBytes: singleVectorBytes,
						}},
					},
				},
			},
		},
	})
}

// mustNewStruct panics if [structpb.NewStruct] returns an error.
func mustNewStruct(m map[string]any) *structpb.Struct {
	s, err := structpb.NewStruct(m)
	if err != nil {
		panic(err)
	}
	return s
}

// ----------------------------------------------------------------------------

type MessageUnmarshalerTest[Out transport.ReplyMessage] struct {
	testkit.Only

	name  string
	reply *Out                              // Protobuf message that needs to be unmarshaled.
	dest  transport.MessageUnmarshaler[Out] // Set dest to a pointer to response struct.
	want  any                               // Expected response value (pointer).
	err   testkit.Error
}

// testMessageMarshaler runs test cases for [transport.MessageUnmarshaler] implementations.
func testMessageUnmarshaler[Out transport.ReplyMessage](t *testing.T, tests []MessageUnmarshalerTest[Out]) {
	t.Helper()
	for _, tt := range testkit.WithOnly(t, tests) {
		t.Run(tt.name, func(t *testing.T) {
			testkit.RequirePointer(t, tt.want, "want")
			require.NotNil(t, tt.dest, "bad dest")

			err := tt.dest.UnmarshalMessage(tt.reply)
			tt.err.Require(t, err, "unmarshal")
			require.EqualExportedValues(t, tt.want, tt.dest, "bad response")
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
							CreatedAt:     &testkit.Now,
							Certainty:     nil, // present == false
							ExplainScore:  nil, // present == false
							LastUpdatedAt: nil, // present == false
							Vectors: api.Vectors{
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
							Vectors: api.Vectors{
								api.DefaultVectorName: api.Vector{
									Name:   api.DefaultVectorName,
									Single: singleVector,
								},
							},
						},
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
						References: map[string][]api.Object{
							"hasAwards": {
								{
									Collection: "GrammyAward",
									Properties: map[string]any{
										"category": "metal",
									},
								},
								{
									Collection: "TonyAward",
									Metadata: api.ObjectMetadata{
										UUID: testkit.UUID,
									},
								},
							},
							"writtenBy": {
								{
									Collection: "Artists",
									References: map[string][]api.Object{
										"belongsToBand": {
											{
												Collection: "MetalBands",
												Properties: map[string]any{
													"name": "Megadeth",
												},
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
										CreatedAt:     &testkit.Now,
										Certainty:     nil, // present == false
										ExplainScore:  nil, // present == false
										LastUpdatedAt: nil, // present == false
										Vectors: api.Vectors{
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
									References: map[string][]api.Object{
										"hasAwards": {
											{
												Collection: "GrammyAward",
												Properties: map[string]any{
													"category": "metal",
												},
											},
											{
												Collection: "TonyAward",
												Metadata: api.ObjectMetadata{
													UUID: testkit.UUID,
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

// TestAggregateResponse_UnmarshalMessage tests that api.AggregateResponse reads
// proto.AggregateRequest correctly when its UnmarshalMessage is called.
func TestAggregateResponse_UnmarshalMessage(t *testing.T) {
	type Aggregations []*proto.AggregateReply_Aggregations_Aggregation

	// reply is a helper function to wrap returned aggregations in the layers or protobuf bureaucracy.
	reply := func(aggs Aggregations) *proto.AggregateReply_SingleResult {
		return &proto.AggregateReply_SingleResult{
			SingleResult: &proto.AggregateReply_Single{
				ObjectsCount: testkit.Ptr(int64(len(aggs))),
				Aggregations: &proto.AggregateReply_Aggregations{
					Aggregations: aggs,
				},
			},
		}
	}

	testMessageUnmarshaler(t, []MessageUnmarshalerTest[proto.AggregateReply]{
		{
			name: "text properties",
			reply: &proto.AggregateReply{
				Took: 92,
				Result: reply(Aggregations{
					{Property: "colour", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Text_{
						Text: &proto.AggregateReply_Aggregations_Aggregation_Text{
							Count: testkit.Ptr[int64](1),
							TopOccurences: &proto.AggregateReply_Aggregations_Aggregation_Text_TopOccurrences{
								Items: []*proto.AggregateReply_Aggregations_Aggregation_Text_TopOccurrences_TopOccurrence{
									{Value: "red", Occurs: 2},
									{Value: "blue", Occurs: 3},
								},
							},
						},
					}},
					{Property: "tag", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Text_{
						Text: &proto.AggregateReply_Aggregations_Aggregation_Text{
							TopOccurences: &proto.AggregateReply_Aggregations_Aggregation_Text_TopOccurrences{
								Items: []*proto.AggregateReply_Aggregations_Aggregation_Text_TopOccurrences_TopOccurrence{
									{Value: "casual", Occurs: 1},
									{Value: "comfy", Occurs: 2},
								},
							},
						},
					}},
				}),
			},
			dest: new(api.AggregateResponse),
			want: &api.AggregateResponse{
				Took: 92 * time.Second,
				Results: api.Aggregations{
					TotalCount: testkit.Ptr[int64](2),
					Text: []api.AggregateTextResult{
						{
							Property: "colour",
							Count:    testkit.Ptr[int64](1),
							TopOccurrences: []api.TopOccurrence{
								{Value: "red", OccursTimes: 2},
								{Value: "blue", OccursTimes: 3},
							},
						},
						{
							Property: "tag",
							TopOccurrences: []api.TopOccurrence{
								{Value: "casual", OccursTimes: 1},
								{Value: "comfy", OccursTimes: 2},
							},
						},
					},
				},
			},
		},
		{
			name: "integer properties",
			reply: &proto.AggregateReply{
				Took: 92,
				Result: reply(Aggregations{
					{Property: "price", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Int{
						Int: &proto.AggregateReply_Aggregations_Aggregation_Integer{
							Sum:     testkit.Ptr[int64](1),
							Minimum: testkit.Ptr[int64](2),
							Maximum: testkit.Ptr[int64](3),
						},
					}},
					{Property: "size", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Int{
						Int: &proto.AggregateReply_Aggregations_Aggregation_Integer{
							Count:  testkit.Ptr[int64](1),
							Mode:   testkit.Ptr[int64](2),
							Median: testkit.Ptr[float64](3),
						},
					}},
				}),
			},
			dest: new(api.AggregateResponse),
			want: &api.AggregateResponse{
				Took: 92 * time.Second,
				Results: api.Aggregations{
					TotalCount: testkit.Ptr[int64](2),
					Integer: []api.AggregateIntegerResult{
						{
							Property: "price",
							Sum:      testkit.Ptr[int64](1),
							Min:      testkit.Ptr[int64](2),
							Max:      testkit.Ptr[int64](3),
						},
						{
							Property: "size",
							Count:    testkit.Ptr[int64](1),
							Mode:     testkit.Ptr[int64](2),
							Median:   testkit.Ptr[float64](3),
						},
					},
				},
			},
		},
		{
			name: "number properties",
			reply: &proto.AggregateReply{
				Took: 92,
				Result: reply(Aggregations{
					{Property: "price", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Number_{
						Number: &proto.AggregateReply_Aggregations_Aggregation_Number{
							Sum:     testkit.Ptr[float64](1),
							Minimum: testkit.Ptr[float64](2),
							Maximum: testkit.Ptr[float64](3),
						},
					}},
					{Property: "size", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Number_{
						Number: &proto.AggregateReply_Aggregations_Aggregation_Number{
							Count:  testkit.Ptr[int64](1),
							Mode:   testkit.Ptr[float64](2),
							Median: testkit.Ptr[float64](3),
						},
					}},
				}),
			},
			dest: new(api.AggregateResponse),
			want: &api.AggregateResponse{
				Took: 92 * time.Second,
				Results: api.Aggregations{
					TotalCount: testkit.Ptr[int64](2),
					Number: []api.AggregateNumberResult{
						{
							Property: "price",
							Sum:      testkit.Ptr[float64](1),
							Min:      testkit.Ptr[float64](2),
							Max:      testkit.Ptr[float64](3),
						},
						{
							Property: "size",
							Count:    testkit.Ptr[int64](1),
							Mode:     testkit.Ptr[float64](2),
							Median:   testkit.Ptr[float64](3),
						},
					},
				},
			},
		},
		{
			name: "boolean properties",
			reply: &proto.AggregateReply{
				Took: 92,
				Result: reply(Aggregations{
					{Property: "onSale", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Boolean_{
						Boolean: &proto.AggregateReply_Aggregations_Aggregation_Boolean{
							Type:            testkit.Ptr("black_friday"),
							PercentageTrue:  testkit.Ptr[float64](1),
							PercentageFalse: testkit.Ptr[float64](2),
						},
					}},
					{Property: "newArrival", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Boolean_{
						Boolean: &proto.AggregateReply_Aggregations_Aggregation_Boolean{
							Count:      testkit.Ptr[int64](1),
							TotalTrue:  testkit.Ptr[int64](2),
							TotalFalse: testkit.Ptr[int64](3),
						},
					}},
				}),
			},
			dest: new(api.AggregateResponse),
			want: &api.AggregateResponse{
				Took: 92 * time.Second,
				Results: api.Aggregations{
					TotalCount: testkit.Ptr[int64](2),
					Boolean: []api.AggregateBooleanResult{
						{
							Property:        "onSale",
							Type:            testkit.Ptr("black_friday"),
							PercentageTrue:  testkit.Ptr[float64](1),
							PercentageFalse: testkit.Ptr[float64](2),
						},
						{
							Property:   "newArrival",
							Count:      testkit.Ptr[int64](1),
							TotalTrue:  testkit.Ptr[int64](2),
							TotalFalse: testkit.Ptr[int64](3),
						},
					},
				},
			},
		},
		{
			name: "date properties",
			reply: &proto.AggregateReply{
				Took: 92,
				Result: reply(Aggregations{
					{Property: "lastPurchase", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Date_{
						Date: &proto.AggregateReply_Aggregations_Aggregation_Date{
							Count:   testkit.Ptr[int64](1),
							Minimum: testkit.Ptr(testkit.Now.Format(time.RFC3339Nano)),
							Maximum: testkit.Ptr(testkit.Now.Format(time.RFC3339Nano)),
						},
					}},
					{Property: "lastReturn", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Date_{
						Date: &proto.AggregateReply_Aggregations_Aggregation_Date{
							Mode:   testkit.Ptr(testkit.Now.Format(time.RFC3339Nano)),
							Median: testkit.Ptr(testkit.Now.Format(time.RFC3339Nano)),
						},
					}},
				}),
			},
			dest: new(api.AggregateResponse),
			want: &api.AggregateResponse{
				Took: 92 * time.Second,
				Results: api.Aggregations{
					TotalCount: testkit.Ptr[int64](2),
					Date: []api.AggregateDateResult{
						{
							Property: "lastPurchase",
							Count:    testkit.Ptr[int64](1),
							Min:      &testkit.Now,
							Max:      &testkit.Now,
						},
						{
							Property: "lastReturn",
							Mode:     &testkit.Now,
							Median:   &testkit.Now,
						},
					},
				},
			},
		},
		{
			name: "grouped result",
			reply: &proto.AggregateReply{
				Took: 92,
				Result: &proto.AggregateReply_GroupedResults{
					GroupedResults: &proto.AggregateReply_Grouped{
						Groups: []*proto.AggregateReply_Group{
							{
								ObjectsCount: testkit.Ptr(int64(1)),
								GroupedBy: &proto.AggregateReply_Group_GroupedBy{
									Path:  []string{"onSale"},
									Value: &proto.AggregateReply_Group_GroupedBy_Boolean{Boolean: true},
								},
								Aggregations: &proto.AggregateReply_Aggregations{
									Aggregations: []*proto.AggregateReply_Aggregations_Aggregation{
										{Property: "onSale", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Boolean_{
											Boolean: &proto.AggregateReply_Aggregations_Aggregation_Boolean{
												Type:            testkit.Ptr("black_friday"),
												PercentageTrue:  testkit.Ptr[float64](1),
												PercentageFalse: testkit.Ptr[float64](2),
											},
										}},
									},
								},
							},
							{
								ObjectsCount: testkit.Ptr(int64(1)),
								GroupedBy: &proto.AggregateReply_Group_GroupedBy{
									Path:  []string{"price"},
									Value: &proto.AggregateReply_Group_GroupedBy_Number{Number: 4},
								},
								Aggregations: &proto.AggregateReply_Aggregations{
									Aggregations: []*proto.AggregateReply_Aggregations_Aggregation{
										{Property: "price", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Number_{
											Number: &proto.AggregateReply_Aggregations_Aggregation_Number{
												Sum:     testkit.Ptr[float64](1),
												Minimum: testkit.Ptr[float64](2),
												Maximum: testkit.Ptr[float64](3),
											},
										}},
									},
								},
							},
						},
					},
				},
			},
			dest: new(api.AggregateResponse),
			want: &api.AggregateResponse{
				Took: 92 * time.Second,
				GroupByResults: []api.AggregateGroup{
					{
						Property: "onSale",
						Value:    true,
						Results: api.Aggregations{
							TotalCount: testkit.Ptr[int64](1),
							Boolean: []api.AggregateBooleanResult{
								{
									Property:        "onSale",
									Type:            testkit.Ptr("black_friday"),
									PercentageTrue:  testkit.Ptr[float64](1),
									PercentageFalse: testkit.Ptr[float64](2),
								},
							},
						},
					},
					{
						Property: "price",
						Value:    float64(4),
						Results: api.Aggregations{
							TotalCount: testkit.Ptr[int64](1),
							Number: []api.AggregateNumberResult{
								{
									Property: "price",
									Sum:      testkit.Ptr[float64](1),
									Min:      testkit.Ptr[float64](2),
									Max:      testkit.Ptr[float64](3),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "object count",
			reply: &proto.AggregateReply{
				Result: &proto.AggregateReply_SingleResult{
					SingleResult: &proto.AggregateReply_Single{
						ObjectsCount: testkit.Ptr[int64](92),
					},
				},
			},
			dest: new(api.CountObjectsResponse),
			want: testkit.Ptr[api.CountObjectsResponse](92),
		},
	})
}

func TestInsertObjectBatchResponse_UnmarshalMessage(t *testing.T) {
	testMessageUnmarshaler(t, []MessageUnmarshalerTest[proto.BatchObjectsReply]{
		{
			name: "has errors",
			reply: &proto.BatchObjectsReply{
				Took: 92,
				Errors: []*proto.BatchObjectsReply_BatchError{
					{Index: 6, Error: "Whaam!"},
					{Index: 22, Error: "Whoops!"},
				},
			},
			dest: new(api.InsertObjectsResponse),
			want: &api.InsertObjectsResponse{
				Took:      92 * time.Second,
				Positions: []int32{6, 22},
				Errors:    []string{"Whaam!", "Whoops!"},
			},
		},
		{
			name:  "no errors",
			reply: &proto.BatchObjectsReply{Took: 92},
			dest:  new(api.InsertObjectsResponse),
			want:  &api.InsertObjectsResponse{Took: 92 * time.Second},
		},
	})
}
