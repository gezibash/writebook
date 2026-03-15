package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"writebook/client"
)

func init() {
	rootCmd.AddCommand(booksCmd)
	booksCmd.AddCommand(booksListCmd)
	booksCmd.AddCommand(booksCreateCmd)
	booksCmd.AddCommand(booksShowCmd)
	booksCmd.AddCommand(booksUpdateCmd)
	booksCmd.AddCommand(booksDeleteCmd)

	booksCreateCmd.Flags().String("title", "", "Book title (required)")
	booksCreateCmd.Flags().String("subtitle", "", "Book subtitle")
	booksCreateCmd.Flags().String("author", "", "Book author")
	booksCreateCmd.Flags().String("theme", "", "Book theme (black, blue, green, magenta, orange, violet, white)")
	booksCreateCmd.MarkFlagRequired("title")

	booksUpdateCmd.Flags().String("title", "", "New title")
	booksUpdateCmd.Flags().String("subtitle", "", "New subtitle")
	booksUpdateCmd.Flags().String("author", "", "New author")
	booksUpdateCmd.Flags().String("theme", "", "New theme")

	booksDeleteCmd.Flags().Bool("force", false, "Skip confirmation prompt")
}

var booksCmd = &cobra.Command{
	Use:   "books",
	Short: "Manage books",
}

var booksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accessible books",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient(getURL(), getToken())
		books, err := c.ListBooks()
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(books, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(books) == 0 {
			fmt.Println("No books found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tAUTHOR\tPUBLISHED\tTHEME")
		for _, b := range books {
			published := "no"
			if b.Published {
				published = "yes"
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", b.ID, b.Title, b.Author, published, b.Theme)
		}
		w.Flush()
		return nil
	},
}

var booksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new book",
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		subtitle, _ := cmd.Flags().GetString("subtitle")
		author, _ := cmd.Flags().GetString("author")
		theme, _ := cmd.Flags().GetString("theme")

		params := client.CreateBookParams{
			Title:    title,
			Subtitle: subtitle,
			Author:   author,
			Theme:    theme,
		}

		c := client.NewClient(getURL(), getToken())
		book, err := c.CreateBook(params)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(book, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Created book #%d: %s\n", book.ID, book.Title)
		return nil
	},
}

var booksShowCmd = &cobra.Command{
	Use:   "show <book_id>",
	Short: "Show a book with its leaves",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}

		c := client.NewClient(getURL(), getToken())
		book, err := c.GetBook(bookID)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(book, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Book #%d: %s\n", book.ID, book.Title)
		if book.Subtitle != "" {
			fmt.Printf("Subtitle: %s\n", book.Subtitle)
		}
		if book.Author != "" {
			fmt.Printf("Author: %s\n", book.Author)
		}
		fmt.Printf("Theme: %s | Published: %v\n", book.Theme, book.Published)
		fmt.Println()

		if len(book.Leaves) == 0 {
			fmt.Println("No leaves.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTYPE\tTITLE")
		for _, l := range book.Leaves {
			fmt.Fprintf(w, "%d\t%s\t%s\n", l.ID, l.Type, l.Title)
		}
		w.Flush()
		return nil
	},
}

var booksUpdateCmd = &cobra.Command{
	Use:   "update <book_id>",
	Short: "Update a book",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}

		params := client.UpdateBookParams{}
		if cmd.Flags().Changed("title") {
			v, _ := cmd.Flags().GetString("title")
			params.Title = &v
		}
		if cmd.Flags().Changed("subtitle") {
			v, _ := cmd.Flags().GetString("subtitle")
			params.Subtitle = &v
		}
		if cmd.Flags().Changed("author") {
			v, _ := cmd.Flags().GetString("author")
			params.Author = &v
		}
		if cmd.Flags().Changed("theme") {
			v, _ := cmd.Flags().GetString("theme")
			params.Theme = &v
		}

		c := client.NewClient(getURL(), getToken())
		book, err := c.UpdateBook(bookID, params)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(book, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Updated book #%d: %s\n", book.ID, book.Title)
		return nil
	},
}

var booksDeleteCmd = &cobra.Command{
	Use:   "delete <book_id>",
	Short: "Delete a book",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			fmt.Printf("Are you sure you want to delete book #%d? [y/N] ", bookID)
			var answer string
			fmt.Scanln(&answer)
			if !strings.EqualFold(answer, "y") && !strings.EqualFold(answer, "yes") {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		c := client.NewClient(getURL(), getToken())
		if err := c.DeleteBook(bookID); err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println("{}")
			return nil
		}

		fmt.Printf("Deleted book #%d\n", bookID)
		return nil
	},
}
