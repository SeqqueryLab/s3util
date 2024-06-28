package adapter

import (
	"context"
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
	// Fail on error
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   r,
	})
	return err
}
