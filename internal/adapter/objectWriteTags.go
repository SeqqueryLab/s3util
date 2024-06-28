package adapter

import (
	"context"
	"errors"
	"path"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// WriteObjectTags updates the tagset of the object on S3
func (a *Adapter) PutObjectTags(bucket, key string, tags map[string]string) error {
	// Define result
	var tagset []types.Tag
	// Clean key
	key = path.Clean(key)
	// If tags is not empty
	if len(tags) == 0 {
		return errors.New("tags are empty: at least one key-value pair is required")
	}
	// Retreive old tags of the object if available
	tagCopy, err := a.GetObjectTags(bucket, key) // TODO check for object existance here
	if err != nil {
		tagCopy = make(map[string]string)
	}

	// Iterate over key-value pairs converting into the []types.Tag
	for key, val := range tags {
		tag := &types.Tag{Key: &key, Value: &val}
		tagset = append(tagset, *tag)
		// delete old tag if present
		delete(tagCopy, key)
	}

	// Preserve old tags which are not modified
	if len(tagCopy) > 0 {
		for key, val := range tagCopy {
			tag := &types.Tag{Key: &key, Value: &val}
			tagset = append(tagset, *tag)
		}
	}

	// Tag object on S3
	_, err = a.client.PutObjectTagging(
		context.TODO(),
		&s3.PutObjectTaggingInput{
			Bucket: &bucket,
			Key:    &key,
			Tagging: &types.Tagging{
				TagSet: tagset,
			},
		},
	)

	return err
}
