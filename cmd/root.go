package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/canpok1/minfo/internal"
	"github.com/canpok1/minfo/internal/misskey"
	"github.com/spf13/cobra"
)

var (
	server          string
	limit           uint
	ignoreUserNames []string
)

type Summary struct {
	Date          string
	NoteCount     int
	ReactionCount uint64
}

type Result struct {
	Summaries       []Summary
	Latest          string
	Oldest          string
	IgnoreUserNames []string
}

var rootCmd = &cobra.Command{
	Use:   "minfo <server_url>",
	Short: "This is a tool to retrieve Misskey information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		server = args[0]

		client, err := misskey.NewClient(server)
		if err != nil {
			return fmt.Errorf("failed to create misskey client: %w", err)
		}

		ignoreSet := misskey.IgnoreUserNameSet{}
		for _, userName := range ignoreUserNames {
			ignoreSet[userName] = struct{}{}
		}

		notes, err := client.FetchLocalTimeline(limit, ignoreSet)
		if err != nil {
			return fmt.Errorf("failed to FetchLocalTimeline: %w", err)
		}

		summaries := make(map[string]Summary)
		for _, note := range notes {
			reactionCount := uint64(0)
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
			Summaries:       nil,
			Latest:          internal.ToJST(notes[0].CreatedAt).Format(time.RFC3339),
			Oldest:          internal.ToJST(notes[len(notes)-1].CreatedAt).Format(time.RFC3339),
			IgnoreUserNames: ignoreUserNames,
		}

		for _, k := range keys {
			result.Summaries = append(result.Summaries, summaries[k])
		}

		if s, err := json.MarshalIndent(result, "", "  "); err != nil {
			return err
		} else {
			fmt.Println(string(s))
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().UintVarP(&limit, "limit", "l", 50, "limit the number of notes")
	rootCmd.Flags().StringSliceVar(&ignoreUserNames, "ignore-usernames", []string{}, "list of usernames to ignore (comma separated). when the display name is xxxxx@yyyy, username is yyyy.")
}
