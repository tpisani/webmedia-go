package webmedia

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"testing"
)

type MockRoundTripper struct{}

func (m MockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	var mock io.Reader
	var err error

	if r.URL.Path == "/videos.json" {
		mock, err = os.Open("mocks/videos.json")
	} else if r.URL.Path == "/videos/5767587.json" {
		mock, err = os.Open("mocks/video-5767587.json")
	} else if r.URL.Path == "/videos/6053793.json" {
		mock, err = os.Open("mocks/video-6053793.json")
	} else if r.URL.Path == "/tags.json" {
		mock, err = os.Open("mocks/tags.json")
	} else if r.URL.Path == "/tags/86.json" {
		mock, err = os.Open("mocks/tag-86.json")
	} else {
		return nil, errors.New("unable to fetch URL, maybe a mock needs some setup")
	}

	rw := httptest.NewRecorder()

	io.Copy(rw, mock)

	resp := rw.Result()

	return resp, err
}

func TestQueryURLBuilding(t *testing.T) {
	c := NewClient("fake-token")

	tests := []struct {
		query    query
		expected string
	}{
		{
			c.Video(857),
			"http://api.video.globoi.com/videos/857.json?access_token=fake-token",
		},
		{
			c.Video(857).Fields("subscriber_only"),
			"http://api.video.globoi.com/videos/857.json?access_token=fake-token&only=subscriber_only",
		},
		{
			c.Videos(),
			"http://api.video.globoi.com/videos.json?access_token=fake-token",
		},
		{
			c.Videos().PerPage(15),
			"http://api.video.globoi.com/videos.json?access_token=fake-token&per_page=15",
		},
		{
			c.Videos().PerPage(15).Page(3),
			"http://api.video.globoi.com/videos.json?access_token=fake-token&page=3&per_page=15",
		},
		{
			c.Videos().Fields("subscriber_only", "extended_metadata").PerPage(15).Page(3),
			"http://api.video.globoi.com/videos.json?access_token=fake-token&only=subscriber_only%7Cextended_metadata&page=3&per_page=15",
		},
		{
			c.Videos().
				PerPage(5).
				AddTags("Flamengo"),
			"http://api.video.globoi.com/videos.json?access_token=fake-token&per_page=5&tags.all=Flamengo",
		},
		{
			c.Videos().
				PerPage(20).
				AddTags("Fluminense", "Vitória"),
			"http://api.video.globoi.com/videos.json?access_token=fake-token&per_page=20&tags.all=Fluminense%7CVit%C3%B3ria",
		},
		{
			c.Videos().
				PerPage(25).
				AddTags("futebol", "Tempo Real", "Flamengo", "Vasco"),
			"http://api.video.globoi.com/videos.json?access_token=fake-token&per_page=25&tags.all=futebol%7CTempo+Real%7CFlamengo%7CVasco",
		},
		{
			c.Videos().
				PerPage(5).
				PublishedSince(time.Date(2017, 3, 30, 0, 0, 0, 0, time.UTC)),
			"http://api.video.globoi.com/videos.json?access_token=fake-token&per_page=5&published_at.gte=2017-03-30T00%3A00%3A00",
		},
		{
			c.Videos().
				PerPage(5).
				PublishedUntil(time.Date(2017, 3, 30, 0, 0, 0, 0, time.UTC)),
			"http://api.video.globoi.com/videos.json?access_token=fake-token&per_page=5&published_at.lte=2017-03-30T00%3A00%3A00",
		},
		{
			c.Videos().
				PerPage(5).
				PublishedSince(time.Date(2017, 3, 25, 0, 0, 0, 0, time.UTC)).
				PublishedUntil(time.Date(2017, 3, 30, 0, 0, 0, 0, time.UTC)),
			"http://api.video.globoi.com/videos.json?access_token=fake-token&per_page=5&published_at.gte=2017-03-25T00%3A00%3A00&published_at.lte=2017-03-30T00%3A00%3A00",
		},
		{
			c.Tag(456),
			"http://api.video.globoi.com/tags/456.json?access_token=fake-token",
		},
		{
			c.Tags(),
			"http://api.video.globoi.com/tags.json?access_token=fake-token",
		},
		{
			c.Tags().Name("Futebol"),
			"http://api.video.globoi.com/tags.json?access_token=fake-token&name=Futebol",
		},
	}

	for _, test := range tests {
		u := c.buildURL(test.query.endpoint(), test.query.params())
		if test.expected != u.String() {
			t.Errorf("query URL mismatch: expected \"%s\" got \"%s\"", test.expected, u)
		}
	}
}

