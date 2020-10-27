package connection

import (
	clientModels "github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponseData_DecodeBodyIntoTarget(t *testing.T) {

	respond := ResponseData{
		Body:       []byte(`{"actions":{"classes":[],"type":"action"},"things":{"classes":[{"class":"Band","description":"Band that plays and produces music","properties":null}],"type":"thing"}}`),
		StatusCode: 0,
	}

	var schema clientModels.SchemaDump
	err := respond.DecodeBodyIntoTarget(&schema)
	assert.Nil(t, err)
	assert.Equal(t, "Band", schema.Things.Classes[0].Class)

}
