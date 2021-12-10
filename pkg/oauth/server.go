package oauth

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
)

const (
	LoginPath    string = "/auth"
	CallbackPath string = "/callback"
)

func randomHex(n int) (string, error) {
	rand.Seed(42)
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *Server) handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	s.state, _ = randomHex(16)
	url := s.conf.OAuth2Config.AuthCodeURL(s.state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	token, err := s.conf.retrieveTokenFromOIDCIssuer(r.FormValue("state"), s.state, r.FormValue("code"))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	s.TokenChan <- token
}

type Server struct {
	TokenChan chan Token
	http      *http.Server
	listener  net.Listener
	conf      *Config
	state     string
}

func NewServer(ctx context.Context, c *Config, redirectUrl *url.URL) (*Server, error) {

	s := &http.Server{
		Addr: redirectUrl.Host,
	}
	// set up a listener on the redirect port
	port := fmt.Sprintf(":%v", redirectUrl.Port())
	l, err := net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("can't listen to port %v: %v", port, err)
	}

	serv := &Server{
		http:      s,
		listener:  l,
		TokenChan: make(chan Token),
		conf:      c,
	}

	mux := http.NewServeMux()

	mux.HandleFunc(LoginPath, serv.handleAuthLogin)
	mux.HandleFunc(CallbackPath, serv.handleAuthCallback)
	s.Handler = mux

	return serv, nil
}

// Start starts the execution of the oauth server, calling the serve function in a goroutine, so non blocking
func (s *Server) Start(errChan chan error) {
	fmt.Printf("Starting oauth listener on %v\n", s.listener.Addr())
	go func() {
		errChan <- s.http.Serve(s.listener)
	}()
}

// Stop stops the execution of the oauth server
func (s *Server) Stop() {
	s.http.Close()
}
