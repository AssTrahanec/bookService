package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	AdminRole = "admin"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if isPublicMethod(info.FullMethod) {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	roles := md.Get("x-user-role")
	if len(roles) == 0 {
		return nil, status.Error(codes.PermissionDenied, "role not provided")
	}

	if isAdminMethod(info.FullMethod) && roles[0] != AdminRole {
		return nil, status.Error(codes.PermissionDenied, "admin role required")
	}

	newCtx := context.WithValue(ctx, "user_role", roles[0])

	return handler(newCtx, req)
}

func isPublicMethod(method string) bool {
	publicMethods := []string{
		"/bookService.BookService/GetBook",
		"/bookService.BookService/ListBooks",
	}
	for _, m := range publicMethods {
		if m == method {
			return true
		}
	}
	return false
}

func isAdminMethod(method string) bool {
	adminMethods := []string{
		"/bookService.BookService/AddBook",
		"/bookService.BookService/UpdateBook",
		"/bookService.BookService/DeleteBook",
	}
	for _, m := range adminMethods {
		if strings.HasPrefix(method, m) {
			return true
		}
	}
	return false
}
