package bundle

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"sort"
	"strings"
	"time"
)

type bundleFS struct {
	lower fs.FS
	index bundleIndex
}

func NewLoader(lower fs.FS) (*bundleFS, error) {
	indexFile, err := lower.Open("Bundles2/_.index.bin")
	if err != nil {
		return nil, err
	}

	// FIXME: It'd be neat to defer this until it's needed.
	idx, err := loadBundleIndex(indexFile.(io.ReaderAt))
	if err != nil {
		return nil, err
	}

	return &bundleFS{
		lower: lower,
		index: idx,
	}, nil
}

func (b *bundleFS) Open(name string) (fs.File, error) {
	files := b.index.files

	// super special case
	if name == "." {
		return &bundleFsDir{
			fs:     b,
			prefix: "",
			offset: 0,
		}, nil
	}

	// binary search for the file
	idx := sort.Search(len(b.index.files), func(i int) bool {
		return files[i].path >= name
	})

	if idx < len(files) && files[idx].path == name {
		return &bundleFsFile{
			fs:   b,
			info: &files[idx],
		}, nil
	}

	// check for a directory separately -- it's possible for a file to have a
	// prefix which masks the directory, e.g. Abc/Def.txt and Abc/Def/Ghi.txt,
	// searching for Abc/Def will return the file first
	dirName := name + "/"
	idx += sort.Search(len(b.index.files)-idx, func(i int) bool {
		return files[idx+i].path >= dirName
	})

	if idx < len(files) && strings.HasPrefix(files[idx].path, dirName) {
		return &bundleFsDir{
			fs:     b,
			prefix: dirName,
			offset: idx,
		}, nil
	}

	// nope, nothing here
	return nil, &fs.PathError{
		Op:   "open",
		Path: name,
		Err:  fs.ErrNotExist,
	}
}

// bundleFsFile

type bundleFsFile struct {
	fs     *bundleFS
	info   *bundleFileInfo
	reader *io.SectionReader
}

func (bff *bundleFsFile) initReader() error {
	if bff.reader != nil {
		return nil
	}

	bundlePath := "Bundles2/" + bff.fs.index.bundles[bff.info.bundleId] + ".bundle.bin"
	bundleFile, err := bff.fs.lower.Open(bundlePath)
	if err != nil {
		return &fs.PathError{
			Op:   "open",
			Path: bff.info.path,
			Err:  fmt.Errorf("unable to open bundle %s: %w", bundlePath, err),
		}
	}

	bundle, err := openBundle(bundleFile.(io.ReaderAt))
	if err != nil {
		return &fs.PathError{
			Op:   "open",
			Path: bff.info.path,
			Err:  fmt.Errorf("unable to load bundle %s: %w", bundlePath, err),
		}
	}

	bff.reader = io.NewSectionReader(
		bundle,
		int64(bff.info.offset),
		int64(bff.info.size),
	)

	return nil
}

func (bff *bundleFsFile) Read(p []byte) (int, error) {
	err := bff.initReader()
	if err != nil {
		return 0, err
	}
	return bff.reader.Read(p)
}

func (bff *bundleFsFile) ReadAt(p []byte, off int64) (int, error) {
	err := bff.initReader()
	if err != nil {
		return 0, err
	}
	return bff.reader.ReadAt(p, off)
}

func (bff *bundleFsFile) Seek(offset int64, whence int) (int64, error) {
	err := bff.initReader()
	if err != nil {
		return 0, err
	}
	return bff.reader.Seek(offset, whence)
}

func (bff *bundleFsFile) Close() error {
	return nil
}

func (bff *bundleFsFile) Stat() (fs.FileInfo, error) {
	return &bundleFsFileInfo{bff}, nil
}

// bundleFsFileInfo

type bundleFsFileInfo struct {
	*bundleFsFile
}

func (bffi bundleFsFileInfo) Name() string {
	return path.Base(bffi.info.path)
}

func (bffi bundleFsFileInfo) Size() int64 {
	return int64(bffi.info.size)
}

func (bffi bundleFsFileInfo) Mode() fs.FileMode {
	return 0o444
}

