package webmedia

import (
	"encoding/json"
)

type Pager struct {
	TotalEntries int  `json:"total_entries"`
	TotalPages   int  `json:"total_pages"`
	PerPage      int  `json:"per_page"`
	Offset       int  `json:"offset"`
	PreviousPage *int `json:"previous_page"`
	CurrentPage  int  `json:"current_page"`
	NextPage     *int `json:"next_page"`
}

type VideosPagerQuery struct {
	VideosQuery
}

func (v VideosPagerQuery) endpoint() string {
	return "videos/pagination.json"
}

func (v VideosPagerQuery) Fetch() (Pager, error) {
	var pager Pager

	resp, err := v.client.fetch(v)
	if err != nil {
		return pager, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&pager)
	return pager, err
}
