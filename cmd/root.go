package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var startedArg string
var finishedArg string
var logLocation string
var durationArg string

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

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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
		viper.AddConfigPath(home)
		viper.SetConfigName(".ttrack")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		if logLocation, err = homedir.Expand(viper.Get("log").(string)); err != nil {
			panic(err)
		}
	}
}
