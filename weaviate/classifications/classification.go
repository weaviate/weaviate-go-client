package classifications

import "github.com/weaviate/weaviate-go-client/v5/weaviate/connection"

// API classifications API
type API struct {
	connection *connection.Connection
}

// New Classification api group from connection
func New(con *connection.Connection) *API {
	return &API{connection: con}
}

// Scheduler get a builder to schedule a classification
func (api *API) Scheduler() *Scheduler {
	return &Scheduler{connection: api.connection}
}

// Getter get a builder to retrieve a classification
func (api *API) Getter() *Getter {
	return &Getter{connection: api.connection}
}

// KNN (k nearest neighbours) a non parametric classification based on training data
const KNN = "knn"

// Contextual classification labels a data object with
// the closest label based on their vector position (which describes the context).
// It can only be used with the text2vec-contextionary vectorizer.
const Contextual = "text2vec-contextionary"

// ZeroShot classification labels a data object with
// the closest label based on their vector position (which describes the context)
// It can be used with any vectorizer or custom vectors.
const ZeroShot = "zeroshot"
