package kafka

import (
	"context"
	"encoding/json"

	"github.com/RozmiDan/gameReviewHub/internal/config"
	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	logger *zap.Logger
}

// NewProducer создаёт нового продьюсера по конфигу.
func NewProducer(cfg *config.KafkaConfig, logger *zap.Logger) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      cfg.Brokers,
		Topic:        cfg.TopicRatings,
		RequiredAcks: int(kafka.RequireOne),
		Async:        true,
		BatchTimeout: 0,
	})

	logger = logger.With(zap.String("component", "kafka-producer"))
	return &Producer{writer: writer, logger: logger}
}

// PublishRating публикует сообщение с оценкой в Kafka.
func (p *Producer) PublishRating(ctx context.Context, msg entity.RatingMessage) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		p.logger.Error("failed to marshal rating message", zap.Error(err))
		return err
	}

	kmsg := kafka.Message{
		Key:   []byte(msg.GameID),
		Value: bytes,
	}

	if err := p.writer.WriteMessages(ctx, kmsg); err != nil {
		p.logger.Error("failed to write message to kafka", zap.Error(err),
			zap.String("topic", p.writer.Topic))
		return err
	}

	p.logger.Info("published rating to kafka",
		zap.String("game_id", msg.GameID),
		zap.String("user_id", msg.UserID),
		zap.Int32("rating", msg.Rating),
	)
	return nil
}

// Close закрывает продьюсера.
func (p *Producer) Close() error {
	return p.writer.Close()
}
