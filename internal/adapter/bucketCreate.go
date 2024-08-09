package adapter

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// CreateBucket
// Creates bucket with bucket id
func (a *Adapter) CreateBucket(id string) error {
	_, err := a.client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(id),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(a.client.Options().Region),
		},
	})
	return err
}
