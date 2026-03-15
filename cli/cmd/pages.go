package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"writebook/client"
)

func init() {
	rootCmd.AddCommand(pagesCmd)
	pagesCmd.AddCommand(pagesListCmd)
	pagesCmd.AddCommand(pagesCreateCmd)
	pagesCmd.AddCommand(pagesShowCmd)
	pagesCmd.AddCommand(pagesUpdateCmd)
	pagesCmd.AddCommand(pagesDeleteCmd)

	pagesCreateCmd.Flags().String("title", "", "Page title (required)")
	pagesCreateCmd.Flags().String("body", "", "Page body (markdown)")
	pagesCreateCmd.Flags().String("body-file", "", "Read body from file")
	pagesCreateCmd.Flags().Int("position", 0, "Position within the book")
	pagesCreateCmd.MarkFlagRequired("title")

	pagesUpdateCmd.Flags().String("title", "", "New title")
	pagesUpdateCmd.Flags().String("body", "", "New body (markdown)")
	pagesUpdateCmd.Flags().String("body-file", "", "Read body from file")
}

var pagesCmd = &cobra.Command{
	Use:   "pages",
	Short: "Manage pages within a book",
}

var pagesListCmd = &cobra.Command{
	Use:   "list <book_id>",
	Short: "List pages in a book",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}

		c := client.NewClient(getURL(), getToken())
		pages, err := c.ListPages(bookID)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(pages, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(pages) == 0 {
			fmt.Println("No pages found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tSTATUS\tPOSITION")
		for _, p := range pages {
			fmt.Fprintf(w, "%d\t%s\t%s\t%.0f\n", p.ID, p.Title, p.Status, p.Position)
		}
		w.Flush()
		return nil
	},
}

var pagesCreateCmd = &cobra.Command{
	Use:   "create <book_id>",
	Short: "Create a new page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}

		title, _ := cmd.Flags().GetString("title")
		body, err := resolveBody(cmd)
		if err != nil {
			return err
		}

		params := client.CreatePageParams{Title: title, Body: body}
		if cmd.Flags().Changed("position") {
			v, _ := cmd.Flags().GetInt("position")
			params.Position = &v
		}
		c := client.NewClient(getURL(), getToken())
		leaf, err := c.CreatePage(bookID, params)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(leaf, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Created page #%d: %s\n", leaf.ID, leaf.Title)
		return nil
	},
}

var pagesShowCmd = &cobra.Command{
	Use:   "show <book_id> <page_id>",
	Short: "Show a page",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}
		pageID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid page_id: %s", args[1])
		}

		c := client.NewClient(getURL(), getToken())
		leaf, err := c.GetPage(bookID, pageID)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(leaf, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Page #%d: %s\n", leaf.ID, leaf.Title)
		fmt.Printf("Status: %s | Position: %.0f\n", leaf.Status, leaf.Position)
		if leaf.Body != "" {
			fmt.Println()
			fmt.Println(leaf.Body)
		}
		return nil
	},
}

var pagesUpdateCmd = &cobra.Command{
	Use:   "update <book_id> <page_id>",
	Short: "Update a page",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}
		pageID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid page_id: %s", args[1])
		}

		params := client.UpdatePageParams{}
		if cmd.Flags().Changed("title") {
			v, _ := cmd.Flags().GetString("title")
			params.Title = &v
		}
		body, err := resolveBody(cmd)
		if err != nil {
			return err
		}
		if cmd.Flags().Changed("body") || cmd.Flags().Changed("body-file") {
			params.Body = &body
		}

		c := client.NewClient(getURL(), getToken())
		leaf, err := c.UpdatePage(bookID, pageID, params)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(leaf, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Updated page #%d: %s\n", leaf.ID, leaf.Title)
		return nil
	},
}

var pagesDeleteCmd = &cobra.Command{
	Use:   "delete <book_id> <page_id>",
	Short: "Delete a page",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}
		pageID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid page_id: %s", args[1])
		}

		c := client.NewClient(getURL(), getToken())
		if err := c.DeletePage(bookID, pageID); err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println("{}")
			return nil
		}

		fmt.Printf("Deleted page #%d from book #%d\n", pageID, bookID)
		return nil
	},
}

func resolveBody(cmd *cobra.Command) (string, error) {
	bodyFile, _ := cmd.Flags().GetString("body-file")
	if bodyFile != "" {
		data, err := os.ReadFile(bodyFile)
		if err != nil {
			return "", fmt.Errorf("reading body file: %w", err)
		}
		return string(data), nil
	}
	body, _ := cmd.Flags().GetString("body")
	return body, nil
}
