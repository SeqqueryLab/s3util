package adapter

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// DeleteBucket
// Deletes bucket with bucket id
func (a *Adapter) DeleteBucket(id string) error {
	_, err := a.client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(id),
	})
	if err != nil {
		return err
	}
	return nil
}
