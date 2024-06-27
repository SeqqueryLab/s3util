// Package s3util provides tools to work with S3 Storage
//
// # Service
//
// ## NewS3Service
//
// This function creates new instance of s3 Service which provides following functionality:
//
// - ListBuckets
package s3util

import (
	"github.com/SeqqueryLab/s3util/internal/application/core/api"
	"github.com/SeqqueryLab/s3util/internal/application/core/domain"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Rook interface
type Rook interface {
	GetBuckets() ([]domain.Bucket, error)
	DirListObjects(bucket, source string) ([]types.Object, error)
	DirDelete(bucket, source string) error
	ObjectDelete(bucket, key string) error
	GetObjectTags(bucket, key string) (map[string]string, error)
	PutObjectTags(bucket, key string, tags map[string]string) error
	GetAllTags(bucket string) ([]map[string]string, error)
}

// NewS3Service returns new S3 Service
func New() (Rook, error) {
	rook, err := api.NewApplication()
	return rook, err
}
