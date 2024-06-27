package adapter

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Adapter struct {
	client *s3.Client
}

func NewAdapter() (*Adapter, error) {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(config, func(o *s3.Options) {
		endpoint := os.Getenv("AWS_ENDPOINT")
		o.BaseEndpoint = &endpoint
	})

	return &Adapter{client: client}, nil
}
