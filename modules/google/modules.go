package google

import (
	"time"

	"github.com/weaviate/weaviate-go-client/v6/modules"
	"github.com/weaviate/weaviate-go-client/v6/modules/internal"
)

func init() {
	modules.Register(*new(Text2Vec))
	modules.Register(*new(Multi2Vec))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-google module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - text2vec-google].
//
// [Weaviate Docs - text2vec-google]: https://docs.weaviate.io/weaviate/model-providers/google/embeddings
type Text2Vec struct {
	// APIEndpoint is a host name without a scheme such as "https://".
	APIEndpoint string `json:"apiEndpoint,omitempty"`
	// Model is the embedding model name.
	Model string `json:"model,omitempty"`
	// Dimensions is the size of the output vectors.
	Dimensions int `json:"dimensions,omitzero"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
	// VertexProjectID is the Google Cloud project ID, required for Vertex AI.
	VertexProjectID string `json:"projectId,omitempty"`
	// VertexLocation is the Vertex AI region to run the model in.
	VertexLocation string `json:"location,omitempty"`
	// TaskType tells the model what the embeddings will be used for.
	TaskType string `json:"taskType,omitempty"`
	// TitleProperty names the property the model uses as the document title.
	TitleProperty string `json:"titleProperty,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-google" }

const (
	RetrievalQuery     = "RETRIEVAL_QUERY"
	QuestionAnswering  = "QUESTION_ANSWERING"
	FactVerification   = "FACT_VERIFICATION"
	CodeRetrievalQuery = "CODE_RETRIEVAL_QUERY"
	Classification     = "CLASSIFICATION"
	Clustering         = "CLUSTERING"
	SemanticSimilarity = "SEMANTIC_SIMILARITY"
)

// Multi2Vec is a vectorizer for text, image, video and audio properties based on the multi2vec-google module.
// Unset fields inherit the server defaults.
//
// See [Weaviate Docs - multi2vec-google].
//
// [Weaviate Docs - multi2vec-google]: https://docs.weaviate.io/weaviate/model-providers/google/embeddings-multimodal
type Multi2Vec struct {
	// APIEndpoint is the host name.
	APIEndpoint string
	// Model is the embedding model name.
	Model string
	// Dimensions is the size of the output vectors.
	Dimensions int
	// VertexProjectID is the Google Cloud project ID, required for Vertex AI.
	VertexProjectID string
	// VertexLocation is the Vertex AI region to run the model in.
	VertexLocation string
	// VideoInterval is the length of the video segments to embed.
	VideoInterval time.Duration

	// TextFields, ImageFields, VideoFields and AudioFields name the properties to vectorize.
	TextFields  []modules.WeightedField `json:"-"`
	ImageFields []modules.WeightedField `json:"-"`
	AudioFields []modules.WeightedField `json:"-"`
	VideoFields []modules.WeightedField `json:"-"`
}

func (Multi2Vec) Name() string { return "multi2vec-google" }

type multi2vecJSON struct {
	APIEndpoint          string `json:"apiEndpoint,omitempty"`
	Model                string `json:"model,omitempty"`
	Dimensions           int    `json:"dimensions,omitzero"`
	VertexProjectID      string `json:"projectId,omitempty"`
	VertexLocation       string `json:"location,omitempty"`
	VideoIntervalSeconds int    `json:"videoIntervalSeconds,omitzero"`
}

func (m2v Multi2Vec) EncodeMap() (map[string]any, error) {
	return internal.EncodeWithWeightedFields(multi2vecJSON{
		APIEndpoint:          m2v.APIEndpoint,
		Model:                m2v.Model,
		Dimensions:           m2v.Dimensions,
		VertexProjectID:      m2v.VertexProjectID,
		VertexLocation:       m2v.VertexLocation,
		VideoIntervalSeconds: int(m2v.VideoInterval.Seconds()),
	}, internal.WeightedFields{
		TextFields:  m2v.TextFields,
		ImageFields: m2v.ImageFields,
		AudioFields: m2v.AudioFields,
		VideoFields: m2v.VideoFields,
	})
}

func (m2v *Multi2Vec) DecodeMap(m map[string]any) error {
	var dest multi2vecJSON
	wfs, err := internal.DecodeWithWeightedFields(m, &dest)
	if err != nil {
		return err
	}

	*m2v = Multi2Vec{
		APIEndpoint:     dest.APIEndpoint,
		Model:           dest.Model,
		Dimensions:      dest.Dimensions,
		VertexProjectID: dest.VertexProjectID,
		VertexLocation:  dest.VertexLocation,
		VideoInterval:   time.Duration(dest.VideoIntervalSeconds) * time.Second,
		TextFields:      wfs.TextFields,
		ImageFields:     wfs.ImageFields,
		AudioFields:     wfs.AudioFields,
		VideoFields:     wfs.VideoFields,
	}
	return nil
}
