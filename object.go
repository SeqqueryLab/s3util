package s3util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// CreateObject
// Creates object in bucket with given id
func (s *Service) CreateJson(bucket string, id string, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)

	s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(id),
		Body:   r,
	})

	return nil
}

// DeleteObject
// Deletes objects in the current bucket
func (s *Service) DeleteObject(bucket, key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

// ListObjects
// Lists object in bucket
func (s *Service) ListObjectsBucket(bucket string) ([]types.Object, error) {
	res, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	var contents []types.Object
	if err != nil {
		return nil, err
	} else {
		contents = res.Contents
	}

	return contents, err
}

// CopyObjectFolder
// Copy object to the folder
func (s *Service) CopyObjectToFolder(bucket string, key string, folder string) error {
	newKey := path.Base(key)
	_, err := s.client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(fmt.Sprintf("%s/%s", bucket, key)),
		Key:        aws.String(fmt.Sprintf("%s/%s", folder, newKey)),
	})
	return err
}

// MoveObjectToFolder
// Moves object from the original destination to another folder at the same bucket
func (s *Service) MoveObjectToFolder(bucket, key, folder string) error {
	err := s.CopyObjectToFolder(bucket, key, folder)
	if err != nil {
		return err
	}

	err = s.DeleteObject(bucket, key)
	return err
}
