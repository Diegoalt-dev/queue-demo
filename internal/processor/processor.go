package serviceprocessor

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Processor struct {
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Process(ctx context.Context, message types.Message) error {
	slog.Info("processing message", slog.String("id", *message.MessageId))
	return nil
}
