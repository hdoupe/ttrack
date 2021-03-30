package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hdoupe/ttrack/oauth"
	"github.com/hdoupe/ttrack/track"
	"github.com/spf13/cobra"
)

var (
	lastArg  string
	sinceArg string
	untilArg string
	limitArg string
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "View time entry log.",
	Long:  `View time entries by querying time periods or description substrings.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			last  time.Duration
			since time.Time
			until time.Time
			limit int
			err   error
		)

		if lastArg != "" {
			last, err = time.ParseDuration(lastArg)
			if err != nil {
				log.Fatal(err)
			}
			since = time.Now().UTC().Add(-last)
		}

		if sinceArg != "" {
			since, err = ParseTimeArg(sinceArg)
			if err != nil {
				log.Fatal(err)
			}
		}

		if untilArg != "" {
			until, err = ParseTimeArg(untilArg)
			if err != nil {
				log.Fatal(err)
			}
		}

		if limitArg != "" {
			limit, err = strconv.Atoi(limitArg)
			if err != nil {
				log.Fatal(err)
			}
		}

		client := oauth.Client{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
		}
		tracker := GetTracker(client)

		entries := tracker.LoadEntries()

		params := track.FilterParameters{
			Since: since,
			Until: until,
			Limit: limit,
		}

		entries = track.FilterEntries(entries, params)

		if len(entries) == 0 {
			fmt.Println("No entries matched the query parameters.")
			return
		}

		var total time.Duration

		for _, entry := range entries {
			fmt.Println(entry.String())
			fmt.Println()
			d, _ := entry.GetDuration()
			total = total + d
		}

		fmt.Println("Total hours recorded: ", total.Round(time.Minute))

	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.
	logCmd.Flags().StringVarP(&lastArg, "last", "l", "", "Show entries over previous time period (eg. --last 1w).")
	logCmd.Flags().StringVarP(&limitArg, "limit", "n", "", "Show entries over previous time period (eg. --last 1w).")
	logCmd.Flags().StringVar(&sinceArg, "since", "", "Show entries starting from some date.")
	logCmd.Flags().StringVar(&untilArg, "until", "", "Show entries until some date.")
}
