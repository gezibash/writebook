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
	rootCmd.AddCommand(picturesCmd)
	picturesCmd.AddCommand(picturesListCmd)
	picturesCmd.AddCommand(picturesCreateCmd)
	picturesCmd.AddCommand(picturesShowCmd)
	picturesCmd.AddCommand(picturesUpdateCmd)
	picturesCmd.AddCommand(picturesDeleteCmd)

	picturesCreateCmd.Flags().String("title", "", "Picture title (required)")
	picturesCreateCmd.Flags().String("caption", "", "Picture caption")
	picturesCreateCmd.Flags().String("image", "", "Path to image file (required)")
	picturesCreateCmd.Flags().Int("position", 0, "Position within the book")
	picturesCreateCmd.MarkFlagRequired("title")
	picturesCreateCmd.MarkFlagRequired("image")

	picturesUpdateCmd.Flags().String("title", "", "New title")
	picturesUpdateCmd.Flags().String("caption", "", "New caption")
	picturesUpdateCmd.Flags().String("image", "", "Path to new image file")
}

var picturesCmd = &cobra.Command{
	Use:   "pictures",
	Short: "Manage pictures within a book",
}

var picturesListCmd = &cobra.Command{
	Use:   "list <book_id>",
	Short: "List pictures in a book",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}

		c := client.NewClient(getURL(), getToken())
		pictures, err := c.ListPictures(bookID)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(pictures, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if len(pictures) == 0 {
			fmt.Println("No pictures found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tCAPTION\tHAS IMAGE\tPOSITION")
		for _, p := range pictures {
			hasImage := "no"
			if p.HasImage {
				hasImage = "yes"
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%.0f\n", p.ID, p.Title, p.Caption, hasImage, p.Position)
		}
		w.Flush()
		return nil
	},
}

var picturesCreateCmd = &cobra.Command{
	Use:   "create <book_id>",
	Short: "Create a new picture",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}

		title, _ := cmd.Flags().GetString("title")
		caption, _ := cmd.Flags().GetString("caption")
		image, _ := cmd.Flags().GetString("image")

		params := client.CreatePictureParams{
			Title:   title,
			Caption: caption,
			Image:   image,
		}
		if cmd.Flags().Changed("position") {
			v, _ := cmd.Flags().GetInt("position")
			params.Position = &v
		}

		c := client.NewClient(getURL(), getToken())
		leaf, err := c.CreatePicture(bookID, params)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(leaf, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Created picture #%d: %s\n", leaf.ID, leaf.Title)
		return nil
	},
}

var picturesShowCmd = &cobra.Command{
	Use:   "show <book_id> <picture_id>",
	Short: "Show a picture",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}
		pictureID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid picture_id: %s", args[1])
		}

		c := client.NewClient(getURL(), getToken())
		leaf, err := c.GetPicture(bookID, pictureID)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(leaf, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Picture #%d: %s\n", leaf.ID, leaf.Title)
		hasImage := "no"
		if leaf.HasImage {
			hasImage = "yes"
		}
		fmt.Printf("Has image: %s | Status: %s | Position: %.0f\n", hasImage, leaf.Status, leaf.Position)
		if leaf.Caption != "" {
			fmt.Printf("Caption: %s\n", leaf.Caption)
		}
		return nil
	},
}

var picturesUpdateCmd = &cobra.Command{
	Use:   "update <book_id> <picture_id>",
	Short: "Update a picture",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}
		pictureID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid picture_id: %s", args[1])
		}

		params := client.UpdatePictureParams{}
		if cmd.Flags().Changed("title") {
			v, _ := cmd.Flags().GetString("title")
			params.Title = &v
		}
		if cmd.Flags().Changed("caption") {
			v, _ := cmd.Flags().GetString("caption")
			params.Caption = &v
		}
		if cmd.Flags().Changed("image") {
			v, _ := cmd.Flags().GetString("image")
			params.Image = &v
		}

		c := client.NewClient(getURL(), getToken())
		leaf, err := c.UpdatePicture(bookID, pictureID, params)
		if err != nil {
			return err
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(leaf, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("Updated picture #%d: %s\n", leaf.ID, leaf.Title)
		return nil
	},
}

var picturesDeleteCmd = &cobra.Command{
	Use:   "delete <book_id> <picture_id>",
	Short: "Delete a picture",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bookID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid book_id: %s", args[0])
		}
		pictureID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid picture_id: %s", args[1])
		}

		c := client.NewClient(getURL(), getToken())
		if err := c.DeletePicture(bookID, pictureID); err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println("{}")
			return nil
		}

		fmt.Printf("Deleted picture #%d from book #%d\n", pictureID, bookID)
		return nil
	},
}
