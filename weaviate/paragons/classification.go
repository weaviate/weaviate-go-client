package paragons

// Classification defines the type of classification
type Classification string

// KNN (k nearest neighbours) a non parametric classification based on training data
const KNN Classification = "knn"

// Contextual classification labels a data object with the closest label based on their vector position (which describes the context)
const Contextual Classification = "contextual"
