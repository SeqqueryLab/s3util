package s3util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
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

// CreateObject
// Creates object in bucket with given id
func (s *Service) CreateJson(bucket string, key string, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	length := int64(len(b))

	res, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        &bucket,
		Key:           &key,
		Body:          bytes.NewBuffer(b),
		ContentLength: &length,
	}, s3.WithAPIOptions(
		v4.SwapComputePayloadSHA256ForUnsignedPayloadMiddleware,
	))

	log.Printf("PutObject result %+v", res)

	return err
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
func (s *Service) ListObjectsBucket(bucket, prefix string) ([]types.Object, error) {
	res, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
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

// UploadObjecct
// UploadObject to the bucket
func (s *Service) UploadObject(bucket, key string, data []byte) error {
	buff := bytes.NewReader(data)
	var parMib int64

	uploader := manager.NewUploader(s.client, func(u *manager.Uploader) {
		u.PartSize = parMib * 1024 * 1024
	})

	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   buff,
	})
	return err
}

// uploadPart
func (s *Service) uploadPart(fileBytes []byte, bucket, key, uploadId string, n int32, ch chan *partUploadResult, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Uploading part %d, size %d ", n, len(fileBytes))

	uploadRes, err := s.client.UploadPart(context.TODO(), &s3.UploadPartInput{
		Body:          bytes.NewReader(fileBytes),
		Bucket:        &bucket,
		Key:           &key,
		PartNumber:    &n,
		UploadId:      &uploadId,
		ContentLength: aws.Int64(int64(len(fileBytes))),
	})
	if err != nil {
		log.Printf("Error uploading data: %s", err)
		ch <- &partUploadResult{
			nil,
			err,
		}
		return
	} else {
		ch <- &partUploadResult{
			&s3.UploadPartOutput{
				ETag: uploadRes.ETag,
			},
			nil,
		}
	}
}

// Multipart upload
// Multipart upload to s3
func (s *Service) UploadObjectMultipart(bucket, key string) error {
	var (
		wg                            sync.WaitGroup
		completed, current, remaining int
		n                             int32 = 1
	)

	ch := make(chan *partUploadResult, 5)

	file, _ := os.Open("example.fastq.gz")
	defer file.Close()

	stat, _ := file.Stat()
	log.Printf("Read the file. Size: %d\n B", stat.Size())
	remaining = int(stat.Size())

	buff := make([]byte, stat.Size())
	_, _ = file.Read(buff)

	id, err := s.client.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return err
	}

	uploaded := func() (int, error) {
		res, err := s.client.ListParts(context.TODO(), &s3.ListPartsInput{
			Bucket:   &bucket,
			Key:      &key,
			UploadId: id.UploadId,
		})
		if err != nil {
			return 0, err
		}

		parts := res.Parts
		var size int
		for _, part := range parts {
			size += int(*part.Size)
		}
		fmt.Printf("Uploaded: %v", res)
		return size, err
	}

	for start := 0; remaining > 0; start += PartSize {
		log.Printf("Uploading part %d, remaining %d\n", n, remaining)
		wg.Add(1)
		if remaining < PartSize {
			current = remaining
		} else {
			current = PartSize
		}
		go s.uploadPart(buff[start:current+start], bucket, key, *id.UploadId, n, ch, &wg)
		remaining = -current
		n++
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for res := range ch {
		if res.err != nil {
			// abort
			_, err := s.client.AbortMultipartUpload(context.TODO(), &s3.AbortMultipartUploadInput{
				Bucket:   &bucket,
				Key:      &key,
				UploadId: id.UploadId,
			})
			if err != nil {
				log.Fatal(err)
			}
		} else {
			temp, err := uploaded()
			if err != nil {
				log.Fatal(err)
			}
			completed += temp
		}
		fmt.Printf("Uploaded: %d B (%d %%)\n", completed, completed/int(stat.Size())*100)
	}

	s.client.CompleteMultipartUpload(context.TODO(), &s3.CompleteMultipartUploadInput{
		Bucket:   &bucket,
		Key:      &key,
		UploadId: id.UploadId,
	})
	return nil
}

// SelectObjectContent (json, csv)
// Given an SQL query select object content and return it to the user
func (s *Service) SelectObjectContent(bucket, key, query string) ([]types.Object, error) {
	// send the request
	res, err := s.client.SelectObjectContent(context.TODO(), &s3.SelectObjectContentInput{
		Bucket:         aws.String(bucket),
		Key:            aws.String(key),
		ExpressionType: types.ExpressionTypeSql,
		Expression:     aws.String("SELECT name FROM S3Object WHERE cast(age as int) > 35"),
		InputSerialization: &types.InputSerialization{
			CSV: &types.CSVInput{
				FileHeaderInfo: types.FileHeaderInfoUse,
			},
		},
		OutputSerialization: &types.OutputSerialization{
			CSV: &types.CSVOutput{},
		},
	})
	if err != nil {
		return nil, err
	}

	log.Printf("%+v", res)

	return nil, err
}
