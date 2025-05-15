package auth

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	adminRole = "admin"
	userRole  = "user_role"
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

	if isAdminMethod(info.FullMethod) && roles[0] != adminRole {
		return nil, status.Error(codes.PermissionDenied, "admin role required")
	}

	newCtx := context.WithValue(ctx, userRole, roles[0])

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
		if method == m {
			return true
		}
	}
	return false
}
