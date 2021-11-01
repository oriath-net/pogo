package poefs

import (
	"io/fs"
	"os"
)

type OpenCache interface {
	Open(p string) (fs.FS, error)
	OpenFile(unipath string) (fs.File, error)
}

type openCache map[string]fs.FS

func NewOpenCache() OpenCache {
	return openCache(make(map[string]fs.FS))
}

func (oc openCache) Open(p string) (fs.FS, error) {
	f, ok := oc[p]
	if ok {
		return f, nil
	}

	f, err := Open(p)
	if err != nil {
		return nil, err
	}

	oc[p] = f
	return f, nil
}

func (oc openCache) OpenFile(unipath string) (fs.File, error) {
	srcPath, localPath := SplitPath(unipath)

	if srcPath == "" {
		return os.Open(localPath)
	}

	f, err := oc.Open(srcPath)
	if err != nil {
		return nil, err
	}
	return f.Open(localPath)
}
