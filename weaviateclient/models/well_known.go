package models

type OpenIDConfiguration struct {
	// The Location to redirect to
	Href string `json:"href,omitempty"`
	// OAuth Client ID
	ClientId string `json:"clientId,omitempty"`
}