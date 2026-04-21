// Package tokenize provides access to Weaviate's /v1/tokenize REST
// endpoints (available in Weaviate >= 1.37.0). It lets clients inspect how
// text is tokenized — either with an arbitrary tokenization method or using
// the configuration of an existing property.
package tokenize

// Tokenization identifies a tokenization method supported by Weaviate.
type Tokenization string

const (
	Word       Tokenization = "word"
	Lowercase  Tokenization = "lowercase"
	Whitespace Tokenization = "whitespace"
	Field      Tokenization = "field"
	Trigram    Tokenization = "trigram"
	Gse        Tokenization = "gse"
	GseCh      Tokenization = "gse_ch"
	KagomeJa   Tokenization = "kagome_ja"
	KagomeKr   Tokenization = "kagome_kr"
)

// AnalyzerConfig configures optional text analyzer behaviour that is applied
// before tokenization (ASCII folding and stopword preset selection).
type AnalyzerConfig struct {
	// AsciiFold, when true, folds accented characters to their ASCII form
	// (e.g. "école" → "ecole").
	AsciiFold *bool `json:"asciiFold,omitempty"`
	// AsciiFoldIgnore lists characters that should be excluded from ASCII
	// folding. Requires AsciiFold to be true.
	AsciiFoldIgnore []string `json:"asciiFoldIgnore,omitempty"`
	// StopwordPreset is the name of a stopword preset to apply. This can be a
	// built-in preset (e.g. "en", "none") or the name of a custom preset
	// provided via TextTokenizer.WithStopwordPresets.
	StopwordPreset string `json:"stopwordPreset,omitempty"`
}

// StopwordConfig is a custom stopword preset definition.
type StopwordConfig struct {
	// Preset is the name of a built-in base preset to extend (e.g. "en" or
	// "none"). If empty, no base preset is used.
	Preset string `json:"preset,omitempty"`
	// Additions are extra stopwords that the preset should include.
	Additions []string `json:"additions,omitempty"`
	// Removals are stopwords that should be removed from the base preset.
	Removals []string `json:"removals,omitempty"`
}

// TokenizeResult is the response of a tokenize request.
type TokenizeResult struct {
	// Tokenization is the method that was applied.
	Tokenization Tokenization `json:"tokenization"`
	// Indexed contains the tokens as they are stored in the inverted index.
	Indexed []string `json:"indexed"`
	// Query contains the tokens as they are used for query matching (after
	// stopword removal, if applicable).
	Query []string `json:"query"`
	// AnalyzerConfig echoes the analyzer configuration that was applied, if
	// any.
	AnalyzerConfig *AnalyzerConfig `json:"analyzerConfig,omitempty"`
	// StopwordConfig echoes the resolved stopword configuration, if any.
	StopwordConfig *StopwordConfig `json:"stopwordConfig,omitempty"`
}

// tokenizeRequest is the on-the-wire payload for POST /v1/tokenize.
type tokenizeRequest struct {
	Text            string                     `json:"text"`
	Tokenization    Tokenization               `json:"tokenization"`
	AnalyzerConfig  *AnalyzerConfig            `json:"analyzerConfig,omitempty"`
	StopwordPresets map[string]*StopwordConfig `json:"stopwordPresets,omitempty"`
}
