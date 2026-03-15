package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "writebook",
	Short: "CLI for managing Writebook",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
}

func initConfig() {
	viper.SetConfigName(".writebook")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")

	viper.SetEnvPrefix("WRITEBOOK")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Error reading config: %s\n", err)
		}
	}
}

func getURL() string {
	url := viper.GetString("url")
	if url == "" {
		fmt.Fprintln(os.Stderr, "Error: WRITEBOOK_URL not set. Set it via environment variable or run 'writebook login'.")
		os.Exit(1)
	}
	return url
}

func getToken() string {
	token := viper.GetString("token")
	if token == "" {
		fmt.Fprintln(os.Stderr, "Error: WRITEBOOK_TOKEN not set. Run 'writebook login' first.")
		os.Exit(1)
	}
	return token
}
