Database schema

I'm using SQLite for this project but the SQLs r pretty much universal so it can work with any other relational DB.

Five tables: book, author, book_author join for the many-to-many relation between them, collection, and a collection_book join for the many-to-many relation between books and
collections.

```sql
CREATE TABLE book (
    id          INTEGER PRIMARY KEY,
    title       TEXT NOT NULL,
    published   TEXT,                    -- YYYY-MM-DD
    edition     TEXT NOT NULL DEFAULT '',
    genre       TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_book_genre     ON book (genre);
CREATE INDEX idx_book_published ON book (published);

CREATE TABLE author (
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE book_author (
    book_id   INTEGER NOT NULL REFERENCES book (id)   ON DELETE CASCADE,
    author_id INTEGER NOT NULL REFERENCES author (id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);

CREATE INDEX idx_book_author_author ON book_author (author_id);

CREATE TABLE collection (
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE collection_book (
    collection_id INTEGER NOT NULL REFERENCES collection (id) ON DELETE CASCADE,
    book_id       INTEGER NOT NULL REFERENCES book (id)       ON DELETE CASCADE,
    PRIMARY KEY (collection_id, book_id)
);

CREATE INDEX idx_collection_book_book ON collection_book (book_id);
```

Please note that:

SQLite has no date type so published is TEXT holding YYYY-MM-DD. That format
sorts chronologically as a plain string, so the date-range filter is an ordinary >=/<=. The server validates the format on input.

author name is unique. Authors are matched by name: creating a book with
"Chetan Bhagat" reuses the existing author or makes a new one.