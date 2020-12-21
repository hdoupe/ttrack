package cmd

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hdoupe/ttrack/oauth"
	"github.com/hdoupe/ttrack/track"
	"github.com/spf13/cobra"
)

var agoArg string

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the values of the most recent entry.",
	Long:  `Edit start, finish, or description values.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("either no arguments or one argument must be set")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var ago int

		if agoArg != "" {
			i, err := strconv.Atoi(agoArg)
			if err != nil {
				log.Fatal(err)
			}
			ago = i
		} else {
			ago = 1
		}

		client := oauth.Client{}
		tracker := GetTracker(client)

		entries := tracker.LoadEntries()
		if len(entries) < ago {
			log.Fatal("There are only ", len(entries), "which is less than ago: ", ago, ".")
		}
		entry := entries[len(entries)-ago]

		if startedArg != "" {
			t, err := ParseTimeArg(startedArg)
			if err != nil {
				log.Fatal(err)
			}
			entry.StartedAt = t
		}
		if finishedArg != "" && durationArg != "" {
			log.Fatal("Only one of finished-at and duration can be specified.")
		}
		if finishedArg != "" {
			t, err := ParseTimeArg(finishedArg)
			if err != nil {
				log.Fatal(err)
			}
			entry.FinishedAt = t
			entry.Duration = int(t.Sub(entry.StartedAt).Seconds())
		}
		if durationArg != "" {
			d, err := time.ParseDuration(durationArg)
			if err != nil {
				log.Fatal(err)
			}
			entry.Duration = int(d.Seconds())
			entry.FinishedAt = entry.StartedAt.Add(d).UTC()
		}
		if len(args) == 1 {
			entry.Description = args[0]
		}

		fmt.Println()
		fmt.Println(entry.String())
		tracker.SaveEntries([]track.Entry{entry})

	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringVarP(&agoArg, "ago", "a", "", "Edit ago-th most recent entry.")
}
