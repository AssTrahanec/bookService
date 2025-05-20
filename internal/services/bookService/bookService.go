package bookService

import (
	"bookService/internal/domain/models"
	"context"
	"fmt"
	"log/slog"
)

type BookService struct {
	log          *slog.Logger
	bookSaver    BookSaver
	bookProvider BookProvider
	bookCache    BookCache
}

type BookSaver interface {
	AddBook(ctx context.Context, book *models.Book) (*models.Book, error)
	UpdateBook(ctx context.Context, book *models.Book) (*models.Book, error)
	DeleteBook(ctx context.Context, id string) (string, error)
	RemoveBookFromUser(ctx context.Context, userID, bookID string) (string, error)
	AddBookToUser(ctx context.Context, userID, bookID string) (string, error)
}
type BookProvider interface {
	GetBook(ctx context.Context, id string) (*models.Book, error)
	ListBooks(ctx context.Context, filter *models.BookFilter) ([]*models.Book, error)
	GetUserBooks(ctx context.Context, userID string, filter *models.BookFilter) ([]*models.Book, error)
}
type BookCache interface {
	GetBook(ctx context.Context, id string) (*models.Book, error)
	SetBook(ctx context.Context, key string, book *models.Book) error
	InvalidateBook(ctx context.Context, key string) error
}

func New(
	bookSaver BookSaver,
	bookProvider BookProvider,
	bookCache BookCache,
	log *slog.Logger,
) *BookService {
	return &BookService{
		bookSaver:    bookSaver,
		bookProvider: bookProvider,
		bookCache:    bookCache,
		log:          log,
	}
}

func (s *BookService) AddBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	const op = "BookService.AddBook"

	log := s.log.With(
		slog.String("op", op),
		slog.String("id", book.ID),
	)

	book, err := s.bookSaver.AddBook(ctx, book)
	if err != nil {
		log.Error("failed AddBook", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("added book")
	return book, nil
}
func (s *BookService) UpdateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	const op = "BookService.UpdateBook"

	log := s.log.With(
		slog.String("op", op),
		slog.String("id", book.ID),
	)

	updatedBook, err := s.bookSaver.UpdateBook(ctx, book)
	if err != nil {
		log.Error("failed to update book", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	cacheKey := fmt.Sprintf("book:%s", book.ID)
	if err := s.bookCache.InvalidateBook(ctx, cacheKey); err != nil {
		log.Warn("failed to invalidate cache", slog.String("error", err.Error()))
	}

	log.Info("book updated successfully")
	return updatedBook, nil
}
func (s *BookService) DeleteBook(ctx context.Context, id string) (string, error) {
	const op = "BookService.DeleteBook"

	log := s.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)

	id, err := s.bookSaver.DeleteBook(ctx, id)
	if err != nil {
		log.Error("failed to delete book", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("book deleted successfully")
	return id, nil
}
func (s *BookService) GetBook(ctx context.Context, id string) (*models.Book, error) {
	const op = "BookService.GetBook"

	log := s.log.With(
		slog.String("op", op),
		slog.String("id", id),
	)
	cacheKey := fmt.Sprintf("book:%s", id)
	cachedBook, err := s.bookCache.GetBook(ctx, cacheKey)
	if err != nil {
		log.Warn("cache get error", slog.String("error", err.Error()))
	}
	if cachedBook != nil {
		log.Debug("book retrieved from cache")
		return cachedBook, nil
	}

	book, err := s.bookProvider.GetBook(ctx, id)
	if err != nil {
		log.Error("failed to get book", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.bookCache.SetBook(ctx, cacheKey, book); err != nil {
		log.Warn("failed to cache book", slog.String("error", err.Error()))
	}
	log.Debug("book retrieved")
	return book, nil
}
func (s *BookService) ListBooks(ctx context.Context, filter *models.BookFilter) ([]*models.Book, error) {
	const op = "BookService.ListBooks"

	log := s.log.With(
		slog.String("op", op),
	)

	books, err := s.bookProvider.ListBooks(ctx, filter)
	if err != nil {
		log.Error("failed to list books", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("listed books", slog.Int("count", len(books)))
	return books, nil
}
func (s *BookService) GetUserBooks(ctx context.Context, userID string, filter *models.BookFilter) ([]*models.Book, error) {
	const op = "BookService.GetUserBooks"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	books, err := s.bookProvider.GetUserBooks(ctx, userID, filter)
	if err != nil {
		log.Error("failed to get user books", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("retrieved user books", slog.Int("count", len(books)))
	return books, nil
}
func (s *BookService) AddBookToUser(ctx context.Context, userID, bookID string) (string, error) {
	const op = "BookService.AddBookToUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", userID))

	savedBookID, err := s.bookSaver.AddBookToUser(ctx, userID, bookID)

	if err != nil {
		log.Error("failed to add book to user", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("added book to user", slog.String("savedBookID", savedBookID))
	return savedBookID, nil
}
func (s *BookService) RemoveBookFromUser(ctx context.Context, userID, bookID string) (string, error) {
	const op = "BookService.RemoveBookFromUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", userID))

	deletedBookId, err := s.bookSaver.RemoveBookFromUser(ctx, userID, bookID)

	if err != nil {
		log.Error("failed to remove book from user", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("removed book from user", slog.String("deletedBookId", deletedBookId))
	return deletedBookId, nil
}
