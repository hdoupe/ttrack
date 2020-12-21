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

// finishCmd represents the finish command
var finishCmd = &cobra.Command{
	Use:   "finish",
	Short: "Finish an entry",
	Long:  `Finish a previously started entry.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("either no arguments or one argument must be set")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Finishing last time entry...")
		var startedAt time.Time
		var finishedAt time.Time
		var description string = ""

		if len(args) == 1 {
			description = args[0]
		}

		var duration time.Duration
		if durationArg != "" {
			var err error
			duration, err = time.ParseDuration(durationArg)
			if err != nil {
				log.Fatal(err)
			}
		} else if finishedArg == "" {
			t := time.Now().UTC()
			finishedAt = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
		} else {
			t, err := ParseTimeArg(finishedArg)
			if err != nil {
				log.Fatal(err)
			}
			finishedAt = t
		}

		entry := track.Entry{
			StartedAt:   startedAt,
			FinishedAt:  finishedAt,
			Description: description,
			Duration:    int(duration.Seconds()),
		}

		client := oauth.Client{}
		tracker := GetTracker(client)
		entry = tracker.Finish(entry)

		fmt.Println()
		fmt.Println(entry.String())
	},
}

func init() {
	rootCmd.AddCommand(finishCmd)
}
