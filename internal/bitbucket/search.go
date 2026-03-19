package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type SearchResult struct {
	ContentMatchCount int `json:"content_match_count"`
	ContentMatches    []struct {
		Lines []struct {
			Line     int `json:"line"`
			Segments []struct {
				Text  string `json:"text"`
				Match bool   `json:"match"`
			} `json:"segments"`
		} `json:"lines"`
	} `json:"content_matches"`
	PathMatches []struct {
		Text  string `json:"text"`
		Match bool   `json:"match"`
	} `json:"path_matches"`
	File struct {
		Path string `json:"path"`
	} `json:"file"`
}

type SearchResponse struct {
	Size    int            `json:"size"`
	Page    int            `json:"page"`
	PageLen int            `json:"pagelen"`
	Next    string         `json:"next"`
	Values  []SearchResult `json:"values"`
}

func (c *Client) SearchCode(workspace, query string, opts *PaginationOptions) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("search_query", query)
	if opts != nil {
		opts.applyParams(params)
	}
	ensurePageLen(params)

	path := fmt.Sprintf("/workspaces/%s/search/code?%s",
		url.PathEscape(workspace), params.Encode())

	if opts != nil && opts.All {
		var allResults []SearchResult
		currentPath := path
		totalSize := 0
		for currentPath != "" {
			data, err := c.get(currentPath)
			if err != nil {
				if len(allResults) == 0 {
					return nil, err
				}
				break
			}
			var result SearchResponse
			if err := json.Unmarshal(data, &result); err != nil {
				return nil, fmt.Errorf("parsing search results: %w", err)
			}
			allResults = append(allResults, result.Values...)
			totalSize = result.Size
			currentPath = result.Next
		}
		return &SearchResponse{
			Size:   totalSize,
			Values: allResults,
		}, nil
	}

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var result SearchResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parsing search results: %w", err)
	}
	return &result, nil
}