func TestVideoFetch(t *testing.T) {
	c := NewClient("fake-token", WithRoundTripper(MockRoundTripper{}))
	video, err := c.Video(5767587).Fetch()
	if err != nil {
		t.Fatal("unable to fetch video:", err)
	}

	id := 5767587
	if id != video.ID {
		t.Error("video ID mismatch: expected", id, "got", video.ID)
		t.Errorf("video ID mismatch: expected \"%d\" got \"%d\"", id, video.ID)
	}

	title := "Clássico entre Brasília e Flamengo, em Manaus, agita rodada do NBB"
	if title != video.Title {
		t.Errorf("video title mismatch: expected \"%s\" got \"%s\"", title, video.Title)
	}

	duration := 17067
	if duration != video.Duration {
		t.Errorf("video duration mismatch: expected \"%d\" got \"%d\"", duration, video.Duration)
	}

	tags := []string{"Flamengo", "Manaus", "Basquete"}
	if tags[0] != video.Tags[0] || tags[1] != video.Tags[1] || tags[2] != video.Tags[2] {
		t.Errorf("video tags mismatch: expected \"%+v\" got \"%+v\"", tags, video.Tags)
	}

	loc, _ := time.LoadLocation("America/Sao_Paulo")

	publishedAt := time.Date(2017, 3, 30, 13, 30, 52, 0, loc)
	if !publishedAt.Equal(video.PublishedAt) {
		t.Errorf("video published_at mismatch: expected \"%s\" got \"%s\"", publishedAt, video.PublishedAt)
	}

	updatedAt := time.Date(2017, 3, 31, 10, 10, 07, 0, loc)
	if !updatedAt.Equal(video.UpdatedAt) {
		t.Errorf("video updated_at mismatch: expected \"%s\" got \"%s\"", updatedAt, video.UpdatedAt)
	}

	if video.ExtendedMetadata != nil {
		t.Error("video should not have any extended metadata")
	}
}
func TestVideoWithExtendedMetadataFetch(t *testing.T) {
	c := NewClient("fake-token", WithRoundTripper(MockRoundTripper{}))
	video, err := c.Video(6053793).Fetch()
	if err != nil {
		t.Fatal("unable to fetch video:", err)
	}

	id := 6053793
	if id != video.ID {
		t.Error("video ID mismatch: expected", id, "got", video.ID)
		t.Errorf("video ID mismatch: expected \"%d\" got \"%d\"", id, video.ID)
	}

	title := "Violência tem Cor, Novos Expedientes, Todas as Vulvas"
	if title != video.Title {
		t.Errorf("video title mismatch: expected \"%s\" got \"%s\"", title, video.Title)
	}

	duration := 2859993
	if duration != video.Duration {
		t.Errorf("video duration mismatch: expected \"%d\" got \"%d\"", duration, video.Duration)
	}

	tags := []string{"Saia Justa", "GNT"}
	if tags[0] != video.Tags[0] || tags[1] != video.Tags[1] {
		t.Errorf("video tags mismatch: expected \"%+v\" got \"%+v\"", tags, video.Tags)
	}

	loc, _ := time.LoadLocation("America/Sao_Paulo")

	publishedAt := time.Date(2017, 8, 3, 17, 17, 41, 0, loc)
	if !publishedAt.Equal(video.PublishedAt) {
		t.Errorf("video published_at mismatch: expected \"%s\" got \"%s\"", publishedAt, video.PublishedAt)
	}

	updatedAt := time.Date(2017, 8, 3, 17, 17, 41, 0, loc)
	if !updatedAt.Equal(video.UpdatedAt) {
		t.Errorf("video updated_at mismatch: expected \"%s\" got \"%s\"", updatedAt, video.UpdatedAt)
	}

	if video.ExtendedMetadata == nil {
		t.Fatal("video should have extended metadata")
	}

	contentRating := "12"
	if contentRating != video.ExtendedMetadata.ContentRating {
		t.Errorf("video content rating mismatch: expected \"%s\" got \"%s\"",
			contentRating, video.ExtendedMetadata.ContentRating)
	}
}

func TestVideosFetch(t *testing.T) {
	c := NewClient("fake-token", WithRoundTripper(MockRoundTripper{}))
	videos, err := c.Videos().Fetch()
	if err != nil {
		t.Fatal("unable to fetch videos:", err)
	}

	expected := 5
	if vlen := len(videos); vlen != expected {
		t.Error("video count mismatch: expected", expected, "got", vlen)
	}
}

func TestTagFetch(t *testing.T) {
	c := NewClient("fake-token", WithRoundTripper(MockRoundTripper{}))
	tag, err := c.Tag(86).Fetch()
	if err != nil {
		t.Fatal("unable to fetch tag:", err)
	}

	id := 86
	if id != tag.ID {
		t.Errorf("tag ID mismatch: expected \"%d\" got \"%d\"", id, tag.ID)
	}

	name := "Futebol"
	if name != tag.Name {
		t.Errorf("tag name mismatch: expected \"%s\" got \"%s\"", name, tag.Name)
	}
}

func TestTagsFetch(t *testing.T) {
	c := NewClient("fake-token", WithRoundTripper(MockRoundTripper{}))
	tags, err := c.Tags().Fetch()
	if err != nil {
		t.Fatal("unable to fetch tags:", err)
	}

	expected := 2
	if tlen := len(tags); tlen != expected {
		t.Errorf("tag count mismatch: expected \"%d\" got \"%d\"", expected, tlen)
	}
}
