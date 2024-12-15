package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/canpok1/minfo/internal/misskey"
)

type Args struct {
	origin string
	limit  int
}

type Summary struct {
	date          string
	reactionCount int
}

type Result struct {
	summaries []Summary
}

func main() {
	args, err := parseArgs(os.Args)
	if err != nil {
		panic(err)
	}

	if result, err := run(args.origin, args.limit); err != nil {
		panic(err)
	} else {
		for _, summary := range result.summaries {
			fmt.Printf("%s : %d\n", summary.date, summary.reactionCount)
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
		origin: origin,
		limit:  limit,
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

	sums := make(map[string]int)
	for _, note := range notes {
		count := 0
		for _, v := range note.Reactions {
			count = count + v
		}

		k := note.GetCreatedAtAsJST()
		if sum, exists := sums[k]; exists {
			sums[k] = sum + count
		} else {
			sums[k] = count
		}
	}

	var keys []string
	for k := range sums {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	result := Result{
		summaries: nil,
	}

	for _, k := range keys {
		result.summaries = append(result.summaries, Summary{date: k, reactionCount: sums[k]})
	}

	return &result, nil
}
