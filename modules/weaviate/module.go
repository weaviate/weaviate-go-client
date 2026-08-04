package weaviate

import "github.com/weaviate/weaviate-go-client/v6/modules"

func init() {
	modules.Register(*new(Text2Vec))
}

// Text2Vec is a vectorizer for text properties
// based on the text2vec-weaviate module.
type Text2Vec struct {
	URL        string   `json:"baseURL"`
	Properties []string `json:"properties"`
	Model      string   `json:"model"`
	Dimensions int      `json:"dimensions"`
}

func (Text2Vec) Name() string { return "text2vec-weaviate" }

const (
	SnowflakeArcticEmbedMv1_5 = "Snowflake/snowflake-arctic-embed-m-v1.5"
	SnowflakeArcticEmbedLv2_0 = "Snowflake/snowflake-arctic-embed-l-v2.0"
)
