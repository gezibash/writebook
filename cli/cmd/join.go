package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"

	"writebook/client"
)

func init() {
	rootCmd.AddCommand(joinCmd)

	joinCmd.Flags().String("name", "", "Your display name")
	joinCmd.Flags().String("email", "", "Your email address")
	joinCmd.Flags().String("password", "", "Your password")
}

var joinCmd = &cobra.Command{
	Use:   "join <join_url>",
	Short: "Join a Writebook instance via invite link",
	Long: `Join a Writebook instance using an invite URL like:
  writebook join http://localhost:3007/join/gFoO-0Lkb-UcFa

Creates your account and saves the API token to ~/.writebook.yaml.
Pass --name, --email, --password flags for non-interactive usage.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		joinURL := args[0]

		parsed, err := url.Parse(joinURL)
		if err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		}

		// Extract join code from the path (last segment of /join/<code>)
		joinCode := path.Base(parsed.Path)
		if joinCode == "" || joinCode == "." || joinCode == "join" {
			return fmt.Errorf("could not extract join code from URL: %s", joinURL)
		}

		// Derive the base URL (scheme + host)
		baseURL := fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)

		reader := bufio.NewReader(os.Stdin)

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fmt.Print("Name: ")
			input, _ := reader.ReadString('\n')
			name = strings.TrimSpace(input)
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

		c := client.NewClient(baseURL, "")
		resp, err := c.Join(joinCode, name, email, password)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		// Save to config
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("finding home directory: %w", err)
		}

		viper.Set("url", baseURL)
		viper.Set("token", resp.Token)
		configPath := home + "/.writebook.yaml"
		if err := viper.WriteConfigAs(configPath); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("Joined as %s (%s)\n", resp.User.Name, resp.User.Role)
		fmt.Printf("Config saved to %s\n", configPath)
		return nil
	},
}
