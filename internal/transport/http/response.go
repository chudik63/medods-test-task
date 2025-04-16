package http

// swagger:model ErrorResponse
type ErrorResponse struct {
	// Error message
	Error string `json:"error"`
}

// swagger:model TokenResponse
type TokenResponse struct {
	// Access token
	AccessToken string `json:"access_token"`

	// Refresh token
	RefreshToken string `json:"refresh_token"`
}
