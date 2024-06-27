package api

import (
	"github.com/SeqqueryLab/s3util/internal/adapter"
	"github.com/SeqqueryLab/s3util/internal/ports/client"
)

// Application holds the core of app functionality
type Application struct {
	client.ClientPort
}

// NewApplication
func NewApplication() (Application, error) {
	client, err := adapter.NewAdapter()
	if err != nil {
		return Application{}, err
	}

	return Application{client}, nil
}

/*
// GetBuckets
func (a *Application) GetBuckets() ([]domain.Bucket, error) {
	res, err := a.GetBuckets()
	return res, err
}

// GetTags
func (a *Application) GetTags(bucket, key string) (map[string]string, error) {
	res, err := a.GetObjectTags(bucket, key)
	return res, err
}

// PutTags
func (a *Application) PutTags(bucket, key string, tags map[string]string) error {
	err := a.PutObjectTags(bucket, key, tags)
	return err
}

// GetAllTags
func (a *Application) GetAllTags(bucket string) ([]map[string]string, error) {
	res, err := a.GetAllTags(bucket)
	return res, err
}
*/
