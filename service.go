package s3util

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

// Service embeds aws S3 client providing S3 API
type Service struct {
	client *s3.Client
}

func New() *Service {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	cl := s3.NewFromConfig(config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("https://s3-eu-central-1.ionoscloud.com")
	})

	return &Service{cl}
}
