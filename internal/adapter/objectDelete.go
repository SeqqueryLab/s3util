package adapter

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ObjectDelete Deletes objects in the current bucket
func (a *Adapter) ObjectDelete(bucket, key string) error {
	_, err := a.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	return err
}
