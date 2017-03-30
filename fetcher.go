package webmedia

import (
	"net/http"
	"net/url"
)

type URLFetcher interface {
	FetchURL(u url.URL) (*http.Response, error)
}

type HTTPFetcher struct{}

func (f HTTPFetcher) FetchURL(u url.URL) (*http.Response, error) {
	return http.Get(u.String())
}
