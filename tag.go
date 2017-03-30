package webmedia

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type TagsQuery struct {
	client *Client

	name string
}

func (t TagsQuery) endpoint() string {
	return "tags.json"
}

func (t TagsQuery) params() *url.Values {
	params := &url.Values{}

	if t.name != "" {
		params.Set("name", t.name)
	}

	return params
}

func (t TagsQuery) clone() TagsQuery {
	return TagsQuery{
		name: t.name,
	}
}

func (t TagsQuery) Name(name string) TagsQuery {
	clone := t.clone()

	clone.name = name

	return clone
}

func (t TagsQuery) Fetch() ([]Tag, error) {
	var tags []Tag

	resp, err := t.client.fetch(t)
	if err != nil {
		return tags, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tags)
	return tags, err
}

type TagQuery struct {
	client *Client

	id int
}

func (t TagQuery) endpoint() string {
	return fmt.Sprintf("tags/%d.json", t.id)
}

func (t TagQuery) params() *url.Values {
	return nil
}

func (t TagQuery) Fetch() (*Tag, error) {
	resp, err := t.client.fetch(t)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tag Tag
	err = json.NewDecoder(resp.Body).Decode(&tag)
	return &tag, err
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
