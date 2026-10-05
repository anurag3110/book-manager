package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"book-manager/internal/store"
)

type Server struct {
	store *store.Store
	mux   *http.ServeMux
}

func New(s *store.Store) *Server {
	srv := &Server{store: s, mux: http.NewServeMux()}

	srv.mux.HandleFunc("GET /books", srv.listBooks)
	srv.mux.HandleFunc("POST /books", srv.createBook)
	srv.mux.HandleFunc("GET /books/{id}", srv.getBook)
	srv.mux.HandleFunc("PUT /books/{id}", srv.updateBook)
	srv.mux.HandleFunc("DELETE /books/{id}", srv.deleteBook)

	srv.mux.HandleFunc("GET /collections", srv.listCollections)
	srv.mux.HandleFunc("POST /collections", srv.createCollection)
	srv.mux.HandleFunc("GET /collections/{name}", srv.getCollection)
	srv.mux.HandleFunc("PUT /collections/{name}", srv.updateCollection)
	srv.mux.HandleFunc("DELETE /collections/{name}", srv.deleteCollection)
	srv.mux.HandleFunc("PUT /collections/{name}/books/{id}", srv.addBookToCollection)
	srv.mux.HandleFunc("DELETE /collections/{name}/books/{id}", srv.removeBookFromCollection)

	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	case errors.Is(err, store.ErrConflict):
		writeError(w, http.StatusConflict, "name already in use")
	default:
		log.Printf("internal error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}
