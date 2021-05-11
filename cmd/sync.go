package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/hdoupe/ttrack/oauth"
	"github.com/hdoupe/ttrack/track"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync time entries on FreshBooks with local time entries.",
	Long:  `Sync time entries on FreshBooks with local time entries.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Syncing time entries...")
		client := oauth.Client{}
		if !client.IsAuthenticated() {
			log.Fatal("Use 'ttrack connect' to log in to Freshbooks")
		}
		creds, err := client.FromCache()
		if err != nil {
			log.Fatal(err)
		}
		if client.IsExpired(creds) {
			fmt.Println("Refreshing expired credentials...")
			var refreshErr error
			creds, refreshErr = client.Refresh(creds)
			if refreshErr != nil {
				log.Fatal(refreshErr)
			}
			client.Cache(creds)
		}
		fbTracker := track.FreshBooks{
			LogLocation: logLocation,
			Credentials: creds,
		}
		fbTracker.SyncEntries()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
