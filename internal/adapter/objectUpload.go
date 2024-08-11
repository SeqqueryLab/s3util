package adapter

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ObjectUpload uploads object to the bucket
func (a *Adapter) ObjectUpload(bucket, key string, r io.Reader, partMiB int64) error {
	// Clean key
	key = path.Clean(key)
	// Prepare uploader
	uploader := manager.NewUploader(a.client, func(u *manager.Uploader) {
		u.PartSize = partMiB * 1024 * 1024
	})
	// read data to the buffer
	buff, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("can not read file: %s", err)
	}
	body := bytes.NewReader(buff)
	// Fail on error
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   body,
	})
	return err
}
