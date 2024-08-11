package adapter

import (
	"context"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// BucketExists returns true if bucket exists, and false otherwise. Returns
// error if request failed
func (a *Adapter) DirExists(bucket, source string) (bool, error) {
	// return false, error if bucket does not exist
	if ok, err := a.BucketExists(bucket); !ok {
		return false, err
	}
	// cliean the prefix the prefix
	prefix := fmt.Sprintf("%s/", path.Clean(source))
	// list objects
	res, err := a.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	})
	if err != nil {
		return false, fmt.Errorf("directory does not exist: %s", err)
	}
	// Get the contents, and check if it is not empty
	contents := res.Contents
	if len(contents) == 0 {
		return false, fmt.Errorf("directory %s does not exist", source)
	}
	return true, nil
}
