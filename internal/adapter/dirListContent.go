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
		return nil, fmt.Errorf("directory does not exist: %s", err)
	}
	// Get the contents of the directory
	contents := res.Contents

	return contents, nil
}
