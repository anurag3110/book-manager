package cli

import (
	"os"

	"github.com/spf13/cobra"
)

const defaultServer = "http://localhost:8080"

func Execute() error {
	var serverURL string

	root := &cobra.Command{
		Use:           "book",
		Short:         "Manage a personal library of books and collections",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().StringVar(&serverURL, "server", envOr("BOOK_SERVER", defaultServer),
		"book-server base URL (or set BOOK_SERVER)")

	client := func() *Client { return NewClient(serverURL) }

	root.AddCommand(bookCommand(client))
	root.AddCommand(collectionCommand(client))

	return root.Execute()
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
