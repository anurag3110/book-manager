package cli

import (
	"fmt"
	"strconv"

	"book-manager/internal/model"

	"github.com/spf13/cobra"
)

func bookCommand(clientFn func() *Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "book",
		Short: "Add, list, show, edit and delete books",
	}
	cmd.AddCommand(
		bookAddCommand(clientFn),
		bookListCommand(clientFn),
		bookShowCommand(clientFn),
		bookEditCommand(clientFn),
		bookDeleteCommand(clientFn),
	)
	return cmd
}

func bookAddCommand(clientFn func() *Client) *cobra.Command {
	var b model.Book
	cmd := &cobra.Command{
		Use:   "add --title TITLE [flags]",
		Short: "Add a book",
		RunE: func(cmd *cobra.Command, args []string) error {
			created, err := clientFn().CreateBook(b)
			if err != nil {
				return err
			}
			printBook(cmd.OutOrStdout(), created)
			return nil
		},
	}
	bookFlags(cmd, &b)
	cmd.MarkFlagRequired("title")
	return cmd
}

func bookListCommand(clientFn func() *Client) *cobra.Command {
	var f model.BookFilter
	cmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "List books, optionally filtered",
		RunE: func(cmd *cobra.Command, args []string) error {
			books, err := clientFn().ListBooks(f)
			if err != nil {
				return err
			}
			printBooks(cmd.OutOrStdout(), books)
			return nil
		},
	}
	cmd.Flags().StringVar(&f.Author, "author", "", "only books by this author")
	cmd.Flags().StringVar(&f.Genre, "genre", "", "only books in this genre")
	cmd.Flags().StringVar(&f.PublishedAfter, "published-after", "", "only books published on or after this date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&f.PublishedBefore, "published-before", "", "only books published on or before this date (YYYY-MM-DD)")
	return cmd
}

func bookShowCommand(clientFn func() *Client) *cobra.Command {
	return &cobra.Command{
		Use:   "show ID",
		Short: "Show one book in full",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			b, err := clientFn().GetBook(id)
			if err != nil {
				return err
			}
			printBook(cmd.OutOrStdout(), b)
			return nil
		},
	}
}

func bookEditCommand(clientFn func() *Client) *cobra.Command {
	var b model.Book
	cmd := &cobra.Command{
		Use:   "edit ID [flags]",
		Short: "Replace a book's details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			updated, err := clientFn().UpdateBook(id, b)
			if err != nil {
				return err
			}
			printBook(cmd.OutOrStdout(), updated)
			return nil
		},
	}
	bookFlags(cmd, &b)
	cmd.MarkFlagRequired("title")
	return cmd
}

func bookDeleteCommand(clientFn func() *Client) *cobra.Command {
	return &cobra.Command{
		Use:   "delete ID",
		Short: "Delete a book",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[0])
			if err != nil {
				return err
			}
			if err := clientFn().DeleteBook(id); err != nil {
				return err
			}
			cmd.Printf("Deleted book %d\n", id)
			return nil
		},
	}
}

func bookFlags(cmd *cobra.Command, b *model.Book) {
	cmd.Flags().StringVar(&b.Title, "title", "", "book title (required)")
	cmd.Flags().StringArrayVar(&b.Authors, "author", nil, "author name; repeat for multiple authors")
	cmd.Flags().StringVar(&b.Published, "published", "", "publication date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&b.Edition, "edition", "", "edition, e.g. \"2nd\"")
	cmd.Flags().StringVar(&b.Genre, "genre", "", "genre")
	cmd.Flags().StringVar(&b.Description, "description", "", "short description")
}

func parseID(s string) (int64, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("book id must be a number, got %q", s)
	}
	return id, nil
}
