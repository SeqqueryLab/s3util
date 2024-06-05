package s3util

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Dir struct {
	Name     string
	Size     int64
	Modified time.Time
}

func ObjectToDir(o types.Object) *Dir {
	return &Dir{
		Name:     path.Base(*o.Key),
		Size:     *o.Size,
		Modified: *o.LastModified,
	}
}

func (s *Service) ListDir(bucket, source string) ([]Dir, error) {
	// prepare the prefix
	prefix := path.Clean(source)
	prefix = fmt.Sprintf("%s/", prefix)

	res, err := s.client.ListObjectsV2(
		context.TODO(),
		&s3.ListObjectsV2Input{
			Bucket: &bucket,
			Prefix: &prefix,
		},
	)
	if err != nil {
		return nil, err
	}

	var content []types.Object = res.Contents
	if len(content) == 0 {
		return nil, errors.New("directory does not exist")
	}

	dirs := make(map[string]Dir)

	for _, val := range content {
		key := strings.TrimPrefix(*val.Key, prefix)
		key = regexp.MustCompile(`^[^/]+`).FindString(key)

		if key != "" {
			d1 := ObjectToDir(val)
			d1.Name = key
			if d2, ok := dirs[key]; ok {
				d1.Size += d2.Size
				if d1.Modified.Before(d2.Modified) {
					d1.Modified = d2.Modified
				}
			}

			dirs[key] = *d1
		}
	}

	var result []Dir
	for _, val := range dirs {
		result = append(result, val)
	}
	return result, nil
}

func (s *Service) CopyDir(bucket, source, destination string) error {
	source = path.Clean(source)
	destination = path.Clean(destination)
	out := make(chan error)
	defer close(out)

	res, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(source),
	})
	if err != nil {
		return err
	}

	var content []types.Object = res.Contents
	if len(content) == 0 {
		return errors.New("directory does not exist")
	}

	var wg sync.WaitGroup

	for _, val := range content {
		wg.Add(1)
		oldKey := val.Key
		// DOES NOt WORK PROPERLY, ADD BLOCKING
		go func(bucket, key, source, destination string) {
			defer wg.Done()
			newKey := fmt.Sprintf("%s/%s", destination, strings.TrimPrefix(key, source))
			newKey = path.Clean(newKey)

			_, err = s.client.CopyObject(context.TODO(), &s3.CopyObjectInput{
				Bucket:     aws.String(bucket),
				CopySource: aws.String(fmt.Sprintf("%s/%s", bucket, key)),
				Key:        aws.String(newKey),
			})
			if err != nil {
				out <- err
				return
			}
		}(bucket, *oldKey, source, destination)
	}
	wg.Wait()

	return nil
}

// MoveDir
// Moves directory with it content to the new location
func (s *Service) MoveDir(bucket, source, destination string) error {
	source = path.Clean(source)
	destination = path.Clean(destination)
	out := make(chan error)
	defer close(out)

	res, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(source),
	})
	if err != nil {
		return err
	}

	var content []types.Object = res.Contents
	if len(content) == 0 {
		return errors.New("directory does not exist")
	}

	var wg sync.WaitGroup

	for _, val := range content {
		wg.Add(1)
		oldKey := val.Key
		// DOES NOt WORK PROPERLY, ADD BLOCKING
		go func(bucket, key, source, destination string) {
			defer wg.Done()
			newKey := fmt.Sprintf("%s/%s", destination, strings.TrimPrefix(key, source))
			newKey = path.Clean(newKey)

			_, err = s.client.CopyObject(context.TODO(), &s3.CopyObjectInput{
				Bucket:     aws.String(bucket),
				CopySource: aws.String(fmt.Sprintf("%s/%s", bucket, key)),
				Key:        aws.String(newKey),
			})
			if err != nil {
				// error is not catch in this function
				return
			}

			err = s.DeleteObject(bucket, key)
			if err != nil {
				// error is not catch in this funcction
				return
			}

		}(bucket, *oldKey, source, destination)
	}
	wg.Wait()

	return nil
}

// DeleteDir
// Deletes the directory with all it's content
func (s *Service) DeleteDir(bucket, source string) error {
	source = path.Clean(source)
	res, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(source),
	})
	if err != nil {
		return err
	}

	var content []types.Object = res.Contents
	if len(content) == 0 {
		return errors.New("directory does not exist")
	}

	var wg sync.WaitGroup

	for _, val := range content {
		wg.Add(1)
		key := *val.Key
		log.Printf("Deleting the object %s\n", key)

		// DOES NOt WORK PROPERLY, ADD BLOCKING
		go func(bucket, key string) {
			defer wg.Done()
			err = s.DeleteObject(bucket, key)
			if err != nil {
				// error is not catch in this funcction
				return
			}

		}(bucket, key)
	}
	wg.Wait()

	return err
}
