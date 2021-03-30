package cmd

import (
	"fmt"
	"log"

	"github.com/hdoupe/ttrack/oauth"
	"github.com/hdoupe/ttrack/track"
)

// GetTracker returns the appropriate tracker by checking whether
// oauth credentials are present or not.
// (TODO: and other configuration.)
func GetTracker(client oauth.Client) track.Tracker {
	var tracker track.Tracker
	if client.IsAuthenticated() {
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
		}
		tracker = &track.FreshBooks{
			Credentials: creds,
			LogLocation: logLocation,
		}
	} else {
		tracker = &track.Local{LogLocation: logLocation}
	}

	return tracker
}
