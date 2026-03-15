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
	rootCmd.AddCommand(sectionsCmd)
	sectionsCmd.AddCommand(sectionsListCmd)
	sectionsCmd.AddCommand(sectionsCreateCmd)
	sectionsCmd.AddCommand(sectionsShowCmd)
	sectionsCmd.AddCommand(sectionsUpdateCmd)
	sectionsCmd.AddCommand(sectionsDeleteCmd)

	sectionsCreateCmd.Flags().String("title", "", "Section title (required)")
	sectionsCreateCmd.Flags().String("body", "", "Section body")
	sectionsCreateCmd.Flags().String("theme", "", "Section theme")
	sectionsCreateCmd.Flags().Int("position", 0, "Position within the book")
	sectionsCreateCmd.MarkFlagRequired("title")

	sectionsUpdateCmd.Flags().String("title", "", "New title")
	sectionsUpdateCmd.Flags().String("body", "", "New body")
	sectionsUpdateCmd.Flags().String("theme", "", "New theme")
}

var sectionsCmd = &cobra.Command{
	Use:   "sections",
	Short: "Manage sections within a book",
}

var sectionsListCmd = &cobra.Command{
	Use:   "list <book_id>",
	Short: "List sections in a book",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}

		c := client.NewClient(getURL(), getToken())
		sections, err := c.ListSections(bookID)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(sections, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(sections) == 0 {
			fmt.Println("No sections found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tTHEME\tSTATUS\tPOSITION")
		for _, s := range sections {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%.0f\n", s.ID, s.Title, s.Theme, s.Status, s.Position)
		}
		w.Flush()
		return nil
	},
}

var sectionsCreateCmd = &cobra.Command{
	Use:   "create <book_id>",
	Short: "Create a new section",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}

		title, _ := cmd.Flags().GetString("title")
		body, _ := cmd.Flags().GetString("body")
		theme, _ := cmd.Flags().GetString("theme")

		params := client.CreateSectionParams{Title: title, Body: body, Theme: theme}
		if cmd.Flags().Changed("position") {
			v, _ := cmd.Flags().GetInt("position")
			params.Position = &v
		}
		c := client.NewClient(getURL(), getToken())
		leaf, err := c.CreateSection(bookID, params)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(leaf, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Created section #%d: %s\n", leaf.ID, leaf.Title)
		return nil
	},
}

var sectionsShowCmd = &cobra.Command{
	Use:   "show <book_id> <section_id>",
	Short: "Show a section",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}
		sectionID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid section_id: %s", args[1])
		}

		c := client.NewClient(getURL(), getToken())
		leaf, err := c.GetSection(bookID, sectionID)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(leaf, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Section #%d: %s\n", leaf.ID, leaf.Title)
		fmt.Printf("Theme: %s | Status: %s | Position: %.0f\n", leaf.Theme, leaf.Status, leaf.Position)
		if leaf.Body != "" {
			fmt.Println()
			fmt.Println(leaf.Body)
		}
		return nil
	},
}

var sectionsUpdateCmd = &cobra.Command{
	Use:   "update <book_id> <section_id>",
	Short: "Update a section",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}
		sectionID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid section_id: %s", args[1])
		}

		params := client.UpdateSectionParams{}
		if cmd.Flags().Changed("title") {
			v, _ := cmd.Flags().GetString("title")
			params.Title = &v
		}
		if cmd.Flags().Changed("body") {
			v, _ := cmd.Flags().GetString("body")
			params.Body = &v
		}
		if cmd.Flags().Changed("theme") {
			v, _ := cmd.Flags().GetString("theme")
			params.Theme = &v
		}

		c := client.NewClient(getURL(), getToken())
		leaf, err := c.UpdateSection(bookID, sectionID, params)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(leaf, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Updated section #%d: %s\n", leaf.ID, leaf.Title)
		return nil
	},
}

var sectionsDeleteCmd = &cobra.Command{
	Use:   "delete <book_id> <section_id>",
	Short: "Delete a section",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}
		sectionID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid section_id: %s", args[1])
		}

		c := client.NewClient(getURL(), getToken())
		if err := c.DeleteSection(bookID, sectionID); err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println("{}")
			return nil
		}

		fmt.Printf("Deleted section #%d from book #%d\n", sectionID, bookID)
		return nil
	},
}
