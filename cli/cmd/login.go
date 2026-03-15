package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"

	"writebook/client"
)

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().String("url", "", "Writebook server URL")
	loginCmd.Flags().String("email", "", "Login email")
	loginCmd.Flags().String("password", "", "Login password")
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and save API token",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			url = viper.GetString("url")
		}
		if url == "" {
			fmt.Print("Writebook URL: ")
			input, _ := reader.ReadString('\n')
			url = strings.TrimSpace(input)
		}

		email, _ := cmd.Flags().GetString("email")
		if email == "" {
			fmt.Print("Email: ")
			input, _ := reader.ReadString('\n')
			email = strings.TrimSpace(input)
		}

		password, _ := cmd.Flags().GetString("password")
		if password == "" {
			fmt.Print("Password: ")
			passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				return fmt.Errorf("reading password: %w", err)
			}
			password = string(passwordBytes)
		}

		c := client.NewClient(url, "")
		resp, err := c.Login(email, password)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		// Save to config file
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("finding home directory: %w", err)
		}

		viper.Set("url", url)
		viper.Set("token", resp.Token)
		configPath := home + "/.writebook.yaml"
		if err := viper.WriteConfigAs(configPath); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("Logged in as %s (%s)\n", resp.User.Name, resp.User.Role)
		fmt.Printf("Config saved to %s\n", configPath)
		return nil
	},
}