func (bffi bundleFsFileInfo) ModTime() time.Time {
	return time.Unix(0, 0)
}

func (bffi bundleFsFileInfo) IsDir() bool {
	return false
}

func (bffi bundleFsFileInfo) Sys() interface{} {
	return nil
}

func (bffi bundleFsFileInfo) Provenance() string {
	return fmt.Sprintf("bundle %s + %x", bffi.fs.index.bundles[bffi.info.bundleId], bffi.info.offset)
}

func (bffi bundleFsFileInfo) Signature() []byte {
	return nil
}

// bundleFsDir

type bundleFsDir struct {
	fs     *bundleFS
	prefix string
	offset int
}

func (bfd *bundleFsDir) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("is a directory")
}

func (bfd *bundleFsDir) Close() error {
	return nil
}

func (bfd *bundleFsDir) Stat() (fs.FileInfo, error) {
	return &bundleFsDirInfo{bfd}, nil
}

func (bfd *bundleFsDir) ReadDir(n int) ([]fs.DirEntry, error) {
	files := bfd.fs.index.files
	prefixLen := len(bfd.prefix)

	dirents := []fs.DirEntry{}

	for {
		fi := &files[bfd.offset]
		if !strings.HasPrefix(fi.path, bfd.prefix) {
			break
		}

		slashIdx := strings.Index(fi.path[prefixLen:], "/")
		if slashIdx != -1 {
			dir := fi.path[:prefixLen+slashIdx]
			dirents = append(dirents, &bundleFsDirEnt{
				fs:   bfd.fs,
				path: dir,
			})
			next := bfd.offset + sort.Search(len(files)-bfd.offset, func(i int) bool {
				return files[bfd.offset+i].path >= dir+"/\xff"
			})
			bfd.offset = next

		} else {
			dirents = append(dirents, &bundleFsDirEnt{
				fs:   bfd.fs,
				path: fi.path,
				file: &bundleFsFile{
					fs:     bfd.fs,
					info:   fi,
					reader: nil, // not needed here
				},
			})
			bfd.offset += 1
		}

		if bfd.offset >= len(files) {
			break
		}

		if n > 0 && len(dirents) >= n {
			return dirents, nil
		}
	}

	if n > 0 {
		return dirents, io.EOF
	}

	return dirents, nil
}

// bundleFsDirInfo

type bundleFsDirInfo struct {
	*bundleFsDir
}

func (bfdi bundleFsDirInfo) Name() string {
	return path.Base(bfdi.prefix)
}

func (bfdi bundleFsDirInfo) Size() int64 {
	return 0
}

func (bfdi bundleFsDirInfo) Mode() fs.FileMode {
	return 0o444 | fs.ModeDir
}

func (bfdi bundleFsDirInfo) ModTime() time.Time {
	return time.Unix(0, 0)
}

func (bfdi bundleFsDirInfo) IsDir() bool {
	return true
}

func (bfdi bundleFsDirInfo) Sys() interface{} {
	return nil
}

func (bfdi bundleFsDirInfo) Provenance() string {
	return "bundle dir"
}

func (bfdi bundleFsDirInfo) Signature() []byte {
	return nil
}

// bundleFsDirEnt

type bundleFsDirEnt struct {
	fs   *bundleFS
	path string
	file *bundleFsFile
}

func (bfde *bundleFsDirEnt) Name() string {
	return path.Base(bfde.path)
}

func (bfde *bundleFsDirEnt) IsDir() bool {
	return bfde.file == nil
}

func (bfde *bundleFsDirEnt) Type() fs.FileMode {
	if bfde.IsDir() {
		return 0o444 | fs.ModeDir
	} else {
		return 0o444
	}
}

func (bfde *bundleFsDirEnt) Info() (fs.FileInfo, error) {
	if bfde.IsDir() {
		return &bundleFsDirInfo{
			&bundleFsDir{
				fs:     bfde.fs,
				prefix: bfde.path,
				offset: -1, // unused
			},
		}, nil
	} else {
		return &bundleFsFileInfo{bfde.file}, nil
	}
}
