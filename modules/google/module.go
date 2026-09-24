package google

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
}

// Text2Vec is a vectorizer for text properties based on the text2vec-google module.
// Unset fields inherit the server defaults.
//
// The server targets Vertex AI by default, which requires ProjectID.
// To use AI Studio (Gemini API), set APIEndpoint to "generativelanguage.googleapis.com"; ProjectID
// is then not required.
//
// See https://docs.weaviate.io/weaviate/model-providers/google/embeddings.
type Text2Vec struct {
	// APIEndpoint is a host name without a scheme such as "https://".
	// Defaults to us-central1-aiplatform.googleapis.com (Vertex AI).
	APIEndpoint string `json:"apiEndpoint,omitempty"`
	// ProjectID is the Google Cloud project ID; required for Vertex AI.
	ProjectID string `json:"projectId,omitempty"`
	// Model defaults to gemini-embedding-001.
	// An invalid name is not rejected at collection create; it fails when objects are vectorized.
	Model string `json:"model,omitempty"`
	// Location is the Vertex AI region to run the model in. Defaults to us-central1.
	Location string `json:"location,omitempty"`
	// Dimensions is the size of the output vectors. Defaults to 768 for gemini-embedding-001.
	Dimensions int `json:"dimensions,omitzero"`
	// TaskType defaults to [TaskTypeRetrievalQuery].
	TaskType TaskType `json:"taskType,omitempty"`
	// TitleProperty names the property the model uses as the document title.
	TitleProperty string `json:"titleProperty,omitempty"`
	// Properties limits vectorization to these properties.
	// By default, all text properties are vectorized.
	Properties []string `json:"properties,omitempty"`
}

func (Text2Vec) Name() string { return "text2vec-google" }

// TaskType tells the model what the embeddings will be used for.
type TaskType string

const (
	TaskTypeRetrievalQuery     TaskType = "RETRIEVAL_QUERY"
	TaskTypeQuestionAnswering  TaskType = "QUESTION_ANSWERING"
	TaskTypeFactVerification   TaskType = "FACT_VERIFICATION"
	TaskTypeCodeRetrievalQuery TaskType = "CODE_RETRIEVAL_QUERY"
	TaskTypeClassification     TaskType = "CLASSIFICATION"
	TaskTypeClustering         TaskType = "CLUSTERING"
	TaskTypeSemanticSimilarity TaskType = "SEMANTIC_SIMILARITY"
)
