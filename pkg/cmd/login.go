package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bachelor-thesis-hown3d/chat-client/pkg/errors"
	"github.com/bachelor-thesis-hown3d/chat-client/pkg/oauth"
	"github.com/pkg/browser"
)

func Login(ctx context.Context, clientID, clientSecret string, issuerUrl string) error {
	// parse the URL strings to issuer format
	issuer, err := url.Parse(issuerUrl)
	if err != nil {
		return err
	}

	// parse the redirect URL for the port number
	redirectBase, err := url.Parse("http://localhost:7070")
	if err != nil {
		return err
	}
	redirectCallback, err := redirectBase.Parse(oauth.CallbackPath)
	if err != nil {
		return err
	}

	// oauth config
	c, err := oauth.NewConfig(ctx, issuer, redirectCallback, clientID, clientSecret)
	if err != nil {
		return err
	}

	// check if we already logged in, so we can just refresh the token
	// get the token from the file
	t, err := oauth.LoadTokenFromFile()

	if err != nil {
		// error is not of type tokenfile not found, so handle it
		if _, ok := err.(errors.TokenFileNotFound); !ok {
			return err
		}
		return getNewToken(ctx, c, redirectBase)
	}
	// error is not nil, so we could load the token
	// checking if we need a refresh
	return useExistingToken(c, t)

}

func getNewToken(ctx context.Context, c *oauth.Config, redirect *url.URL) error {
	// create oauth server
	s, err := oauth.NewServer(ctx, c, redirect)
	if err != nil {
		return err
	}

	redirectLogin, err := redirect.Parse(oauth.LoginPath)
	if err != nil {
		return err
	}
	errChan := make(chan error, 1)
	s.Start(errChan)
	browser.OpenURL(redirectLogin.String())
	// stop the http server when done retrieving token
	defer s.Stop()

	for {
		select {
		case err := <-errChan:
			return err
		case token := <-s.TokenChan:
			return oauth.SafeTokenToFile(token)
		}
	}
}

func useExistingToken(c *oauth.Config, t oauth.Token) error {
	fmt.Println("Using existing token...")
	refresh, err := oauth.NeedTokenRefresh(t)
	if err != nil {
		return err
	}
	if refresh {
		fmt.Println("Token needs a refresh...")
		newToken, err := c.RefreshToken(t)
		if err != nil {
			return err
		}
		return oauth.SafeTokenToFile(newToken)
	}
	return nil
}
