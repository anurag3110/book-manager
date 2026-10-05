package server

import (
	"encoding/json"
	"net/http"
	"net/url"

	"book-manager/internal/model"
)

// GET /collections
func (s *Server) listCollections(w http.ResponseWriter, r *http.Request) {
	cols, err := s.store.ListCollections(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cols)
}

// POST /collections
func (s *Server) createCollection(w http.ResponseWriter, r *http.Request) {
	c, ok := decodeCollection(w, r)
	if !ok {
		return
	}
	created, err := s.store.CreateCollection(r.Context(), c)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/collections/"+url.PathEscape(created.Name))
	writeJSON(w, http.StatusCreated, created)
}

// GET /collections/{name}
func (s *Server) getCollection(w http.ResponseWriter, r *http.Request) {
	c, err := s.store.GetCollection(r.Context(), r.PathValue("name"))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// PUT /collections/{name}
func (s *Server) updateCollection(w http.ResponseWriter, r *http.Request) {
	c, ok := decodeCollection(w, r)
	if !ok {
		return
	}
	updated, err := s.store.UpdateCollection(r.Context(), r.PathValue("name"), c)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// DELETE /collections/{name}
func (s *Server) deleteCollection(w http.ResponseWriter, r *http.Request) {
	if err := s.store.DeleteCollection(r.Context(), r.PathValue("name")); err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

// PUT /collections/{name}/books/{id}
func (s *Server) addBookToCollection(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := s.store.AddBook(r.Context(), r.PathValue("name"), id); err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, nil)
}

// DELETE /collections/{name}/books/{id}
func (s *Server) removeBookFromCollection(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := s.store.RemoveBook(r.Context(), r.PathValue("name"), id); err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func decodeCollection(w http.ResponseWriter, r *http.Request) (model.Collection, bool) {
	var c model.Collection
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return model.Collection{}, false
	}
	if c.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return model.Collection{}, false
	}
	return c, true
}
