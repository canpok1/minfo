package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/canpok1/minfo/internal"
	"github.com/canpok1/minfo/internal/misskey"
)

type Args struct {
	Origin string
	Limit  int
}

type Summary struct {
	Date          string
	NoteCount     int
	ReactionCount int
}

type Result struct {
	Summaries []Summary
	Latest    string
	Oldest    string
}

func main() {
	args, err := parseArgs(os.Args)
	if err != nil {
		panic(err)
	}

	if result, err := run(args.Origin, args.Limit); err != nil {
		panic(err)
	} else {
		if s, err := json.Marshal(*result); err != nil {
			panic(err)
		} else {
			fmt.Println(string(s))
		}
	}
}

func parseArgs(args []string) (*Args, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("not enough arguments")
	}

	origin := args[1]
	limit, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, err
	}
	return &Args{
		Origin: origin,
		Limit:  limit,
	}, nil
}

func run(origin string, limit int) (*Result, error) {
	client, err := misskey.NewClient(origin)
	if err != nil {
		return nil, fmt.Errorf("failed to create misskey client: %w", err)
	}

	notes, err := client.FetchLocalTimeline(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to FetchLocalTimeline: %w", err)
	}

	summaries := make(map[string]Summary)
	for _, note := range notes {
		reactionCount := 0
		for _, v := range note.Reactions {
			reactionCount = reactionCount + v
		}

		k := internal.FormatTime(internal.ToJST(note.CreatedAt), internal.YYYYMMDD)
		if summary, exists := summaries[k]; exists {
			summaries[k] = Summary{
				Date:          k,
				ReactionCount: summary.ReactionCount + reactionCount,
				NoteCount:     summary.NoteCount + 1,
			}
		} else {
			summaries[k] = Summary{
				Date:          k,
				ReactionCount: reactionCount,
				NoteCount:     1,
			}
		}
	}

	var keys []string
	for k := range summaries {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	result := Result{
		Summaries: nil,
		Latest:    internal.ToJST(notes[0].CreatedAt).Format(time.RFC3339),
		Oldest:    internal.ToJST(notes[len(notes)-1].CreatedAt).Format(time.RFC3339),
	}

	for _, k := range keys {
		result.Summaries = append(result.Summaries, summaries[k])
	}

	return &result, nil
}
