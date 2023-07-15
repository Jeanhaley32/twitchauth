// Library for authenticating with Twitch using OAuth2 and the client credentials grant flow
// https://dev.twitch.tv/docs/authentication/getting-tokens-oauth#oauth-client-credentials-flow
package twitchauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Constants for the Twitch API
const (
	// Twitch API URL
	twitchAuthTokenURL = "https://id.twitch.tv/oauth2/token"
)

// TwitchAuth is the struct for the Twitch API
type TwitchAuth struct {
	ClientID string
	Secret   string
	Token    token
}

// token is the response from the Twitch API
type token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type TwitchAuthInterface interface {
	NewTokenSet()
	remainingTime()
}

// Returns time until token expires
func (self *TwitchAuth) remainingTime() int64 {
	return self.Token.ExpiresIn
}

// Obtains a new Token set from the Twitch API
// Token set includes access token, Type, expiration time
func (self *TwitchAuth) NewTokenSet() error {
	var t token
	// Client credentials grant flow
	// https://dev.twitch.tv/docs/authentication/getting-tokens-oauth#oauth-client-credentials-flow
	data := url.Values{}
	data.Set("client_id", self.ClientID)
	data.Set("client_secret", self.Secret)
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", twitchAuthTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		fmt.Println("error decoding json")
		return err
	}

	self.Token = t

	return nil
}
