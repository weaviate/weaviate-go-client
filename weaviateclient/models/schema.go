package models
import (
	weaviateModels "github.com/semi-technologies/weaviate/entities/models"
)

type SchemaDump struct{
	Things *weaviateModels.Schema ` json:"things"`
	Actions *weaviateModels.Schema `json:"actions"`
}
