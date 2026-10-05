refer the /docs folder for documents of steps 1, 2, 3.
- cli.md : Step 1: Describe the expected user experience
- rest-api.md: Step 2: Describe the expected REST API
- db-schema.md: Step 3: Describe the database structure

build and deploy the REST APIs:
```
go build ./...
go run ./cmd/book-server --db books.db --addr :8080
```

example cli commands: (first "book" is the cli handler and the other is module name)
```
book book add --title "The Three Mistakes of My Life" \
    --author "Chetan Bhagat" --author "Krish Malhotra" \
    --published 2008-01-01 --genre fiction

book book list --genre fiction

book collection add "To read"
book collection add-book "To read" 1
book collection show "To read"
```