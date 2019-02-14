package aengine

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/pkg/errors"
	"io/ioutil"
)

type GcsFileStore struct {
	Bucket string
	Prefix string
}

func (fs *GcsFileStore) Load(ctx context.Context, path string) ([]byte, error) {
	storageService, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	rc, err := storageService.Bucket(fs.Bucket).Object(fs.Prefix + path).NewReader(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return data, nil
}

func (fs *GcsFileStore) Save(ctx context.Context, path string, content []byte) error {
	storageService, err := storage.NewClient(ctx)
	if err != nil {
		return nil
	}

	w := storageService.Bucket(fs.Bucket).Object(fs.Prefix + path).NewWriter(ctx)
	if err != nil {
		return nil
	}

	_, err = w.Write(content)
	if err != nil {
		w.Close()
		return errors.Wrap(err, "")
	}

	err = w.Close()
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
