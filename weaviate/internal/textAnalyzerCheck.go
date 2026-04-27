package internal

import (
	"strconv"
	"strings"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// CheckTextAnalyzerSupport returns a WeaviateClientError when class uses
// Weaviate 1.37.0 schema features (per-property TextAnalyzer or collection
// StopwordPresets) against an older server.
func CheckTextAnalyzerSupport(dbVersionProvider *db.VersionProvider, class *models.Class) error {
	if class == nil || dbVersionProvider == nil {
		return nil
	}
	feature := detectTextAnalyzerFeature(class)
	if feature == "" {
		return nil
	}
	version := dbVersionProvider.Version()
	if supportsTextAnalyzer(version) {
		return nil
	}
	return except.NewWeaviateClientErrorf(0,
		"%s requires Weaviate >= 1.37.0 (connected server: %q)", feature, version)
}

func supportsTextAnalyzer(version string) bool {
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}
	minor, err := strconv.Atoi(strings.SplitN(parts[1], "-", 2)[0])
	if err != nil {
		return false
	}
	return major > 1 || (major == 1 && minor >= 37)
}

func detectTextAnalyzerFeature(class *models.Class) string {
	if class.InvertedIndexConfig != nil && len(class.InvertedIndexConfig.StopwordPresets) > 0 {
		return "InvertedIndexConfig.stopwordPresets"
	}
	for _, p := range class.Properties {
		if propertyUsesTextAnalyzer(p) {
			return "Property.textAnalyzer"
		}
	}
	return ""
}

func propertyUsesTextAnalyzer(p *models.Property) bool {
	if p == nil {
		return false
	}
	if p.TextAnalyzer != nil {
		return true
	}
	for _, np := range p.NestedProperties {
		if nestedPropertyUsesTextAnalyzer(np) {
			return true
		}
	}
	return false
}

func nestedPropertyUsesTextAnalyzer(p *models.NestedProperty) bool {
	if p == nil {
		return false
	}
	if p.TextAnalyzer != nil {
		return true
	}
	for _, np := range p.NestedProperties {
		if nestedPropertyUsesTextAnalyzer(np) {
			return true
		}
	}
	return false
}
