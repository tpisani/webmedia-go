package webmedia

import (
	"net/url"
)

type ClientOption func(c *Client)

func WithBaseURL(u url.URL) ClientOption {
	return func(c *Client) {
		c.url = u
	}
}

func WithTransport(transport Transport) ClientOption {
	return func(c *Client) {
		c.transport = transport
	}
}
