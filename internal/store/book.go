package store

import (
	"context"
	"database/sql"
	"strings"

	"book-manager/internal/model"
)

func (s *Store) CreateBook(ctx context.Context, b model.Book) (model.Book, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Book{}, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`INSERT INTO book (title, published, edition, genre, description)
		 VALUES (?, ?, ?, ?, ?)`,
		b.Title, nullableDate(b.Published), b.Edition, b.Genre, b.Description)
	if err != nil {
		return model.Book{}, err
	}
	id, _ := res.LastInsertId()
	b.ID = id

	if err := linkAuthors(ctx, tx, id, b.Authors); err != nil {
		return model.Book{}, err
	}
	if err := tx.Commit(); err != nil {
		return model.Book{}, err
	}
	return b, nil
}

func (s *Store) GetBook(ctx context.Context, id int64) (model.Book, error) {
	var b model.Book
	var published sql.NullString
	err := s.db.QueryRowContext(ctx,
		`SELECT id, title, published, edition, genre, description FROM book WHERE id = ?`, id).
		Scan(&b.ID, &b.Title, &published, &b.Edition, &b.Genre, &b.Description)
	if err == sql.ErrNoRows {
		return model.Book{}, ErrNotFound
	}
	if err != nil {
		return model.Book{}, err
	}
	b.Published = published.String

	authors, err := s.authorsOf(ctx, b.ID)
	if err != nil {
		return model.Book{}, err
	}
	b.Authors = authors
	return b, nil
}

func (s *Store) ListBooks(ctx context.Context, f model.BookFilter) ([]model.Book, error) {
	var where []string
	var args []any

	if f.Author != "" {
		where = append(where, `id IN (
			SELECT ba.book_id FROM book_author ba
			JOIN author a ON a.id = ba.author_id
			WHERE a.name = ?)`)
		args = append(args, f.Author)
	}
	if f.Genre != "" {
		where = append(where, `genre = ?`)
		args = append(args, f.Genre)
	}
	if f.PublishedAfter != "" {
		where = append(where, `published IS NOT NULL AND published >= ?`)
		args = append(args, f.PublishedAfter)
	}
	if f.PublishedBefore != "" {
		where = append(where, `published IS NOT NULL AND published <= ?`)
		args = append(args, f.PublishedBefore)
	}

	query := `SELECT id, title, published, edition, genre, description FROM book`
	if len(where) > 0 {
		query += ` WHERE ` + strings.Join(where, ` AND `)
	}
	query += ` ORDER BY id DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []model.Book{}
	for rows.Next() {
		var b model.Book
		var published sql.NullString
		if err := rows.Scan(&b.ID, &b.Title, &published, &b.Edition, &b.Genre, &b.Description); err != nil {
			return nil, err
		}
		b.Published = published.String
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range books {
		authors, err := s.authorsOf(ctx, books[i].ID)
		if err != nil {
			return nil, err
		}
		books[i].Authors = authors
	}
	return books, nil
}

func (s *Store) UpdateBook(ctx context.Context, id int64, b model.Book) (model.Book, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Book{}, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`UPDATE book SET title = ?, published = ?, edition = ?, genre = ?, description = ?
		 WHERE id = ?`,
		b.Title, nullableDate(b.Published), b.Edition, b.Genre, b.Description, id)
	if err != nil {
		return model.Book{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.Book{}, ErrNotFound
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM book_author WHERE book_id = ?`, id); err != nil {
		return model.Book{}, err
	}
	if err := linkAuthors(ctx, tx, id, b.Authors); err != nil {
		return model.Book{}, err
	}
	if err := tx.Commit(); err != nil {
		return model.Book{}, err
	}

	b.ID = id
	return b, nil
}

func (s *Store) DeleteBook(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM book WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) authorsOf(ctx context.Context, bookID int64) ([]string, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT a.name FROM author a
		 JOIN book_author ba ON ba.author_id = a.id
		 WHERE ba.book_id = ?
		 ORDER BY a.id`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	names := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

func linkAuthors(ctx context.Context, tx *sql.Tx, bookID int64, names []string) error {
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		authorID, err := findOrCreateAuthor(ctx, tx, name)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO book_author (book_id, author_id) VALUES (?, ?)`,
			bookID, authorID); err != nil {
			return err
		}
	}
	return nil
}

func findOrCreateAuthor(ctx context.Context, tx *sql.Tx, name string) (int64, error) {
	if _, err := tx.ExecContext(ctx,
		`INSERT OR IGNORE INTO author (name) VALUES (?)`, name); err != nil {
		return 0, err
	}
	var id int64
	if err := tx.QueryRowContext(ctx,
		`SELECT id FROM author WHERE name = ?`, name).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func nullableDate(s string) any {
	if s == "" {
		return nil
	}
	return s
}
