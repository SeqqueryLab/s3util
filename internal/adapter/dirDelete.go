package adapter

import (
	"fmt"
	"path"
	"sync"

	"github.com/SeqqueryLab/s3util/internal/application/core/domain"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// DirDelete Deletes the directory with all it's content
func (a *Adapter) DirDelete(bucket, source string) (*domain.Directory, error) {
	// Define result
	dir := &domain.Directory{}
	// Define waitgroup
	var wg sync.WaitGroup
	// Define lock
	var lock sync.Mutex
	// Define output channel
	ch := make(chan error)
	// List the objects in the directory
	res, err := a.DirListObjects(bucket, source)
	if err != nil {
		return nil, fmt.Errorf("failed to delete directory: %s", err)
	}
	// Iterate over the directory content, and delete objects
	for _, v := range res {
		wg.Add(1)
		go func(o types.Object) {
			defer wg.Done()
			defer lock.Unlock()
			err := a.ObjectDelete(bucket, *v.Key)
			if err != nil {
				ch <- err
				return
			}
			lock.Lock()
			dir.Size += *v.Size
			if v.LastModified.After(dir.Modified) {
				dir.Modified = *v.LastModified
			}
		}(v)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	// Update directory data
	dir.Key = source
	dir.Name = path.Base(source)

	// Check for errors and return first error if any
	for err = range ch {
		if err != nil {
			return dir, fmt.Errorf("failed to delete the object. first error: %s", err)
		}
	}

	// return results
	return dir, nil
}
