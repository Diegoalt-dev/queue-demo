package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Options struct {
	QueueURL string
	BaseURL  string
}

type Client struct {
	options Options
	client  *sqs.Client
}

func New(cfg aws.Config, options Options) *Client {
	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		if options.QueueURL != "" {
			o.BaseEndpoint = aws.String(options.BaseURL)
		}
	})

	return &Client{
		client:  client,
		options: options,
	}
}

func (c *Client) SendMessage(ctx context.Context, message string) (*sqs.SendMessageOutput, error) {
	return c.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(c.options.QueueURL),
		MessageBody: aws.String(message),
	})
}

func (c *Client) ReceiveMessage(ctx context.Context) (*sqs.ReceiveMessageOutput, error) {
	return c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.options.QueueURL),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
	})
}

func (c *Client) DeleteMessage(ctx context.Context, message types.Message) (*sqs.DeleteMessageOutput, error) {
	return c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.options.QueueURL),
		ReceiptHandle: message.ReceiptHandle,
	})
}
