package model

type Book struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Authors     []string `json:"authors"`
	Published   string   `json:"published"`
	Edition     string   `json:"edition"`
	Genre       string   `json:"genre"`
	Description string   `json:"description"`
}

type Collection struct {
	Name      string `json:"name"`
	BookCount int    `json:"book_count,omitempty"`
	Books     []Book `json:"books,omitempty"`
}

type BookFilter struct {
	Author          string
	Genre           string
	PublishedAfter  string
	PublishedBefore string
}
