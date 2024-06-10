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
	// Waitgroup and channel
	var wg sync.WaitGroup
	ch := make(chan Dir)

	// Map to track dirs
	dirs := make(map[string]Dir)

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
		return nil, fmt.Errorf("directory %s does not exist", source)
	}

	if len(content) == 1 && (content[0].Key == &source) {
		return nil, fmt.Errorf("%s is not a directory", source)
	}

	for _, val := range content {
		wg.Add(1)
		go func(val types.Object, prefix string) {
			log.Printf("Go routine received object: %s\n", *val.Key)
			defer wg.Done()
			key := strings.TrimPrefix(*val.Key, prefix)
			key = regexp.MustCompile(`^[^/]+`).FindString(key)
			d := ObjectToDir(val)
			d.Name = key
			ch <- *d
		}(val, prefix)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for d := range ch {
		if v, ok := dirs[d.Name]; ok {
			d.Size += v.Size
			if d.Modified.After(v.Modified) {
				d.Modified = v.Modified
			}
			dirs[d.Name] = d
		}
		dirs[d.Name] = d
	}

	var result []Dir
	for _, val := range dirs {
		result = append(result, val)
	}
	return result, nil
}

// CopyDir
// Copy the content of the source directory to the destination directory
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
// Moves the source directory with it content to the destination directory
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
	// Waitgroup and channel
	var wg sync.WaitGroup
	ch := make(chan error)

	res, err := s.ListObjectDir(bucket, source)
	if err != nil {
		return err
	}

	for _, v := range res {
		wg.Add(1)
		go func(o types.Object) {
			defer wg.Done()
			ch <- s.DeleteObject(bucket, *v.Key)
		}(v)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for err = range ch {
		if err != nil {
			return fmt.Errorf("failed to delete the object: %s", err)
		}
	}

	return nil
}
