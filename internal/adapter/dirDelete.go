package adapter

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// DirDelete Deletes the directory with all it's content
func (a *Adapter) DirDelete(bucket, source string) error {
	// Define waitgroup
	var wg sync.WaitGroup
	// Define output channel
	ch := make(chan error)
	// List the objects in the directory
	res, err := a.DirListObjects(bucket, source)
	if err != nil {
		return err
	}
	// Iterate over the directory content, and delete objects
	for _, v := range res {
		wg.Add(1)
		go func(o types.Object) {
			defer wg.Done()
			ch <- a.ObjectDelete(bucket, *v.Key)
		}(v)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for err = range ch {
		if err != nil {
			return fmt.Errorf("failed to delete the object: %s", err)
		}
	}

	return nil
}
