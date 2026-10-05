package cli

import (
	"book-manager/internal/model"

	"github.com/spf13/cobra"
)

func collectionCommand(clientFn func() *Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "collection",
		Aliases: []string{"collections"},
		Short:   "Manage collections and their books",
	}
	cmd.AddCommand(
		collectionAddCommand(clientFn),
		collectionListCommand(clientFn),
		collectionShowCommand(clientFn),
		collectionEditCommand(clientFn),
		collectionDeleteCommand(clientFn),
		collectionAddBookCommand(clientFn),
		collectionRemoveBookCommand(clientFn),
	)
	return cmd
}

func collectionAddCommand(clientFn func() *Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add NAME",
		Short: "Create a collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			created, err := clientFn().CreateCollection(model.Collection{Name: args[0]})
			if err != nil {
				return err
			}
			cmd.Printf("Created collection %q\n", created.Name)
			return nil
		},
	}
	return cmd
}

func collectionListCommand(clientFn func() *Client) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List collections",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cols, err := clientFn().ListCollections()
			if err != nil {
				return err
			}
			printCollections(cmd.OutOrStdout(), cols)
			return nil
		},
	}
}

func collectionShowCommand(clientFn func() *Client) *cobra.Command {
	return &cobra.Command{
		Use:   "show NAME",
		Short: "Show a collection and the books in it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := clientFn().GetCollection(args[0])
			if err != nil {
				return err
			}
			printCollection(cmd.OutOrStdout(), c)
			return nil
		},
	}
}

func collectionEditCommand(clientFn func() *Client) *cobra.Command {
	var newName string
	cmd := &cobra.Command{
		Use:   "edit NAME --name NEW",
		Short: "Rename a collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			current := args[0]
			updated, err := clientFn().UpdateCollection(current, model.Collection{Name: newName})
			if err != nil {
				return err
			}
			cmd.Printf("Updated collection %q\n", updated.Name)
			return nil
		},
	}
	cmd.Flags().StringVar(&newName, "name", "", "new name for the collection")
	cmd.MarkFlagRequired("name")
	return cmd
}

func collectionDeleteCommand(clientFn func() *Client) *cobra.Command {
	return &cobra.Command{
		Use:   "delete NAME",
		Short: "Delete a collection (the books in it are kept)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := clientFn().DeleteCollection(args[0]); err != nil {
				return err
			}
			cmd.Printf("Deleted collection %q\n", args[0])
			return nil
		},
	}
}

func collectionAddBookCommand(clientFn func() *Client) *cobra.Command {
	return &cobra.Command{
		Use:   "add-book NAME BOOK_ID",
		Short: "Put a book in a collection (idempotent)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[1])
			if err != nil {
				return err
			}
			if err := clientFn().AddBookToCollection(args[0], id); err != nil {
				return err
			}
			cmd.Printf("Added book %d to %q\n", id, args[0])
			return nil
		},
	}
}

func collectionRemoveBookCommand(clientFn func() *Client) *cobra.Command {
	return &cobra.Command{
		Use:   "remove-book NAME BOOK_ID",
		Short: "Take a book out of a collection",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := parseID(args[1])
			if err != nil {
				return err
			}
			if err := clientFn().RemoveBookFromCollection(args[0], id); err != nil {
				return err
			}
			cmd.Printf("Removed book %d from %q\n", id, args[0])
			return nil
		},
	}
}
