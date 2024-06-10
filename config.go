package s3util

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewCredentials(key, secret, session string) *aws.CredentialsCache {
	return aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(key, secret, session))
}

func NewConfig(id, key, secret, session, region, endpoint string, retry int, accelerate bool) *s3.Options {
	// clean endpoint
	//endpoint = path.Clean(endpoint)

	return &s3.Options{
		Credentials:      NewCredentials(key, secret, session),
		Region:           region,
		BaseEndpoint:     &endpoint,
		AppID:            id,
		RetryMaxAttempts: retry,
		UseAccelerate:    accelerate,
	}
}
