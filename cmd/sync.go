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
		oauthClient := oauth.Client{}
		if !oauthClient.IsAuthenticated() {
			log.Fatal("Use 'ttrack connect' to log in to Freshbooks")
		}
		creds, err := oauthClient.FromCache()
		if err != nil {
			log.Fatal(err)
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
