package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"github.com/hdoupe/ttrack/track"
)

// Config describes the structure of the ttrack configuration.
type Config struct {
	ClientID      string         `mapstructure:"clientID"`
	ClientSecret  string         `mapstructure:"clientSecret"`
	LogLocation   string         `mapstructure:"logLocation"`
	CurrentClient track.Client   `mapstructure:"currentClient"`
	Clients       []track.Client `mapstructure:"clients"`
}

var (
	cfgFile     string
	cfg         Config
	startedArg  string
	finishedArg string
	logLocation string
	durationArg string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ttrack",
	Short: "Time tracking CLI application",
	Long:  `A tool for tracking time.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.ttrack.yaml)")

	rootCmd.PersistentFlags().StringVarP(&startedArg, "started-at", "s", "", "start time for entry")
	rootCmd.PersistentFlags().StringVarP(&finishedArg, "finished-at", "f", "", "finish time for entry")
	rootCmd.PersistentFlags().StringVarP(&durationArg, "duration", "d", "", "entry duration e.g. 30m (can be used instead of finished-at)")
	rootCmd.PersistentFlags().StringVar(&logLocation, "log-path", "~/.ttrack.log.json", "path to time entry log")
}

func loadConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ttrack" (without extension).
		viper.SetConfigName(".ttrack")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	loadConfig()
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Fatal("unable to decode into struct", err)
		}
		if len(cfg.Clients) == 0 {
			cfg.CurrentClient = track.Client{Nickname: "default", ProjectID: 0, ClientID: 0}
			cfg.Clients = []track.Client{cfg.CurrentClient}

			if err := WriteConfig(cfg); err != nil {
				log.Fatal("unable to decode into struct", err)
			}
		}
		fmt.Printf("Using client: %s\n\n", cfg.CurrentClient.Nickname)
	}
}

// WriteConfig writes the current config to the config file.
func WriteConfig(newConfig Config) error {
	viper.Set("logLocation", newConfig.LogLocation)
	viper.Set("clients", newConfig.Clients)
	viper.Set("currentClient", newConfig.CurrentClient)
	return viper.WriteConfig()
}
