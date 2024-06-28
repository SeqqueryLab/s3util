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
