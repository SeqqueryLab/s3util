package s3util

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

// ListBucket function lists all the buckets associated with the account
func (s *Service) GetBuckets() ([]types.Bucket, error) {
	var buckets []types.Bucket

	res, err := s.client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	buckets = res.Buckets
	return buckets, err
}

// BucketExists returns true if bucket exists, and false otherwise. Returns
// error if request failed
func (s *Service) BucketExists(id string) (bool, error) {
	res := true
	_, err := s.client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
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

// CreateBucket
// Creates bucket with bucket id
func (s *Service) CreateBucket(id string) error {
	_, err := s.client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(id),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(s.client.Options().Region),
		},
	})
	return err
}

// DeleteBucket
// Deletes bucket with bucket id
func (s *Service) DeleteBucket(id string) error {
	_, err := s.client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(id),
	})
	if err != nil {
		return err
	}
	return nil
}
