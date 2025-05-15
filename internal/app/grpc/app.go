package grpcapp

import (
	"bookService/internal/auth"
	bookServicegrpc "bookService/internal/grpc/book-service"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"time"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	port int,
	bookService bookServicegrpc.BookService,
) *App {
	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(auth.AuthInterceptor))

	bookServicegrpc.Register(gRPCServer, bookService)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server started", slog.String("addr", lis.Addr().String()))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
func (a *App) Stop() { //chat gpt
	const op = "grpcapp.Stop"
	const timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	a.log.With(slog.String("op", op)).
		Info("grpc server stopping", slog.Int("port", a.port))

	stopped := make(chan struct{})
	go func() {
		a.gRPCServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		a.log.Info("grpc server stopped gracefully")
	case <-ctx.Done():
		a.log.Warn("forcing grpc server shutdown due to timeout")
		a.gRPCServer.Stop()
	}
}
