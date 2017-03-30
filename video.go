package webmedia

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type VideosQuery struct {
	client *Client

	page    int
	perPage int

	tags []string

	publishedSince time.Time
	publishedUntil time.Time
}

func (v VideosQuery) endpoint() string {
	return "videos.json"
}

func (v VideosQuery) params() *url.Values {
	params := &url.Values{}

	if v.page != 0 {
		params.Set("page", strconv.Itoa(v.page))
	}

	if v.perPage != 0 {
		params.Set("per_page", strconv.Itoa(v.perPage))
	}

	if len(v.tags) != 0 {
		params.Set("tags.all", strings.Join(v.tags, "|"))
	}

	if !v.publishedSince.IsZero() {
		params.Set("published_at.gte", v.publishedSince.Format(dateLayout))
	}

	if !v.publishedUntil.IsZero() {
		params.Set("published_at.lte", v.publishedUntil.Format(dateLayout))
	}

	return params
}

func (v VideosQuery) clone() VideosQuery {
	return VideosQuery{
		client:         v.client,
		perPage:        v.perPage,
		tags:           v.tags,
		publishedSince: v.publishedSince,
		publishedUntil: v.publishedUntil,
	}
}

func (v VideosQuery) PerPage(n int) VideosQuery {
	clone := v.clone()
	clone.perPage = n

	return clone
}

func (v VideosQuery) Page(n int) VideosQuery {
	clone := v.clone()
	clone.page = n

	return clone
}

func (v VideosQuery) AddTags(tags ...string) VideosQuery {
	clone := v.clone()
	clone.tags = append(clone.tags, tags...)

	return clone
}

func (v VideosQuery) PublishedSince(t time.Time) VideosQuery {
	clone := v.clone()
	clone.publishedSince = t

	return clone
}

func (v VideosQuery) PublishedUntil(t time.Time) VideosQuery {
	clone := v.clone()
	clone.publishedUntil = t

	return clone
}

func (v VideosQuery) Fetch() ([]Video, error) {
	var videos []Video

	resp, err := v.client.fetch(v)
	if err != nil {
		return videos, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&videos)
	return videos, err
}

type VideoQuery struct {
	client *Client

	ID int
}

func (v VideoQuery) endpoint() string {
	return fmt.Sprintf("videos/%d.json", v.ID)
}

func (v VideoQuery) params() *url.Values {
	return nil
}

func (v VideoQuery) Fetch() (*Video, error) {
	resp, err := v.client.fetch(v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var video Video
	err = json.NewDecoder(resp.Body).Decode(&video)
	return &video, err
}

type Video struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	PublishedAt time.Time `json:"published_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ExhibitedAt time.Time `json:"exhibited_at"`
	Tags        []string  `json:"tags"`
}
