package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"book-manager/internal/model"
)

type Client struct {
	base string
	http *http.Client
}

func NewClient(base string) *Client {
	return &Client{base: strings.TrimRight(base, "/"), http: http.DefaultClient}
}

func (c *Client) CreateBook(b model.Book) (model.Book, error) {
	var out model.Book
	err := c.do("POST", "/books", b, &out)
	return out, err
}

func (c *Client) ListBooks(f model.BookFilter) ([]model.Book, error) {
	q := url.Values{}
	if f.Author != "" {
		q.Set("author", f.Author)
	}
	if f.Genre != "" {
		q.Set("genre", f.Genre)
	}
	if f.PublishedAfter != "" {
		q.Set("published_after", f.PublishedAfter)
	}
	if f.PublishedBefore != "" {
		q.Set("published_before", f.PublishedBefore)
	}
	path := "/books"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	var out []model.Book
	err := c.do("GET", path, nil, &out)
	return out, err
}

func (c *Client) GetBook(id int64) (model.Book, error) {
	var out model.Book
	err := c.do("GET", "/books/"+strconv.FormatInt(id, 10), nil, &out)
	return out, err
}

func (c *Client) UpdateBook(id int64, b model.Book) (model.Book, error) {
	var out model.Book
	err := c.do("PUT", "/books/"+strconv.FormatInt(id, 10), b, &out)
	return out, err
}

func (c *Client) DeleteBook(id int64) error {
	return c.do("DELETE", "/books/"+strconv.FormatInt(id, 10), nil, nil)
}

func (c *Client) CreateCollection(col model.Collection) (model.Collection, error) {
	var out model.Collection
	err := c.do("POST", "/collections", col, &out)
	return out, err
}

func (c *Client) ListCollections() ([]model.Collection, error) {
	var out []model.Collection
	err := c.do("GET", "/collections", nil, &out)
	return out, err
}

func (c *Client) GetCollection(name string) (model.Collection, error) {
	var out model.Collection
	err := c.do("GET", "/collections/"+url.PathEscape(name), nil, &out)
	return out, err
}

func (c *Client) UpdateCollection(name string, col model.Collection) (model.Collection, error) {
	var out model.Collection
	err := c.do("PUT", "/collections/"+url.PathEscape(name), col, &out)
	return out, err
}

func (c *Client) DeleteCollection(name string) error {
	return c.do("DELETE", "/collections/"+url.PathEscape(name), nil, nil)
}

func (c *Client) AddBookToCollection(name string, bookID int64) error {
	return c.do("PUT", "/collections/"+url.PathEscape(name)+"/books/"+strconv.FormatInt(bookID, 10), nil, nil)
}

func (c *Client) RemoveBookFromCollection(name string, bookID int64) error {
	return c.do("DELETE", "/collections/"+url.PathEscape(name)+"/books/"+strconv.FormatInt(bookID, 10), nil, nil)
}

func (c *Client) do(method, path string, body, out any) error {
	var reqBody io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(buf)
	}

	req, err := http.NewRequest(method, c.base+path, reqBody)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("cannot reach server at %s: %w", c.base, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s", serverError(resp))
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

func serverError(resp *http.Response) string {
	var payload struct {
		Error string `json:"error"`
	}
	data, _ := io.ReadAll(resp.Body)
	if json.Unmarshal(data, &payload) == nil && payload.Error != "" {
		return payload.Error
	}
	return resp.Status
}
