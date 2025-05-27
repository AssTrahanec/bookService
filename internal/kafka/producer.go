package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
)

type BookCreatedEvent struct {
	BookID string `json:"bookId"`
	Title  string `json:"title"`
}

var writer *kafka.Writer

func InitKafkaProducer(brokers []string, topic string) {
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
	})
	slog.Info("kafka producer initialized", slog.String("topic", topic))
}

func PublishBookCreatedEvent(bookID string, title string) error {
	event := BookCreatedEvent{
		BookID: bookID,
		Title:  title,
	}

	data, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal book-created event", slog.String("error", err.Error()))
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", bookID)),
		Value: data,
	}

	err = writer.WriteMessages(context.Background(), msg)
	if err != nil {
		slog.Error("failed to send kafka message", slog.String("error", err.Error()))
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	slog.Info("published book-created event", slog.String("bookId", bookID))
	return nil
}
