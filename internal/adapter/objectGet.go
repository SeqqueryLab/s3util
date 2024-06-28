package adapter

import (
	"context"
	"io"
	"path"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ObjectGet Get object from s3
func (a *Adapter) ObjectGet(bucket, key string) (io.Reader, error) {
	// Clean key
	key = path.Clean(key)
	// Retreive object from S3
	res, err := a.client.GetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		},
	)
	if err != nil {
		return nil, err
	}
	// Get object body
	reader := res.Body

	return reader, err
}
