package app

import (
	"bookService/config"
	grpcapp "bookService/internal/app/grpc"
	bookService "bookService/internal/services/bookService"
	"bookService/internal/storage/postres"
	"bookService/internal/storage/redis"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	config *config.Config,
) *App {
	storage, err := postres.New(config.DB)
	cache, err := redis.New(config.Cache)
	if err != nil {
		panic(err)
	}
	libraryService := bookService.New(storage, storage, cache, log)
	grpcApp := grpcapp.New(log, grpcPort, libraryService)
	return &App{
		GRPCSrv: grpcApp,
	}
}
