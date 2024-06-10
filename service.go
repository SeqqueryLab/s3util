package s3util

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Service embeds aws S3 client providing S3 API
type Service struct {
	client *s3.Client
}

func NewFromEnv() *Service {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	cl := s3.NewFromConfig(config, func(o *s3.Options) {
		endpoint := os.Getenv("AWS_ENDPOINT")
		o.BaseEndpoint = &endpoint
	})

	return &Service{cl}
}

// NewFromConfig
// Returns new S3 client from user-provided config file
func NewFromConfig(id, key, secret, session, region, endpoint string, retry int, accelerate bool) *Service {
	config := NewConfig(id, key, secret, session, region, endpoint, retry, accelerate)

	cl := s3.New(*config)

	return &Service{cl}
}
