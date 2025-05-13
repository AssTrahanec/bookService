package book_service

import (
	gen "bookService/internal/delivery/protos/gen/go"
	"context"
	"google.golang.org/grpc"
)

type serverAPI struct {
	gen.UnimplementedBookServiceServer
}

func Register(gRPC *grpc.Server) {
	gen.RegisterBookServiceServer(gRPC, &serverAPI{})
}

func (s *serverAPI) AddBook(
	ctx context.Context,
	req *gen.AddBookRequest,
) (*gen.Book, error) {
	panic("implement me")
}
func (s *serverAPI) GetBook(
	ctx context.Context,
	req *gen.GetBookRequest,
) (*gen.Book, error) {
	panic("implement me")
}

func (s *serverAPI) UpdateBook(
	ctx context.Context,
	req *gen.UpdateBookRequest,
) (*gen.Book, error) {
	panic("implement me")
}

func (s *serverAPI) DeleteBook(
	ctx context.Context,
	req *gen.DeleteBookRequest,
) (*gen.DeleteBookResponse, error) {
	panic("implement me")
}

func (s *serverAPI) ListBooks(
	ctx context.Context,
	req *gen.ListBooksRequest,
) (*gen.ListBooksResponse, error) {
	panic("implement me")
}

func (s *serverAPI) AddBookToUser(
	ctx context.Context,
	req *gen.UserBookRequest,
) (*gen.AddUserBookResponse, error) {
	panic("implement me")
}

func (s *serverAPI) RemoveBookFromUser(
	ctx context.Context,
	req *gen.UserBookRequest,
) (*gen.RemoveBookFromUserResponse, error) {
	panic("implement me")
}
func (s *serverAPI) GetUserBooks(
	ctx context.Context,
	req *gen.GetUserBooksRequest,
) (*gen.ListBooksResponse, error) {
	panic("implement me")
}
