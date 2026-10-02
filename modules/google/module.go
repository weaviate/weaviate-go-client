package google

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
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
	// ProjectID is the Google Cloud project ID, required for Vertex AI.
	ProjectID string `json:"projectId,omitempty"`
	// Location is the Vertex AI region to run the model in.
	Location string `json:"location,omitempty"`
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
