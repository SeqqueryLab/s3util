package client

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/SeqqueryLab/s3util/internal/application/core/domain"
)

type ClientPort interface {
	GetBuckets() ([]domain.Bucket, error)
	DirListObjects(bucket, source string) ([]types.Object, error)
	DirDelete(bucket, source string) error
	JsonWrite(bucket, key string, body interface{}) error
	JsonRead(bucket, key string) ([]byte, error)
	ObjectDelete(bucket, key string) error
	GetObjectTags(bucket, key string) (map[string]string, error)
	PutObjectTags(bucket, key string, tags map[string]string) error
	GetAllTags(bucket string) ([]map[string]string, error)
}
