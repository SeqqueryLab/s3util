package adapter

import (
	"context"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// ListObjectDir
// Lists objects in source directory
func (a *Adapter) DirListObjects(bucket, source string) ([]types.Object, error) {
	// cliean the prefix the prefix
	prefix := fmt.Sprintf("%s/", path.Clean(source))
	// list objects
	res, err := a.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	})
	if err != nil {
		return nil, err
	}
	// Get the contents, and check if it is not empty
	contents := res.Contents
	if len(contents) == 0 {
		return nil, fmt.Errorf("directory %s does not exist", source)
	}

	return contents, nil
}
