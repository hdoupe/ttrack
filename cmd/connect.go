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
		oauthClient := oauth.Client{}
		creds, err := oauthClient.FromCache()
		if err != nil && !os.IsNotExist(err) {
			log.Fatal(err)
		} else if os.IsNotExist(err) {
			fmt.Println("Go to link: ", oauth.AuthURL)

			fmt.Print("Enter authorization code:")
			var code string
			fmt.Scanf("%s", &code)

			creds = oauthClient.Exchange(code)
			oauthClient.Cache(creds)
		}

		if oauthClient.IsExpired(creds) {
			fmt.Println("Credentials have expired. Refreshing credentials now.")
			creds = oauthClient.Refresh(creds)
			oauthClient.Cache(creds)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
