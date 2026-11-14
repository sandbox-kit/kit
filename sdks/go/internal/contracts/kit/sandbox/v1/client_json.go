package sandboxv1

import "encoding/json"

// MarshalJSON prevents the API key from entering ordinary Go JSON output.
func (APIKeyCredentials) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}

// MarshalJSON prevents the token identifier and secret from entering ordinary
// Go JSON output.
func (TokenPairCredentials) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}

// MarshalJSON prevents the bearer token from entering ordinary Go JSON output.
func (BearerTokenCredentials) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}

// MarshalJSON retains the non-sensitive OAuth client identifier and prevents
// the refresh token, client secret, and JWT key from entering ordinary Go JSON
// output.
func (credentials OAuthCredentials) MarshalJSON() ([]byte, error) {
	type publicOAuthCredentials struct {
		ClientID string `json:"client_id,omitempty"`
	}
	return json.Marshal(publicOAuthCredentials{ClientID: credentials.ClientId})
}
