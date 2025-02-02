package misskey

import (
	"fmt"
	"net/url"
	"time"

	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/services/notes/timeline"
)

type Note struct {
	ID        string
	CreatedAt time.Time
	Reactions ReactionCountMap
	Text      string
	User      User
}

type User struct {
	Username string
}

type ReactionCountMap map[string]uint64

type IgnoreUserNameSet map[string]struct{}

type Client struct {
	client *misskey.Client
}

func NewClient(origin string) (*Client, error) {
	o, err := url.Parse(origin)
	if err != nil {
		return nil, fmt.Errorf("failed to parse origin: %w", err)
	}
	client, err := misskey.NewClientWithOptions(
		misskey.WithBaseURL(o.Scheme, o.Host, ""),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return &Client{client: client}, nil
}

const localTimelineMaxLimit = 100

func (c *Client) FetchLocalTimeline(limit uint, ignoreSet IgnoreUserNameSet) ([]Note, error) {
	var filteredNotes []Note

	beforeId := ""
	for uint(len(filteredNotes)) < limit {
		requestLimit := min(limit-uint(len(filteredNotes)), localTimelineMaxLimit)

		notes, err := c.client.Notes().Timeline().Local(timeline.LocalRequest{Limit: requestLimit, UntilID: beforeId})
		if err != nil {
			return nil, fmt.Errorf("failed to fetchLocalTimeline(%d,%s): %w", limit, beforeId, err)
		}

		if len(notes) == 0 {
			break
		}

		for _, note := range notes {
			if _, exists := ignoreSet[note.User.Username]; !exists {
				filteredNote := Note{
					ID:        note.ID,
					CreatedAt: note.CreatedAt,
					Reactions: note.Reactions,
					Text:      note.Text,
					User:      User{Username: note.User.Username},
				}

				filteredNotes = append(filteredNotes, filteredNote)
			}
		}
		beforeId = notes[len(notes)-1].ID
	}

	return filteredNotes, nil
}
