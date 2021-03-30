package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/hdoupe/ttrack/oauth"
	"github.com/spf13/cobra"
)

// var serviceName string

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to Freshbooks.",
	Long:  `Connect to a Freshbooks with Oauth.`,
	Run: func(cmd *cobra.Command, args []string) {
		oauthClient := oauth.Client{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
		}
		var credsError error
		creds, err := oauthClient.FromCache()
		if err != nil && !os.IsNotExist(err) {
			log.Fatal(err)
		} else if os.IsNotExist(err) {
			authURL := fmt.Sprintf("https://my.freshbooks.com/service/auth/oauth/authorize/?response_type=code&redirect_uri=%s&client_id=%s", oauth.RedirectURI, cfg.ClientID)
			fmt.Println("Go to link: ", authURL)

			fmt.Print("Enter authorization code: ")
			var code string
			fmt.Scanf("%s", &code)

			creds, credsError = oauthClient.Exchange(code)
			if credsError != nil {
				log.Fatal(credsError)
			}
			oauthClient.Cache(creds)
		}

		if oauthClient.IsExpired(creds) {
			fmt.Println("Credentials have expired. Refreshing credentials now.")
			creds, credsError = oauthClient.Refresh(creds)
			if credsError != nil {
				log.Fatal(credsError)
			}
			oauthClient.Cache(creds)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
