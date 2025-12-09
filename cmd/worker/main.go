package main

import (
	"context"
	"demoproject/internal/platform/sqs"
	serviceprocessor "demoproject/internal/processor"
	processor "demoproject/internal/taskprocessor"
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		return
	}
}

func main() {
	queueURL := os.Getenv("QUEUE_URL")
	baseURL := os.Getenv("BASE_URL")
	if queueURL == "" {
		log.Fatal("QUEUE_URL environment variable not set")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("eu-central-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("x", "x", "")))
	if err != nil {
		log.Fatal("error getting config")
	}

	sqsClient := sqs.New(cfg, sqs.Options{BaseURL: baseURL, QueueURL: queueURL})

	s := serviceprocessor.NewProcessor()
	sqsProcessor := processor.NewTaskProcessor(sqsClient, s)

	slog.Info("started worker")
	sqsProcessor.Process(context.Background())

}
