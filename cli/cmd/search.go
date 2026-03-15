package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"writebook/client"
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().Int("book", 0, "Limit search to a specific book ID")
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search across book content",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		bookID, _ := cmd.Flags().GetInt("book")

		c := client.NewClient(getURL(), getToken())

		if bookID > 0 {
			return searchBook(c, bookID, query)
		}
		return searchGlobal(c, query)
	},
}

func searchGlobal(c *client.Client, query string) error {
	groups, err := c.SearchGlobal(query)
	if err != nil {
		return err
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(groups, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	if len(groups) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	for i, group := range groups {
		if i > 0 {
			fmt.Println()
		}
		count := len(group.Results)
		noun := "result"
		if count != 1 {
			noun = "results"
		}
		fmt.Printf("%s (%d %s)\n", group.BookTitle, count, noun)
		for _, r := range group.Results {
			fmt.Printf("  %-4d %-8s %s\n", r.ID, r.Type, r.Title)
			if r.ContentSnippet != "" {
				fmt.Printf("  %s%s\n", strings.Repeat(" ", 13), r.ContentSnippet)
			}
		}
	}

	return nil
}

func searchBook(c *client.Client, bookID int, query string) error {
	results, err := c.SearchBook(bookID, query)
	if err != nil {
		return err
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	if len(results) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTYPE\tTITLE\tSNIPPET")
	for _, r := range results {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", r.ID, r.Type, r.Title, r.ContentSnippet)
	}
	w.Flush()
	return nil
}
