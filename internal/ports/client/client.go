package client

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/SeqqueryLab/s3util/internal/application/core/domain"
)

type ClientPort interface {
	GetBuckets() ([]domain.Bucket, error)
	BucketExists(id string) (bool, error)
	CreateBucket(id string) error
	DeleteBucket(id string) error
	DirExists(bucket, source string) (bool, error)
	DirListObjects(bucket, source string) ([]types.Object, error)
	DirDelete(string, string) (*domain.Directory, error)
	JsonWrite(bucket, key string, body interface{}) error
	JsonRead(bucket, key string) ([]byte, error)
	ObjectGet(bucket, key string) (io.Reader, error)
	ObjectUpload(bucket, key string, body io.Reader, partMiB int64) error
	ObjectDelete(bucket, key string) error
	GetObjectTags(bucket, key string) (map[string]string, error)
	PutObjectTags(bucket, key string, tags map[string]string) error
	GetAllTags(bucket string) ([]map[string]string, error)
}
