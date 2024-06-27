package adapter

import (
	"context"

	"github.com/SeqqueryLab/s3util/internal/application/core/domain"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GetObjectsBucket returns list of objects in bucket
func (a *Adapter) GetObjectsBucket(bucket string) (*[]domain.Object, error) {
	// Call S3 API
	res, err := a.client.ListObjectsV2(
		context.TODO(),
		&s3.ListObjectsV2Input{Bucket: &bucket},
	)
	if err != nil {
		return nil, err
	}
	// convert to []domain.Objects
	result := domain.S3ObjectsToObject(res.Contents)

	return result, err
}
