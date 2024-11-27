package bucket

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/bufbuild/buf/private/pkg/storage"
)

var rootPath string = "/"
var rootbucket *buckethasroot
var rootbucketonce sync.Once

func init() {
	if erp := os.Getenv("BUF_LSP_ROOT_PATH"); erp != "" {
		rootPath = erp
	} else if runtime.GOOS == "windows" {
		rootPath = filepath.VolumeName("C:\\") + "\\"
	}
}

func RootBucket() storage.ReadWriteBucket {
	rootbucketonce.Do(func() {
		b, err := NewBucket(rootPath)
		if err != nil {
			panic(fmt.Errorf("init root bucket: %w", err))
		}

		rootbucket = b.(*buckethasroot)
	})

	return rootbucket
}
