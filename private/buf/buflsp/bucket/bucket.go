package bucket

import (
	"errors"
	"fmt"

	"github.com/bufbuild/buf/private/pkg/storage"
	"github.com/bufbuild/buf/private/pkg/storage/storageos"
)

var bucketProvider = storageos.NewProvider(storageos.ProviderWithSymlinks())

type HasRoot interface {
	Root() string
}

func GetBucketRoot(b storage.ReadWriteBucket) (string, bool) {
	if rb, ok := b.(HasRoot); ok {
		return rb.Root(), true
	} else {
		return "", false
	}
}

type buckethasroot struct {
	storage.ReadWriteBucket
	root string
}

func NewBucket(root string) (storage.ReadWriteBucket, error) {
	if root == "" {
		return nil, errors.New("`root` cannot be empty")
	}

	b, err := bucketProvider.NewReadWriteBucket(root)
	if err != nil {
		return nil, fmt.Errorf("bucket provider new: %w", err)
	}

	return &buckethasroot{
		root:            root,
		ReadWriteBucket: b,
	}, nil
}

func (b *buckethasroot) Root() string {
	return b.root
}
