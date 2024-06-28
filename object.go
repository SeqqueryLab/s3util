package s3util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	PartSize = 50_000_000
)

type partUploadResult struct {
	completedPart *s3.UploadPartOutput
	err           error
}

// WriteJson
// Writes JSON object in bucket with given id
func (s *Service) WriteJson(bucket string, key string, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewBuffer(b),
	})

	log.Printf("PutObject result %+v", string(b))

	return err
}

// ReadJson DONE
// Reads json file from the storage
func (s *Service) ReadJson(bucket string, key string) ([]byte, error) {

	// Get the object from S3
	res, err := s.client.GetObject(
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

// GetObject
// Reads object from S3 storage with provided bucket, and key
func (s *Service) GetObject(bucket, key string) (io.Reader, error) {
	res, err := s.client.GetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		},
	)
	if err != nil {
		return nil, err
	}
	reader := res.Body

	return reader, err
}

// DeleteObject DONE
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
func (s *Service) ListObjectsBucket(bucket, prefix string) ([]types.Object, error) {
	var contents []types.Object
	res, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	})
	if err != nil {
		return nil, err
	}
	contents = res.Contents

	return contents, err
}

// ListObjectDir
// Lists objects in source directory
func (s *Service) ListObjectDir(bucket, source string) ([]types.Object, error) {
	// prepare the prefix
	prefix := path.Clean(source)
	prefix = fmt.Sprintf("%s/", prefix)

	// list objects
	res, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	})
	if err != nil {
		return nil, err
	}
	contents := res.Contents
	if len(contents) == 0 {
		return nil, fmt.Errorf("directory %s does not exist", source)
	}

	return contents, nil
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

// UploadObjecct DONE
// UploadObject to the bucket
func (s *Service) UploadObject(bucket, key string, r io.Reader, partMiB int64) error {
	// Clean key
	key = path.Clean(key)
	// Prepare uploader
	uploader := manager.NewUploader(s.client, func(u *manager.Uploader) {
		u.PartSize = partMiB * 1024 * 1024
	})
	// Fail on error
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   r,
	})
	return err
}
