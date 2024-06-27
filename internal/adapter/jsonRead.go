package adapter

import (
	"context"
	"io"
	"path"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// JsonRead Reads json file from the storage
func (a *Adapter) JsonRead(bucket string, key string) ([]byte, error) {
	// Clean the key
	key = path.Clean(key)
	// Get the object from S3
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

	// read the response body to tye bytes buffer
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// close the response body
	res.Body.Close()

	return b, nil
}
