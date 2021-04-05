package cmd

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hdoupe/ttrack/track"
	"github.com/spf13/cobra"
)

var (
	clientNickname string
	clientIDArg    string
	projectIDArg   string
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "clients",
	Short: "Manage clients",
}

var addClientCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new client",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			clientID  int
			projectID int
			err       error
		)

		if clientID, err = strconv.Atoi(clientIDArg); err != nil {
			log.Fatal("Client ID must be an integer.")
		}
		if projectID, err = strconv.Atoi(projectIDArg); err != nil {
			log.Fatal("Project ID must be an integer.")
		}

		newClient := track.Client{
			Nickname:  clientNickname,
			ClientID:  clientID,
			ProjectID: projectID,
		}
		clients, newClientErr := track.AddClient(cfg.Clients, newClient)
		if newClientErr != nil {
			log.Fatal(newClientErr)
		}
		cfg.Clients = clients
		if configErr := WriteConfig(cfg); configErr != nil {
			log.Fatal(configErr)
		}

		fmt.Println("Added new client:")
		fmt.Println(newClient.String())
	},
}

var setCurrentClientCmd = &cobra.Command{
	Use:   "set-current",
	Short: "Set current client",
	Long:  `Set the client to be used when logging time entries.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("either no arguments or one argument must be set")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			clientNickname = args[0]
		} else if len(args) != 0 {
			log.Fatal("Only one client may be used at a time.")
		}
		clients := track.FilterClients(cfg.Clients, track.Client{Nickname: clientNickname})
		if len(clients) == 0 {
			log.Fatal("No clients found with nickname: ", clientNickname)
		}
		if len(clients) > 1 {
			log.Fatal("More than one client found with nickname: ", clientNickname)
		}
		cfg.CurrentClient = clients[0]
		if configErr := WriteConfig(cfg); configErr != nil {
			log.Fatal(configErr)
		}
		fmt.Println("Current client set to:", cfg.CurrentClient.Nickname)
	},
}

var listClientsCommand = &cobra.Command{
	Use:   "list",
	Short: "List clients",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			clientID  int
			projectID int
			err       error
		)

		if clientIDArg != "" {
			if clientID, err = strconv.Atoi(clientIDArg); err != nil {
				log.Fatal("Client ID must be an integer.")
			}
		}
		if projectIDArg != "" {
			if projectID, err = strconv.Atoi(projectIDArg); err != nil {
				log.Fatal("Project ID must be an integer.")
			}
		}

		clients := track.FilterClients(cfg.Clients, track.Client{Nickname: clientNickname, ClientID: clientID, ProjectID: projectID})

		if len(clients) == 0 {
			fmt.Println("No clients found.")
		} else {
			for _, client := range clients {
				fmt.Println(client.String())
			}
		}
	},
}

var getCurrentClientCmd = &cobra.Command{
	Use:   "get-current",
	Short: "Get current client",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if (cfg.CurrentClient == track.Client{}) {
			fmt.Println("Current client has not been configured yet.")
		} else {
			fmt.Println(cfg.CurrentClient.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.AddCommand(addClientCmd)
	clientCmd.AddCommand(setCurrentClientCmd)
	clientCmd.AddCommand(listClientsCommand)
	clientCmd.AddCommand(getCurrentClientCmd)

	// Here you will define your flags and configuration settings.

	clientCmd.PersistentFlags().StringVar(&clientNickname, "nickname", "", "Nickname for client")
	clientCmd.PersistentFlags().StringVar(&clientIDArg, "client-id", "", "ID for client")
	clientCmd.PersistentFlags().StringVar(&projectIDArg, "project-id", "", "ID for project")
}
