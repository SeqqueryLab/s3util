package adapter

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// JsonWrite write json to the bucket with the given ID
func (a *Adapter) JsonWrite(bucket string, key string, body interface{}) error {
	// parse object to json
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	// Write object to S3
	_, err = a.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewBuffer(b),
	})

	return err
}
