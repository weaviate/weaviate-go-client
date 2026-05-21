package internal

import (
	"strconv"
	"strings"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
	"github.com/weaviate/weaviate/entities/models"
)

// supportsServerSideDefaultVectorIndexType returns true iff the parsed
// version is >= 1.37.5. Empty or malformed -> false (treat as old server).
func supportsServerSideDefaultVectorIndexType(version string) bool {
	if version == "" {
		return false
	}
	// strip "-suffix" (e.g. "1.37.5-e0fe0d5.amd64")
	base := strings.SplitN(version, "-", 2)[0]
	parts := strings.Split(base, ".")
	if len(parts) < 3 {
		return false
	}
	major, err1 := strconv.Atoi(parts[0])
	minor, err2 := strconv.Atoi(parts[1])
	patch, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}
	if major > 1 {
		return true
	}
	if major < 1 {
		return false
	}
	if minor > 37 {
		return true
	}
	if minor < 37 {
		return false
	}
	return patch >= 5
}

// InjectLegacyDefaultVectorIndexType injects "hnsw" into any empty
// VectorIndexType field on the given class — but only when the server is
// older than 1.37.5. Servers >= 1.37.5 apply DEFAULT_VECTOR_INDEX themselves.
func InjectLegacyDefaultVectorIndexType(p *db.VersionProvider, class *models.Class) {
	if class == nil || p == nil {
		return
	}
	if supportsServerSideDefaultVectorIndexType(p.Version()) {
		return
	}
	if class.VectorIndexType == "" {
		class.VectorIndexType = "hnsw"
	}
	for k, vc := range class.VectorConfig {
		if vc.VectorIndexType == "" {
			vc.VectorIndexType = "hnsw"
			class.VectorConfig[k] = vc
		}
	}
}
