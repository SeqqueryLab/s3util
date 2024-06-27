package adapter

import (
	"context"
	"sync"

	"github.com/SeqqueryLab/s3util/internal/application/core/domain"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (a *Adapter) GetBuckets() ([]domain.Bucket, error) {
	// Define result
	var result []domain.Bucket
	// Define wait group
	var wg sync.WaitGroup
	// Define ouput channel
	out := make(chan *domain.Bucket)

	// Call S3 API
	res, err := a.client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	for _, bucket := range res.Buckets {
		wg.Add(1)
		go func(b types.Bucket) {
			defer wg.Done()
			result := domain.S3BucketToBucket(&bucket)
			out <- result
		}(bucket)
	}

	// Wait
	go func() {
		wg.Wait()
		close(out)
	}()

	// Collect results
	for bucket := range out {
		result = append(result, *bucket)
	}

	return result, nil
}
