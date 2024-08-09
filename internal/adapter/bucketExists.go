package adapter

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

// BucketExists returns true if bucket exists, and false otherwise. Returns
// error if request failed
func (a *Adapter) BucketExists(id string) (bool, error) {
	res := true
	_, err := a.client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(id),
	})
	if err != nil {
		var awsError smithy.APIError
		if errors.As(err, &awsError) {
			switch awsError.(type) {
			case *types.NotFound:
				res, err = false, nil
			default:
				res, err = false, errors.New("access denied")
			}
		}
	}

	return res, err
}
