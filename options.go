package webmedia

import (
	"net/http"
	"net/url"
)

type ClientOption func(c *Client)

func WithBaseURL(u url.URL) ClientOption {
	return func(c *Client) {
		c.url = u
	}
}

func WithRoundTripper(roundTripper http.RoundTripper) ClientOption {
	return func(c *Client) {
		c.roundTripper = roundTripper
	}
}
