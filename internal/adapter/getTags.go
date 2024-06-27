package adapter

import (
	"context"
	"sync"

	"github.com/SeqqueryLab/s3util/internal/application/core/domain"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GetObjectTags extracts tagging of the object with key key, at bucket buckt
func (a *Adapter) GetObjectTags(bucket, key string) (map[string]string, error) {
	// Define result
	result := make(map[string]string)
	res, err := a.client.GetObjectTagging(
		context.TODO(),
		&s3.GetObjectTaggingInput{
			Bucket: &bucket,
			Key:    &key,
		},
	)
	if err != nil {
		return nil, err
	}
	// Loop over S3.TagSet and fillout the result
	for _, tag := range res.TagSet {
		result[*tag.Key] = *tag.Value
	}

	return result, nil
}

// GetAllTags returns list of tagsets for each object in the bucket
func (a *Adapter) GetAllTags(bucket string) ([]map[string]string, error) {
	// define result
	var result []map[string]string
	// Define wait group
	var wg sync.WaitGroup
	// Define output channel
	out := make(chan map[string]string)

	// Get all objects in the bucket
	objects, err := a.GetObjectsBucket(bucket)
	if err != nil {
		return nil, err
	}
	// iterate over objects and return tagging
	for _, o := range *objects {
		wg.Add(1)
		go func(object domain.Object) {
			defer wg.Done()
			tag, _ := a.GetObjectTags(bucket, object.Key)
			if len(tag) == 0 {
				return
			}
			out <- tag
		}(o)
	}

	// Wait
	go func() {
		wg.Wait()
		close(out)
	}()
	// Iterate over channel and collect results
	for tagset := range out {
		result = append(result, tagset)
	}

	return result, nil
}
