package poefs

import (
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"sort"
	"time"
)

var errIsDir = errors.New("cannot read a directory")

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
	return &unionFSMergedDirectoryStat{path.Base(udir.path)}, nil
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

	sort.SliceStable(dirents, func(i, j int) bool {
		return dirents[i].Name() < dirents[j].Name()
	})

	uniqDirents := make([]fs.DirEntry, 0, len(dirents))
	for i := 0; i < len(dirents); i++ {
		if i+1 < len(dirents) && dirents[i].Name() == dirents[i+1].Name() {
			de1, de2 := dirents[i], dirents[i+1]
			i += 1
			if de1.IsDir() && de2.IsDir() {
				uniqDirents = append(uniqDirents, &unionFSMergedDirectoryStat{de1.Name()})
			} else {
				log.Fatalf("Directory %s contains multiple non-unifiable files named %s", udir.path, de1.Name())
			}
			continue
		}
		uniqDirents = append(uniqDirents, dirents[i])
	}

	dirents = uniqDirents

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
	name string
}

func (uds *unionFSMergedDirectoryStat) Name() string {
	return uds.name
}

func (uds *unionFSMergedDirectoryStat) Size() int64 {
	return 0
}

func (uds *unionFSMergedDirectoryStat) Mode() fs.FileMode {
	return 0o444 | fs.ModeDir
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

func (uds *unionFSMergedDirectoryStat) Provenance() string {
	return "union"
}

func (uds *unionFSMergedDirectoryStat) Signature() []byte {
	return nil
}

func (uds *unionFSMergedDirectoryStat) Type() fs.FileMode {
	return uds.Mode()
}

func (uds *unionFSMergedDirectoryStat) Info() (fs.FileInfo, error) {
	return uds, nil
}
