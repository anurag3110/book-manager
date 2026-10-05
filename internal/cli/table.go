package cli

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"book-manager/internal/model"
)

func printBooks(w io.Writer, books []model.Book) {
	tw := newTab(w)
	fmt.Fprintln(tw, "ID\tTITLE\tAUTHORS\tPUBLISHED\tGENRE")
	for _, b := range books {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\n",
			b.ID, b.Title, strings.Join(b.Authors, ", "), dash(b.Published), dash(b.Genre))
	}
	tw.Flush()
}

func printBook(w io.Writer, b model.Book) {
	tw := newTab(w)
	fmt.Fprintf(tw, "ID\t%d\n", b.ID)
	fmt.Fprintf(tw, "Title\t%s\n", b.Title)
	fmt.Fprintf(tw, "Authors\t%s\n", dash(strings.Join(b.Authors, ", ")))
	fmt.Fprintf(tw, "Published\t%s\n", dash(b.Published))
	fmt.Fprintf(tw, "Edition\t%s\n", dash(b.Edition))
	fmt.Fprintf(tw, "Genre\t%s\n", dash(b.Genre))
	fmt.Fprintf(tw, "Description\t%s\n", dash(b.Description))
	tw.Flush()
}

func printCollections(w io.Writer, cols []model.Collection) {
	tw := newTab(w)
	fmt.Fprintln(tw, "NAME\tBOOKS")
	for _, c := range cols {
		fmt.Fprintf(tw, "%s\t%d\n", c.Name, c.BookCount)
	}
	tw.Flush()
}

func printCollection(w io.Writer, c model.Collection) {
	fmt.Fprintf(w, "%s", c.Name)
	fmt.Fprintf(w, "  (%s)\n\n", plural(len(c.Books), "book"))
	if len(c.Books) == 0 {
		fmt.Fprintln(w, "(empty)")
		return
	}
	printBooks(w, c.Books)
}

func newTab(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
}

func dash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return strconv.Itoa(n) + " " + noun + "s"
}
