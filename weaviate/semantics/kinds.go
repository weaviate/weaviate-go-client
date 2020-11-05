package semantics

// Kind defining if a weaviate object is a thing or an action
type Kind string

// Things usually indicated by substantives
const Things Kind = "things"

// Actions usually indicated by verbs
const Actions Kind = "actions"
