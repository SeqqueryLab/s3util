package adapter

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"sync"
)

func (a *Adapter) DirGet(bucket, source string) (io.Reader, error) {
	// define wait group
	var wg sync.WaitGroup
	// define error channel
	ch := make(chan error)
	// return error if bucket does not exist
	_, err := a.BucketExists(bucket)
	if err != nil {
		return nil, fmt.Errorf("bucket %s does not exist: %w", bucket, err)
	}
	// return error if dir does not exist
	_, err = a.DirExists(bucket, source)
	if err != nil {
		return nil, fmt.Errorf("directory %s does not exist: %w", source, err)
	}
	// list objects in directory
	res, err := a.DirListObjects(bucket, source)
	if err != nil {
		return nil, fmt.Errorf("error listing the contents of the directory %s: %w", source, err)
	}
	// return error if directory is empty
	if len(res) == 0 {
		return nil, fmt.Errorf("directory %s is emtpy", source)
	}
	// create new buffer to write an archive
	buf := new(bytes.Buffer)
	// create writer
	w := zip.NewWriter(buf)
	// iterate over the objects and write them to the buffer
	for _, item := range res {
		wg.Add(1)
		go func(bucket, key string) {
			defer wg.Done()
			obj, err := a.ObjectGet(bucket, key)
			if err != nil {
				ch <- err
				return
			}
			f, err := w.Create(key)
			if err != nil {
				ch <- err
				return
			}
			body, err := io.ReadAll(obj)
			if err != nil {
				ch <- err
				return
			}
			_, err = f.Write(body)
			if err != nil {
				ch <- err
				return
			}

		}(bucket, *item.Key)
	}

	// wat goroutins to complete
	go func() {
		wg.Wait()
		close(ch)
	}()

	// check for errors, return if any
	for err := range ch {
		if err != nil {
			_ = w.Close()
			return nil, fmt.Errorf("failed to compress directory %s: %w", source, err)
		}
	}

	// check for errors before closing the writer
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to compress directory %s: %w", source, err)
	}

	return buf, nil
}
