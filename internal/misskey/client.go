package misskey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Message string      `json:"message"`
	Code    string      `json:"code"`
	ID      string      `json:"id"`
	Info    interface{} `json:"info"`
}

type Note struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Reactions Reaction  `json:"reactions"`
	Text      string    `json:"text"`
}

type Reaction map[string]int

type Client struct {
	origin *url.URL
}

func NewClient(origin string) (*Client, error) {
	o, err := url.Parse(origin)
	if err != nil {
		return nil, fmt.Errorf("failed to parse origin: %w", err)
	}
	return &Client{origin: o}, nil
}

const localTimelineMaxLimit = 100

func (c *Client) FetchLocalTimeline(limit int) ([]Note, error) {
	var notes []Note

	beforeId := ""
	requestedCount := 0
	for requestedCount < limit {
		count := limit - requestedCount
		if count > localTimelineMaxLimit {
			count = localTimelineMaxLimit
		}

		if n, err := c.fetchLocalTimeline(count, beforeId); err != nil {
			return nil, fmt.Errorf("failed to fetchLocalTimeline(%d,%s): %w", limit, beforeId, err)
		} else {
			notes = append(notes, n...)
			beforeId = n[len(n)-1].ID
			requestedCount = requestedCount + count
		}
	}

	return notes, nil
}

func (c *Client) fetchLocalTimeline(limit int, untilID string) ([]Note, error) {
	type requestBody struct {
		Limit   int    `json:"limit"`
		UntilID string `json:"untilId,omitempty"`
	}

	data := requestBody{Limit: limit, UntilID: untilID}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshal: %w", err)
	}

	url := *c.origin
	url.Path = path.Join(url.Path, "api/notes/local-timeline")

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to NewRequest: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errorResponse ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			fmt.Printf("response body: %s\n", body)
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		return nil, fmt.Errorf("%v", errorResponse.Error)
	}

	var notes []Note
	err = json.Unmarshal(body, &notes)
	if err != nil {
		fmt.Printf("response: %s\n", body)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return notes, nil
}
