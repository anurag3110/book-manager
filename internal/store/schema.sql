CREATE TABLE IF NOT EXISTS book (
    id          INTEGER PRIMARY KEY,
    title       TEXT NOT NULL,
    published   TEXT,                    -- YYYY-MM-DD
    edition     TEXT NOT NULL DEFAULT '',
    genre       TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_book_genre     ON book (genre);
CREATE INDEX IF NOT EXISTS idx_book_published ON book (published);

CREATE TABLE IF NOT EXISTS author (
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS book_author (
    book_id   INTEGER NOT NULL REFERENCES book (id)   ON DELETE CASCADE,
    author_id INTEGER NOT NULL REFERENCES author (id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);

CREATE INDEX IF NOT EXISTS idx_book_author_author ON book_author (author_id);

CREATE TABLE IF NOT EXISTS collection (
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS collection_book (
    collection_id INTEGER NOT NULL REFERENCES collection (id) ON DELETE CASCADE,
    book_id       INTEGER NOT NULL REFERENCES book (id)       ON DELETE CASCADE,
    PRIMARY KEY (collection_id, book_id)
);

CREATE INDEX IF NOT EXISTS idx_collection_book_book ON collection_book (book_id);
