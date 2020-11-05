package paragons

// SemanticKind defining if a weaviate object is a thing or an action
type SemanticKind string

// SemanticKindThings usually indicated by substantives
const SemanticKindThings SemanticKind = "things"

// SemanticKindActions usually indicated by verbs
const SemanticKindActions SemanticKind = "actions"
