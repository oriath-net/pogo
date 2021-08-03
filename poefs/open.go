package poefs

import (
	"archive/zip"
	"errors"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/oriath-net/pogo/poefs/bundle"
	"github.com/oriath-net/pogo/poefs/ggpk"
)

var (
	errBadPath = errors.New("invalid source path")
)

func Open(p string) (fs.FS, error) {
	fi, err := os.Stat(p)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		dirfs := os.DirFS(p)

		bundles, err := bundle.NewLoader(dirfs)
		if err != nil {
			return nil, err
		}
		return newUnionFS(dirfs, bundles), nil
	}

	ext := path.Ext(p)

	if ext == ".ggpk" {
		f, err := os.Open(p)
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

	if ext == ".zip" {
		return zip.OpenReader(p)
	}

	return nil, errBadPath
}

// Split a path into two parts on a colon. If no colon is present (or if it's
// part of a Windows drive prefix), return an empty string and the input.
func SplitPath(p string) (string, string) {
	colon := strings.LastIndex(p, ":")

	// no colon in path
	if colon < 0 {
		return "", p
	}

	// there was only one colon, and it's in a Windows drive prefix
	if colon == 1 && len(p) > colon+1 && p[colon+1] == '\\' {
		return "", p
	}

	srcPath := p[:colon]
	localPath := p[colon+1:]

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
