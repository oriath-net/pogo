package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/oriath-net/pogo/ggpk"
)

func openGgpkPath(path string) (io.Reader, error) {
	colonIdx := strings.LastIndex(path, ":")
	if colonIdx < 0 {
		return os.Open(path)
	}

	gf, err := ggpk.Open(path[0:colonIdx])
	if err != nil {
		return nil, err
	}

	return gf.Open(path[colonIdx+1:])
}
