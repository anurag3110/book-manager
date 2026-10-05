REST APIs:


status codes:
- 200: request succeeded
- 201: resource created
- 204: request succeeded (delete)
- 400: malformed JSON or a missing field
- 404: not found (book or collection)

Books

1. GET /books

Lists books

```
$ curl http://localhost:8080/books
```
```json
[
  {
    "id": 4,
    "title": "The Three Mistakes of My Life",
    "authors": ["Chetan Bhagat", "Krish Malhotra"],
    "published": "2008-01-01",
    "edition": "1st",
    "genre": "fiction",
    "description": "Three friends in Ahmedabad chase cricket, business and love."
  },
  {
    "id": 3,
    "title": "What Young India Wants",
    "authors": ["Chetan Bhagat"],
    "published": "2012-08-06",
    "edition": "",
    "genre": "non fiction",
    "description": ""
  },
  {
    "id": 2,
    "title": "A Game of Thrones",
    "authors": ["George R. R. Martin"],
    "published": "1996-08-01",
    "edition": "",
    "genre": "fantasy",
    "description": ""
  },
  {
    "id": 1,
    "title": "The Adventures of Sherlock Holmes",
    "authors": ["Arthur Conan Doyle"],
    "published": "1890-01-01",
    "edition": "",
    "genre": "fiction",
    "description": ""
  }
]
```

The query params are all optional and combine with AND:
- author: books by this exact author name
- genre: books with this exact genre
- published_after: published on or after this date 
- published_before: published on or before this date 

```
$ curl "http://localhost:8080/books?genre=fiction&published_after=2000-01-01&published_before=2010-12-31"
```
```json
[
  {
    "id": 4,
    "title": "The Three Mistakes of My Life",
    "authors": ["Chetan Bhagat", "Krish Malhotra"],
    "published": "2008-01-01",
    "edition": "1st",
    "genre": "fiction",
    "description": "Three friends in Ahmedabad chase cricket, business and love."
  }
]
```


2. POST /books

Creates a book. title is required; everything else is optional. authors is an array:

```
$ curl -X POST http://localhost:8080/books \
    -H 'Content-Type: application/json' \
    -d '{
      "title": "A Game of Thrones",
      "authors": ["George R. R. Martin"],
      "published": "1996-08-01",
      "genre": "fantasy"
    }'
```
```
HTTP/1.1 201 Created
Location: /books/5
```
```json
{
  "id": 5,
  "title": "A Game of Thrones",
  "authors": ["George R. R. Martin"],
  "published": "1996-08-01",
  "edition": "",
  "genre": "fantasy",
  "description": ""
}
```

format for published must be a YYYY-MM-DD date; any other format is invalid for now.

3. GET /books/{id}

return a book with given id

```
$ curl http://localhost:8080/books/2
```
```json
{
  "id": 2,
  "title": "A Game of Thrones",
  "authors": ["George R. R. Martin"],
  "published": "1996-08-01",
  "edition": "",
  "genre": "fantasy",
  "description": ""
}
```

4. PUT /books/{id}

completely replaces a full book object: fields left out are cleared, so send the whole object.

```
$ curl -X PUT http://localhost:8080/books/2 \
    -H 'Content-Type: application/json' \
    -d '{ "title": "A Game of Thrones", "authors": ["George R. R. Martin"], "published": "1996-08-01", "genre": "fantasy" }'
```

Returns 200 and the updated book. 404 if the id doesn't exist.

5. DELETE /books/{id}

Deletes the book and its entry in any collection. Returns 204.

```
$ curl -X DELETE http://localhost:8080/books/4
```


--------------------------------------------------------------



Collections:

Collections are addressed by name in the path.

1. GET /collections

Lists collections. 

```
$ curl http://localhost:8080/collections
```
```json
[
  { "name": "Owned" },
  { "name": "To read", "book_count": 2 }
]
```

2. POST /collections

Creates a collection. name is required and must be unique; 409 if it's taken.

```
$ curl -X POST http://localhost:8080/collections \
    -H 'Content-Type: application/json' \
    -d '{ "name": "To read" }'
```
```
HTTP/1.1 201 Created
Location: /collections/To%20read
```

3. GET /collections/{name}

returns the collection and its book.

```
$ curl http://localhost:8080/collections/To%20read
```
```json
{
  "name": "To read",
  "books": [
    { "id": 2, "title": "A Game of Thrones", "authors": ["George R. R. Martin"], "published": "1996-08-01", "edition": "", "genre": "fantasy", "description": "" },
    { "id": 1, "title": "The Adventures of Sherlock Holmes", "authors": ["Arthur Conan Doyle"], "published": "1890-01-01", "edition": "", "genre": "fiction", "description": "" }
  ]
}
```


4. DELETE /collections/{name}

Deletes the collection. Returns 204.

5. PUT /collections/{name}/books/{id}

Puts a book in the collection. Adding a book again to the same collection will be ignored so there are no duplicate books in a collection.

```
$ curl -X PUT http://localhost:8080/collections/To%20read/books/2
```

6. DELETE /collections/{name}/books/{id}

remove the book from this collection and returns 204. If the book isn't in the collection,
or the collection or book doesn't exist, returns 404.

```
$ curl -X DELETE http://localhost:8080/collections/To%20read/books/2
```
