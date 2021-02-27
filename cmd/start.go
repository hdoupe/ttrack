package cmd

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hdoupe/ttrack/oauth"
	"github.com/hdoupe/ttrack/track"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start an entry",
	Long:  `Log the start time for an entry`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("either no arguments or one argument must be set")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating new time entry...")
		var startedAt time.Time
		var finishedAt time.Time
		var description string = ""

		if startedArg == "" {
			t := time.Now().UTC()
			startedAt = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
		} else {
			t, err := ParseTimeArg(startedArg)
			if err != nil {
				log.Fatal(err)
			}
			startedAt = t
		}
		if len(args) == 1 {
			description = args[0]
		}

		entry := track.Entry{
			StartedAt:   startedAt,
			FinishedAt:  finishedAt,
			Description: description,
			ClientID:    cfg.CurrentClient.ClientID,
			ProjectID:   cfg.CurrentClient.ProjectID,
		}

		client := oauth.Client{}
		tracker := GetTracker(client)
		entry = tracker.Start(entry)

		fmt.Println()
		fmt.Println(entry.String())
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
