package poefs

import (
	"errors"
	"io/fs"
	"os"
	"strings"

	"github.com/oriath-net/pogo/poefs/bundle"
	"github.com/oriath-net/pogo/poefs/ggpk"
)

var (
	errBadPath = errors.New("invalid source path")
)

func Open(path string) (fs.FS, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.Mode().IsRegular() {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		ggfs, err := ggpk.NewLoader(f)
		if err != nil {
			return nil, err
		}

		if !ggfs.UseBundles() {
			return ggfs, nil
		}

		bundles, err := bundle.NewLoader(ggfs)
		if err != nil {
			return nil, err
		}
		return newUnionFS(ggfs, bundles), nil
	}

	if fi.IsDir() {
		dirfs := os.DirFS(path)

		bundles, err := bundle.NewLoader(dirfs)
		if err != nil {
			return nil, err
		}
		return newUnionFS(dirfs, bundles), nil
	}

	return nil, errBadPath
}

// Split a path into two parts on a colon. If no colon is present (or if it's
// part of a Windows drive prefix), return an empty string and the input.
func SplitPath(path string) (string, string) {
	colon := strings.LastIndex(path, ":")

	// no colon in path
	if colon < 0 {
		return "", path
	}

	// there was only one colon, and it's in a Windows drive prefix
	if colon == 1 && len(path) > colon+1 && path[colon+1] == '\\' {
		return "", path
	}

	srcPath := path[:colon]
	localPath := path[colon+1:]

	localPath = strings.Trim(localPath, "/")
	if localPath == "" {
		localPath = "."
	}

	return srcPath, localPath
}

func OpenFile(unipath string) (fs.File, error) {
	srcPath, localPath := SplitPath(unipath)

	if srcPath == "" {
		return os.Open(localPath)
	}

	src, err := Open(srcPath)
	if err != nil {
		return nil, err
	}
	return src.Open(localPath)
}
