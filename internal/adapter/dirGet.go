package adapter

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"path"
	"strings"
)

func (a *Adapter) DirGet(bucket, source string) (io.Reader, error) {
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
		// get the object key
		key := *item.Key
		// get object from the store
		obj, err := a.ObjectGet(bucket, key)
		if err != nil {
			return nil, err
		}
		// remove leading prefix for the directory
		prefix := fmt.Sprintf("%s/", path.Dir(source))
		fname, _ := strings.CutPrefix(key, prefix)
		// create a file with fname in the zip archive
		f, err := w.Create(fname)
		if err != nil {
			return nil, err
		}
		// get the body of the object
		body, err := io.ReadAll(obj)
		if err != nil {
			return nil, err
		}
		// write the body to the zip archive
		_, err = f.Write(body)
		if err != nil {
			return nil, err
		}
	}

	// check for errors before closing the writer
	if err := w.Close(); err != nil {
		log.Printf("failed to close writer: %s", err)
		return nil, fmt.Errorf("failed to compress directory %s: %w", source, err)
	}

	return buf, nil
}
