package redis

import (
	"bookService/config"
	"bookService/internal/domain/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func New(cfg config.RedisConfig) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Cache{
		client: client,
		ttl:    time.Duration(cfg.TTL) * time.Minute,
	}, nil
}

func (c *Cache) Close() error {
	return c.client.Close()
}
func (c *Cache) GetBook(ctx context.Context, key string) (*models.Book, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var book models.Book
	if err := json.Unmarshal(data, &book); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &book, nil
}

func (c *Cache) SetBook(ctx context.Context, key string, book *models.Book) error {
	data, err := json.Marshal(book)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *Cache) InvalidateBook(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
