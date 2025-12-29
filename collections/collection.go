package collections

type (
	Collection struct {
		Name          string
		Description   string
		Properties    []Property
		Sharding      *ShardingConfig
		InvertedIndex *InvertedIndexConfig
	}
	Property struct {
		Name              string
		Description       string
		DataType          DataType
		IndexFilterable   bool
		IndexRangeFilters bool
		IndexSearchable   bool
	}
	ShardingConfig struct {
		VirtualPerPhysical  int
		DesiredCount        int
		DesiredVirtualCount int
	}
	InvertedIndexConfig struct {
		IndexNullState         bool  // Index each object with the null state.
		IndexPropertyLength    bool  // Index length of properties.
		IndexTimestamps        bool  // Index each object by its internal timestamps.
		UsingBlockMaxWAND      bool  // Toggle UsingBlockMaxWAND usage for BM25 search.
		CleanupIntervalSeconds int32 // Asynchronous index cleanup internal.
	}
)

// DataType defines supported property data types.
type DataType string

const (
	DataTypeText DataType = "text"
	DataTypeBool DataType = "bool"
	DataTypeInt  DataType = "int"
)
