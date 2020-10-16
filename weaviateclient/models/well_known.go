package models

// OpenIDConfiguration of weaviate
type OpenIDConfiguration struct {
	// The Location to redirect to
	Href string `json:"href,omitempty"`
	// OAuth Client ID
	ClientID string `json:"clientId,omitempty"`
}