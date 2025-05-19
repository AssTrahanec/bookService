package postres

import (
	"bookService/config"
	"bookService/internal/domain/models"
	"bookService/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strings"
)

type Storage struct {
	db *sqlx.DB
}

func New(cfg config.DBConfig) (*Storage, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) GetBook(ctx context.Context, id string) (*models.Book, error) {
	const op = "postgres.GetBook"
	const query = `
		SELECT 
			book_id as id, 
			title, 
			author, 
			publication_year as publicationyear, 
			genre 
		FROM books 
		WHERE book_id = $1
	`

	var book models.Book
	err := s.db.GetContext(ctx, &book, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrBookNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &book, nil
}
func (s *Storage) ListBooks(ctx context.Context, filter *models.BookFilter) ([]*models.Book, error) {
	const op = "postgres.ListBooks"
	baseQuery := `
		SELECT 
			book_id as id, 
			title, 
			author, 
			publication_year as publicationyear, 
			genre 
		FROM books 
		WHERE 1=1
	`
	var args []interface{}
	var conditions []string

	if filter != nil {
		if filter.Author != nil {
			args = append(args, *filter.Author)
			conditions = append(conditions, fmt.Sprintf("author = $%d", len(args)))
		}
		if filter.PublicationYear != nil {
			args = append(args, *filter.PublicationYear)
			conditions = append(conditions, fmt.Sprintf("publication_year = $%d", len(args)))
		}
		if filter.Genre != nil {
			args = append(args, *filter.Genre)
			conditions = append(conditions, fmt.Sprintf("genre = $%d", len(args)))
		}
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY title ASC"

	var books []*models.Book
	err := s.db.SelectContext(ctx, &books, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return books, nil
}
func (s *Storage) GetUserBooks(ctx context.Context, userID string, filter *models.BookFilter) ([]*models.Book, error) {
	const op = "postgres.GetUserBooks"
	baseQuery := `
		SELECT 
			b.book_id as id, 
			b.title, 
			b.author, 
			b.publication_year as publicationyear, 
			b.genre 
		FROM books b
		JOIN users_books ub ON b.book_id = ub.book_id
		WHERE ub.user_id = $1
	`
	args := []interface{}{userID}
	paramCounter := 2

	if filter != nil {
		if filter.Author != nil {
			args = append(args, *filter.Author)
			baseQuery += fmt.Sprintf(" AND b.author = $%d", paramCounter)
			paramCounter++
		}
		if filter.PublicationYear != nil {
			args = append(args, *filter.PublicationYear)
			baseQuery += fmt.Sprintf(" AND b.publication_year = $%d", paramCounter)
			paramCounter++
		}
		if filter.Genre != nil {
			args = append(args, *filter.Genre)
			baseQuery += fmt.Sprintf(" AND b.genre = $%d", paramCounter)
			paramCounter++
		}
	}

	var books []*models.Book
	err := s.db.SelectContext(ctx, &books, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return books, nil
}
func (s *Storage) AddBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	const op = "postgres.AddBook"
	const query = `
		INSERT INTO books (
			book_id, 
			title, 
			author, 
			publication_year, 
			genre
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING 
			book_id as id, 
			title, 
			author, 
			publication_year as publicationyear, 
			genre
	`

	if book.ID == "" {
		book.ID = uuid.New().String()
	}

	var result models.Book
	err := s.db.QueryRowxContext(ctx, query,
		book.ID,
		book.Title,
		book.Author,
		book.PublicationYear,
		book.Genre,
	).StructScan(&result)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &result, nil
}
func (s *Storage) UpdateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	const op = "postgres.UpdateBook"
	const query = `
		UPDATE books 
		SET 
			title = $1, 
			author = $2, 
			publication_year = $3, 
			genre = $4
		WHERE book_id = $5
		RETURNING 
			book_id as id, 
			title, 
			author, 
			publication_year as publicationyear, 
			genre
	`

	var result models.Book
	err := s.db.QueryRowxContext(ctx, query,
		book.Title,
		book.Author,
		book.PublicationYear,
		book.Genre,
		book.ID,
	).StructScan(&result)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrBookNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &result, nil
}
func (s *Storage) DeleteBook(ctx context.Context, id string) (string, error) {
	const op = "postgres.DeleteBook"
	const query = `
		DELETE FROM books 
		WHERE book_id = $1
		RETURNING book_id
	`

	var deletedID string
	err := s.db.QueryRowContext(ctx, query, id).Scan(&deletedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, storage.ErrBookNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return deletedID, nil
}

func (s *Storage) AddBookToUser(ctx context.Context, userID, bookID string) (string, error) {
	const op = "postgres.AddBookToUser"
	const query = `
        DELETE FROM users_books 
        WHERE user_id = $1 AND book_id = $2
        RETURNING book_id
    `

	var returnedBookID string
	err := s.db.QueryRowContext(ctx, query, userID, bookID).Scan(&returnedBookID)
	if err != nil {
		if err == sql.ErrNoRows {
			return bookID, nil
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return returnedBookID, nil
}
func (s *Storage) RemoveBookFromUser(ctx context.Context, userID, bookID string) (string, error) {
	const op = "postgres.RemoveBookFromUser"
	const query = `
        DELETE FROM users_books 
        WHERE user_id = $1 AND book_id = $2
        RETURNING book_id
    `

	var returnedBookID string
	err := s.db.QueryRowContext(ctx, query, userID, bookID).Scan(&returnedBookID)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, storage.ErrBookNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return returnedBookID, nil
}
