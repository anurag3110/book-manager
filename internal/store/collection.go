package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"book-manager/internal/model"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func (s *Store) CreateCollection(ctx context.Context, c model.Collection) (model.Collection, error) {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO collection (name) VALUES (?)`, c.Name)
	if err != nil {
		if isUniqueViolation(err) {
			return model.Collection{}, ErrConflict
		}
		return model.Collection{}, err
	}
	return c, nil
}

func (s *Store) ListCollections(ctx context.Context) ([]model.Collection, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT c.name, COUNT(cb.book_id)
		 FROM collection c
		 LEFT JOIN collection_book cb ON cb.collection_id = c.id
		 GROUP BY c.id
		 ORDER BY c.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.Collection{}
	for rows.Next() {
		var c model.Collection
		if err := rows.Scan(&c.Name, &c.BookCount); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Store) GetCollection(ctx context.Context, name string) (model.Collection, error) {
	var c model.Collection
	var id int64
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name FROM collection WHERE name = ?`, name).
		Scan(&id, &c.Name)
	if err == sql.ErrNoRows {
		return model.Collection{}, ErrNotFound
	}
	if err != nil {
		return model.Collection{}, err
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT b.id, b.title, b.published, b.edition, b.genre, b.description
		 FROM book b
		 JOIN collection_book cb ON cb.book_id = b.id
		 WHERE cb.collection_id = ?
		 ORDER BY b.id DESC`, id)
	if err != nil {
		return model.Collection{}, err
	}
	defer rows.Close()

	c.Books = []model.Book{}
	for rows.Next() {
		var b model.Book
		var published sql.NullString
		if err := rows.Scan(&b.ID, &b.Title, &published, &b.Edition, &b.Genre, &b.Description); err != nil {
			return model.Collection{}, err
		}
		b.Published = published.String
		c.Books = append(c.Books, b)
	}
	if err := rows.Err(); err != nil {
		return model.Collection{}, err
	}

	for i := range c.Books {
		authors, err := s.authorsOf(ctx, c.Books[i].ID)
		if err != nil {
			return model.Collection{}, err
		}
		c.Books[i].Authors = authors
	}
	return c, nil
}

func (s *Store) UpdateCollection(ctx context.Context, name string, c model.Collection) (model.Collection, error) {
	res, err := s.db.ExecContext(ctx,
		`UPDATE collection SET name = ? WHERE name = ?`,
		c.Name, name)
	if err != nil {
		if isUniqueViolation(err) {
			return model.Collection{}, ErrConflict
		}
		return model.Collection{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.Collection{}, ErrNotFound
	}
	return c, nil
}

func (s *Store) DeleteCollection(ctx context.Context, name string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM collection WHERE name = ?`, name)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) AddBook(ctx context.Context, collectionName string, bookID int64) error {
	collectionID, err := s.collectionID(ctx, collectionName)
	if err != nil {
		return err
	}
	if err := s.requireBook(ctx, bookID); err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO collection_book (collection_id, book_id) VALUES (?, ?)`,
		collectionID, bookID)
	return err
}

func (s *Store) RemoveBook(ctx context.Context, collectionName string, bookID int64) error {
	collectionID, err := s.collectionID(ctx, collectionName)
	if err != nil {
		return err
	}
	if err := s.requireBook(ctx, bookID); err != nil {
		return err
	}
	res, err := s.db.ExecContext(ctx,
		`DELETE FROM collection_book WHERE collection_id = ? AND book_id = ?`,
		collectionID, bookID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) collectionID(ctx context.Context, name string) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx, `SELECT id FROM collection WHERE name = ?`, name).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, ErrNotFound
	}
	return id, err
}

func (s *Store) requireBook(ctx context.Context, id int64) error {
	var one int
	err := s.db.QueryRowContext(ctx, `SELECT 1 FROM book WHERE id = ?`, id).Scan(&one)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}

func isUniqueViolation(err error) bool {
	var se *sqlite.Error
	if errors.As(err, &se) {
		code := se.Code()
		return code == sqlite3.SQLITE_CONSTRAINT_UNIQUE || code == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
