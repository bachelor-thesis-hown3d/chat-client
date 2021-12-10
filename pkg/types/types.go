package types

import (
	"errors"
	"log"
	"net/url"
)

type HttpURL url.URL

func NewHttpURL(defaultVal string, u *HttpURL) *HttpURL {
	d, err := url.Parse(defaultVal)
	if err != nil {
		log.Fatal(err)
	}
	defaultUri := (HttpURL)(*d)
	*u = defaultUri
	return u
}

func (p *HttpURL) Set(in string) error {
	u, err := url.Parse(in)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "http", "https":
	default:
		return errors.New("unexpected scheme in URL")
	}

	*p = HttpURL(*u)
	return nil
}

func (p HttpURL) String() string {
	return (*url.URL)(&p).String()
}

func (p HttpURL) Type() string {
	return "url"
}
