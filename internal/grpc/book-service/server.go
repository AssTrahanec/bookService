package book_service

import (
	gen "bookService/internal/delivery/protos/gen/go"
	"bookService/internal/domain"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookService interface {
	AddBook(ctx context.Context, book *domain.Book) (*domain.Book, error)
	GetBook(ctx context.Context, id string) (*domain.Book, error)
	UpdateBook(ctx context.Context, book *domain.Book) (*domain.Book, error)
	DeleteBook(ctx context.Context, id string) (string, error)
	ListBooks(ctx context.Context, filter *domain.BookFilter) ([]*domain.Book, error)
	AddBookToUser(ctx context.Context, userID, bookID string) (string, error)
	RemoveBookFromUser(ctx context.Context, userID, bookID string) (string, error)
	GetUserBooks(ctx context.Context, userID string, filter *domain.BookFilter) ([]*domain.Book, error)
}

type serverAPI struct {
	gen.UnimplementedBookServiceServer
	bookService BookService
}

func Register(gRPC *grpc.Server, bookService BookService) {
	gen.RegisterBookServiceServer(gRPC, &serverAPI{bookService: bookService})
}

func (s *serverAPI) AddBook(
	ctx context.Context,
	req *gen.AddBookRequest,
) (*gen.Book, error) {

	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.GetAuthor() == "" {
		return nil, status.Error(codes.InvalidArgument, "author is required")
	}
	book, err := s.bookService.AddBook(ctx, &domain.Book{
		Title:           req.Title,
		Author:          req.Author,
		PublicationYear: req.GetPublicationYear(),
		Genre:           req.GetGenre(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.Book{
		BookId:          book.ID,
		Title:           book.Title,
		Author:          book.Author,
		PublicationYear: &book.PublicationYear,
		Genre:           &book.Genre,
	}, nil
}
func (s *serverAPI) GetBook(
	ctx context.Context,
	req *gen.GetBookRequest,
) (*gen.Book, error) {

	if req.GetBookId() == "" {
		return nil, status.Error(codes.InvalidArgument, "book id is required")
	}

	book, err := s.bookService.GetBook(ctx, req.GetBookId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.Book{
		BookId:          book.ID,
		Title:           book.Title,
		Author:          book.Author,
		PublicationYear: &book.PublicationYear,
		Genre:           &book.Genre,
	}, nil
}

func (s *serverAPI) UpdateBook(
	ctx context.Context,
	req *gen.UpdateBookRequest,
) (*gen.Book, error) {
	if req.GetBookId() == "" {
		return nil, status.Error(codes.InvalidArgument, "book id is required")
	}

	book, err := s.bookService.UpdateBook(ctx, &domain.Book{
		ID:              req.GetBookId(),
		Title:           req.GetTitle(),
		Author:          req.GetAuthor(),
		PublicationYear: req.GetPublicationYear(),
		Genre:           req.GetGenre(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.Book{
		BookId:          book.ID,
		Title:           book.Title,
		Author:          book.Author,
		PublicationYear: &book.PublicationYear,
		Genre:           &book.Genre,
	}, nil
}

func (s *serverAPI) DeleteBook(
	ctx context.Context,
	req *gen.DeleteBookRequest,
) (*gen.DeleteBookResponse, error) {
	if req.GetBookId() == "" {
		return nil, status.Error(codes.InvalidArgument, "book id is required")
	}

	id, err := s.bookService.DeleteBook(ctx, req.GetBookId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.DeleteBookResponse{BookId: id}, nil
}

func (s *serverAPI) ListBooks(
	ctx context.Context,
	req *gen.ListBooksRequest,
) (*gen.ListBooksResponse, error) {
	filter := &domain.BookFilter{}
	if req.GetAuthor() != "" {
		filter.Author = req.Author
	}
	if req.GetPublicationYear() != 0 {
		filter.PublicationYear = req.PublicationYear
	}
	if req.GetGenre() != "" {
		filter.Genre = req.Genre
	}

	books, err := s.bookService.ListBooks(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &gen.ListBooksResponse{}
	for _, book := range books {
		response.Books = append(response.Books, &gen.Book{
			BookId:          book.ID,
			Title:           book.Title,
			Author:          book.Author,
			PublicationYear: &book.PublicationYear,
			Genre:           &book.Genre,
		})
	}

	return response, nil
}

func (s *serverAPI) AddBookToUser(
	ctx context.Context,
	req *gen.UserBookRequest,
) (*gen.AddUserBookResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}
	if req.GetBookId() == "" {
		return nil, status.Error(codes.InvalidArgument, "book id is required")
	}

	id, err := s.bookService.AddBookToUser(ctx, req.GetUserId(), req.GetBookId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.AddUserBookResponse{BookId: id}, nil
}

func (s *serverAPI) RemoveBookFromUser(
	ctx context.Context,
	req *gen.UserBookRequest,
) (*gen.RemoveBookFromUserResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	if req.GetBookId() == "" {
		return nil, status.Error(codes.InvalidArgument, "book id is required")
	}
	id, err := s.bookService.RemoveBookFromUser(ctx, req.GetUserId(), req.GetBookId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.RemoveBookFromUserResponse{BookId: id}, nil
}
func (s *serverAPI) GetUserBooks(
	ctx context.Context,
	req *gen.GetUserBooksRequest,
) (*gen.ListBooksResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	filter := &domain.BookFilter{}
	if req.GetAuthor() != "" {
		filter.Author = req.Author
	}
	if req.GetPublicationYear() != 0 {
		filter.PublicationYear = req.PublicationYear
	}
	if req.GetGenre() != "" {
		filter.Genre = req.Genre
	}

	books, err := s.bookService.GetUserBooks(ctx, req.GetUserId(), filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &gen.ListBooksResponse{}
	for _, book := range books {
		response.Books = append(response.Books, &gen.Book{
			BookId:          book.ID,
			Title:           book.Title,
			Author:          book.Author,
			PublicationYear: &book.PublicationYear,
			Genre:           &book.Genre,
		})
	}

	return response, nil
}
