package webmedia

import (
	"net/http"
	"net/url"
)

type query interface {
	endpoint() string
	params() *url.Values
}

type Transport interface {
	Do(r *http.Request) (*http.Response, error)
}

type Client struct {
	url         url.URL
	accessToken string
	transport   Transport
}

func NewClient(accessToken string, options ...ClientOption) *Client {
	c := &Client{
		url: url.URL{
			Scheme: "https",
			Host:   "api.video.globoi.com",
		},
		accessToken: accessToken,
		transport:   http.DefaultClient,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func (c *Client) baseURL() url.URL {
	return c.url
}

func (c *Client) buildURL(endpoint string, params *url.Values) url.URL {
	if params == nil {
		params = &url.Values{}
	}

	params.Set("access_token", c.accessToken)

	u := c.baseURL()
	u.Path = endpoint
	u.RawQuery = params.Encode()

	return u
}

func (c *Client) fetch(q query) (*http.Response, error) {
	u := c.buildURL(q.endpoint(), q.params())
	r, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.transport.Do(r)
}

func (c *Client) Video(id int) VideoQuery {
	return VideoQuery{
		client: c,
		id:     id,
	}
}

func (c *Client) Videos() VideosQuery {
	return VideosQuery{
		client: c,
	}
}

func (c *Client) Tag(id int) TagQuery {
	return TagQuery{
		client: c,
		id:     id,
	}
}

func (c *Client) Tags() TagsQuery {
	return TagsQuery{
		client: c,
	}
}
