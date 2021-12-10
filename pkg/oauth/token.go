package oauth

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/bachelor-thesis-hown3d/chat-api-server/pkg/oauth"
	"github.com/bachelor-thesis-hown3d/chat-client/pkg/errors"
	"github.com/bachelor-thesis-hown3d/chat-client/pkg/util"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/metadata"
	"gopkg.in/yaml.v2"
)

const (
	tokenKey string = "token"
)

type Config struct {
	OAuth2Config *oauth2.Config
	OIDCConfig   *oauth.Config
}

func NewConfig(ctx context.Context, issuerURL, redirectURL *url.URL, clientID, clientSecret string) (*Config, error) {
	oidcConf, err := oauth.NewConfig(ctx, issuerURL, clientID)
	if err != nil {
		return nil, err
	}
	oauth2conf := &oauth2.Config{
		RedirectURL:  redirectURL.String(),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{oidc.ScopeOpenID},
		Endpoint:     oidcConf.Provider.Endpoint(),
	}

	return &Config{
		OAuth2Config: oauth2conf,
		OIDCConfig:   oidcConf,
	}, nil
}

type Token struct {
	IDToken     string
	OAuth2Token *oauth2.Token
}

func NeedTokenRefresh(t Token) (bool, error) {
	return t.OAuth2Token.Expiry.Before(time.Now()), nil
}

// LoadTokenFromFile loads the token.
// Error is only returned if the token file exists but could not be loaded.
func LoadTokenFromFile() (Token, error) {
	file := util.TokenFile()
	if _, err := os.Stat(file); err != nil {
		// token file does not exist
		return Token{}, errors.TokenFileNotFound("test")
	}

	b, err := os.ReadFile(file)
	if err != nil {
		return Token{}, fmt.Errorf("could not load previous token: %w", err)
	}

	var t Token
	err = yaml.Unmarshal(b, &t)
	if err != nil {
		return Token{}, fmt.Errorf("could not load previous token: %w", err)
	}
	return t, nil
}

// SafeTokenToFile will persist a token on filesystem
func SafeTokenToFile(t Token) error {
	fmt.Println("Saving OIDC Token to file")
	return util.WriteYAML(t, util.TokenFile())
}

// LoadTokenIntoContext creates metadata for authentication to the apiServer with the token from filesystem
// will return an error, if the Token could not be loaded
func LoadTokenIntoContext(ctx context.Context) (context.Context, error) {
	t, err := LoadTokenFromFile()
	if err != nil {
		return ctx, err
	}

	md := metadata.Pairs("authorization", "bearer "+t.IDToken)
	return metadata.NewOutgoingContext(ctx, md), nil
}

// RefreshToken refreshes the token for a given Token
func (c *Config) RefreshToken(t Token) (Token, error) {
	tokenSource := c.OAuth2Config.TokenSource(oauth2.NoContext, t.OAuth2Token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return t, err
	}

	newIDToken, ok := newToken.Extra("id_token").(string)
	if !ok {
		return t, fmt.Errorf("id_token was missing from oauth2 response")
	}
	if newIDToken == t.IDToken {
		return t, nil
	}
	return Token{IDToken: newIDToken, OAuth2Token: newToken}, nil
}

// retrieveTokenFromOIDCIssuer returns the id_token from the oidc issuer
func (c *Config) retrieveTokenFromOIDCIssuer(curState, serverState, code string) (Token, error) {
	if curState != serverState {
		return Token{}, fmt.Errorf("invalid oauth state")
	}

	// since we use a bad ssl certificate, embed a Insecure HTTP Client for oauth to use
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, http.DefaultClient)

	token, err := c.OAuth2Config.Exchange(ctx, code)
	if err != nil {
		return Token{}, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return Token{}, fmt.Errorf("id_token was missing from oauth2 response")
	}
	_, err = c.OIDCConfig.Verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		return Token{}, fmt.Errorf("Can't verify idToken: %v", err)
	}
	return Token{IDToken: rawIDToken, OAuth2Token: token}, nil
}
