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

	orderBy string

	tags []string

	fields []string

	publishedSince time.Time
	publishedUntil time.Time
}

func (v VideosQuery) endpoint() string {
	return "videos/with_pagination.json"
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

	if len(v.fields) != 0 {
		params.Set("only", strings.Join(v.fields, "|"))
	}

	if !v.publishedSince.IsZero() {
		params.Set("published_at.gte", v.publishedSince.Format(dateLayout))
	}

	if !v.publishedUntil.IsZero() {
		params.Set("published_at.lte", v.publishedUntil.Format(dateLayout))
	}

	if v.orderBy != "" {
		params.Set("order_by", v.orderBy)
	}

	return params
}

func (v VideosQuery) clone() VideosQuery {
	return VideosQuery{
		client:         v.client,
		perPage:        v.perPage,
		tags:           v.tags,
		fields:         v.fields,
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

func (v VideosQuery) OrderBy(order string) VideosQuery {
	clone := v.clone()
	clone.orderBy = order

	return clone
}

func (v VideosQuery) WithTags(tags ...string) VideosQuery {
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

func (v VideosQuery) Fields(fields ...string) VideosQuery {
	clone := v.clone()
	clone.fields = append(clone.fields, fields...)

	return clone
}

func (v VideosQuery) Pager() VideosPagerQuery {
	return VideosPagerQuery{v}
}

func (v VideosQuery) Fetch() (VideoResults, error) {
	var results VideoResults

	resp, err := v.client.fetch(v)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&results)
	return results, err
}

type VideoQuery struct {
	client *Client

	id int

	fields []string
}

func (v VideoQuery) endpoint() string {
	return fmt.Sprintf("videos/%d.json", v.id)
}

func (v VideoQuery) params() *url.Values {
	params := &url.Values{}

	if len(v.fields) != 0 {
		params.Set("only", strings.Join(v.fields, "|"))
	}

	return params
}

func (v VideoQuery) clone() VideoQuery {
	return VideoQuery{
		client: v.client,
		id:     v.id,
	}
}

func (v VideoQuery) Fields(fields ...string) VideoQuery {
	clone := v.clone()
	clone.fields = append(clone.fields, fields...)

	return clone
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

type VideoMetadata struct {
	ContentRating string `json:"content_rating"`
}

type Video struct {
	ID               int            `json:"id"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	Duration         int            `json:"duration"`
	PublishedAt      time.Time      `json:"published_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	ExhibitedAt      time.Time      `json:"exhibited_at"`
	SubscriberOnly   bool           `json:"subscriber_only"`
	Tags             []string       `json:"tags"`
	ExtendedMetadata *VideoMetadata `json:"extended_metadata"`
}

type VideoResults struct {
	Pager  Pager   `json:"pager"`
	Videos []Video `json:"videos"`
}
