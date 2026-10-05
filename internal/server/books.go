package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"book-manager/internal/model"
)

// GET /books
func (s *Server) listBooks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := model.BookFilter{
		Author:          q.Get("author"),
		Genre:           q.Get("genre"),
		PublishedAfter:  q.Get("published_after"),
		PublishedBefore: q.Get("published_before"),
	}
	for _, d := range []string{filter.PublishedAfter, filter.PublishedBefore} {
		if d != "" && !validDate(d) {
			writeError(w, http.StatusBadRequest, "date filters must be YYYY-MM-DD")
			return
		}
	}

	books, err := s.store.ListBooks(r.Context(), filter)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, books)
}

// POST /books
func (s *Server) createBook(w http.ResponseWriter, r *http.Request) {
	b, ok := decodeBook(w, r)
	if !ok {
		return
	}
	created, err := s.store.CreateBook(r.Context(), b)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/books/"+strconv.FormatInt(created.ID, 10))
	writeJSON(w, http.StatusCreated, created)
}

// GET /books/{id}
func (s *Server) getBook(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	b, err := s.store.GetBook(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

// PUT /books/{id}
func (s *Server) updateBook(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	b, ok := decodeBook(w, r)
	if !ok {
		return
	}
	updated, err := s.store.UpdateBook(r.Context(), id, b)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// DELETE /books/{id}
func (s *Server) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := s.store.DeleteBook(r.Context(), id); err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func decodeBook(w http.ResponseWriter, r *http.Request) (model.Book, bool) {
	var b model.Book
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&b); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return model.Book{}, false
	}
	if b.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return model.Book{}, false
	}
	if b.Published != "" && !validDate(b.Published) {
		writeError(w, http.StatusBadRequest, "published must be YYYY-MM-DD")
		return model.Book{}, false
	}
	return b, true
}

func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be an integer")
		return 0, false
	}
	return id, true
}

func validDate(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}
