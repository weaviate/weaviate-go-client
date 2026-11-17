package api

import (
	"encoding/json"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

type TokenizeRequest struct {
	transports.BaseEndpoint
	RequestDefaults

	PropetyName     string
	Text            string
	Tokenization    Tokenization
	TextAnalyzer    *TextAnalyzerConfig
	Stopwords       *StopwordConfig
	StopwordPresets map[string][]string
}

var _ transports.Endpoint = (*TokenizeRequest)(nil)

func (*TokenizeRequest) Method() string { return http.MethodPost }
func (r *TokenizeRequest) Path() string {
	if r.CollectionName != "" && r.PropetyName != "" {
		return "/schema/" + r.CollectionName + "/properties/" + r.PropetyName + "/tokenize"
	}
	return "/tokenize"
}

func (r *TokenizeRequest) Body() any {
	if r.CollectionName != "" {
		return rest.TokenizeRequest{
			Text: r.Text,
		}
	}

	req := rest.TokenizeRequest{
		Text:            r.Text,
		Tokenization:    rest.TokenizeRequestTokenization(r.Tokenization),
		StopwordPresets: r.StopwordPresets,
	}
	if t := r.TextAnalyzer; t != nil {
		req.AnalyzerConfig = rest.TextAnalyzerConfig{
			AsciiFold:       t.FoldASCII,
			AsciiFoldIgnore: t.KeepASCII,
			StopwordPreset:  t.StopwordPreset,
		}
	}
	if r.Stopwords != nil {
		req.Stopwords = rest.StopwordConfig(*r.Stopwords)
	}
	return req
}

type TokenizeResponse rest.TokenizeResponse

var _ json.Unmarshaler = (*TokenizeResponse)(nil)

func (r *TokenizeResponse) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*rest.TokenizeResponse)(r))
}
