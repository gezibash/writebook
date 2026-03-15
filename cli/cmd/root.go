package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "writebook",
	Short: "CLI for managing Writebook",
}

func SetVersion(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
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

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding home directory: %s\n", err)
		os.Exit(1)
	}
	return filepath.Join(home, ".config", "writebook")
}

func configPath() string {
	return filepath.Join(configDir(), "config.toml")
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(configDir())

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
