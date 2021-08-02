package poefs

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"time"
)

var (
	errIsDir = errors.New("cannot read a directory")
)

type unionFS struct {
	members []fs.FS
}

func newUnionFS(members ...fs.FS) unionFS {
	return unionFS{members}
}

func (f unionFS) Open(path string) (fs.File, error) {
	dirs := []fs.File{}
	for _, m := range f.members {
		f, err := m.Open(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}

		st, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if st.IsDir() {
			dirs = append(dirs, f)
		} else {
			return f, nil
		}
	}
	if len(dirs) > 0 {
		return &unionFSMergedDirectory{path, dirs, 0}, nil
	}
	return nil, &fs.PathError{
		Op:   "open",
		Path: path,
		Err:  fs.ErrNotExist,
	}
}

type unionFSMergedDirectory struct {
	path    string
	members []fs.File
	offset  int
}

func (udir *unionFSMergedDirectory) Read([]byte) (int, error) {
	return 0, errIsDir
}

func (udir *unionFSMergedDirectory) Stat() (fs.FileInfo, error) {
	return &unionFSMergedDirectoryStat{udir}, nil
}

func (udir *unionFSMergedDirectory) Close() error {
	for _, m := range udir.members {
		err := m.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (udir *unionFSMergedDirectory) ReadDir(n int) ([]fs.DirEntry, error) {
	dirents := []fs.DirEntry{}

	for _, m := range udir.members {
		mrdr, ok := m.(fs.ReadDirFile)
		if !ok {
			continue
		}
		ents, err := mrdr.ReadDir(0)
		if err != nil {
			return nil, err
		}
		dirents = append(dirents, ents...)
	}

	if n <= 0 {
		udir.offset = 0
		return dirents, nil
	}

	dirents = dirents[udir.offset:]
	if len(dirents) > n {
		udir.offset += n
		return dirents[:n], nil
	} else {
		udir.offset += len(dirents)
		return dirents, io.EOF
	}
}

type unionFSMergedDirectoryStat struct {
	*unionFSMergedDirectory
}

func (uds *unionFSMergedDirectoryStat) Name() string {
	return uds.path
}

func (uds *unionFSMergedDirectoryStat) Size() int64 {
	return 0
}

func (uds *unionFSMergedDirectoryStat) Mode() fs.FileMode {
	return 0444 | fs.ModeDir
}

func (uds *unionFSMergedDirectoryStat) ModTime() time.Time {
	return time.Unix(0, 0)
}

func (uds *unionFSMergedDirectoryStat) IsDir() bool {
	return true
}

func (uds *unionFSMergedDirectoryStat) Sys() interface{} {
	return nil
}
