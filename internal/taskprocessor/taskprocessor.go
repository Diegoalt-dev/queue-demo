package taskprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

type Task struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

type SQSClient interface {
	SendMessage(ctx context.Context, message string) (*sqs.SendMessageOutput, error)
	ReceiveMessage(ctx context.Context) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, message types.Message) (*sqs.DeleteMessageOutput, error)
}

type Processor interface {
	Process(ctx context.Context, message types.Message) error
}

type TaskProcessor struct {
	client    SQSClient
	processor Processor
}

func NewTaskProcessor(client SQSClient, processor Processor) *TaskProcessor {
	return &TaskProcessor{
		client:    client,
		processor: processor,
	}
}

func (s *TaskProcessor) Send(ctx context.Context, payload string) error {
	body, err := json.Marshal(Task{
		ID:      uuid.NewString(),
		Payload: payload,
	})

	if err != nil {
		return fmt.Errorf("failed to serialize task: %w", err)
	}

	output, err := s.client.SendMessage(ctx, string(body))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	jsonOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonOutput))

	return nil
}

func worker(id int, funcs <-chan func() error) {
	for f := range funcs {
		fmt.Println("worker", id)
		_ = f()
	}
}

func (s *TaskProcessor) Process(ctx context.Context) {

	funcs := make(chan func() error, 20)

	for w := 1; w <= 4; w++ {
		go worker(w, funcs)
	}

	for {
		msgs, err := s.client.ReceiveMessage(ctx)
		if err != nil {
			slog.Error("error processing", slog.Any("error", err))
			continue
		}

		for _, m := range msgs.Messages {
			funcs <- func() error {
				err = s.processor.Process(ctx, m)
				if err != nil {
					return err
				}
				_, err = s.client.DeleteMessage(ctx, m)
				if err != nil {
					slog.Error("error deleting a message", slog.Any("error", err))
					return err
				}

				slog.Info("processed and deleted message", slog.String("id", *m.MessageId))
				return nil
			}
		}
	}

}
