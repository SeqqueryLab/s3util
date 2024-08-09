package domain

import (
	"path"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Bucket
type Bucket struct {
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

// S3BucketToBucket
func S3BucketToBucket(b *types.Bucket) *Bucket {
	return &Bucket{
		Name:    *b.Name,
		Created: *b.CreationDate,
	}
}

// Directory
type Directory struct {
	Key      string    `json:"key"`
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
}

// Object an object on S3 store
type Object struct {
	Name     string            `json:"name"`
	Key      string            `json:"key"`
	Size     int64             `json:"size"`
	Modified time.Time         `json:"modified"`
	Tags     map[string]string `json:"tags"`
}

// S3ToObject converts s3 object to Object
func S3ToObject(o types.Object) *Object {
	return &Object{
		Name:     path.Base(*o.Key),
		Key:      *o.Key,
		Size:     *o.Size,
		Modified: *o.LastModified,
	}
}

// S3ObjectsToObject
func S3ObjectsToObject(o []types.Object) *[]Object {
	// Define result
	var result []Object
	// define wait group
	var wg sync.WaitGroup
	// Define output channel
	out := make(chan Object)
	// Loop over and make a conversion
	for _, object := range o {
		wg.Add(1)
		go func(object types.Object) {
			defer wg.Done()
			converted := S3ToObject(object)
			out <- *converted
		}(object)
	}

	// Wait
	go func() {
		wg.Wait()
		close(out)
	}()

	// collect results
	for val := range out {
		result = append(result, val)
	}

	return &result
}

// Tags
type Tags map[string]string
