package paragons
import (
	weaviateModels "github.com/semi-technologies/weaviate/entities/models"
)

// SchemaDump Contains all semantic types and respective classes of the schema
type SchemaDump struct{
	Things *weaviateModels.Schema ` json:"things"`
	Actions *weaviateModels.Schema `json:"actions"`
}
