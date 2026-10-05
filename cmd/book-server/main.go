package main

import (
	"flag"
	"log"
	"net/http"

	"book-manager/internal/server"
	"book-manager/internal/store"
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	dbPath := flag.String("db", "books.db", "path to the SQLite database file")
	flag.Parse()

	st, err := store.Open(*dbPath)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer st.Close()

	srv := server.New(st)
	log.Printf("book-server listening on %s (db: %s)", *addr, *dbPath)
	if err := http.ListenAndServe(*addr, srv); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
