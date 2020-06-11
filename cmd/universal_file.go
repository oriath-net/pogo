package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/oriath-net/pogo/ggpk"
)

func openGgpkPath(path string) (io.Reader, error) {
	parts := strings.SplitN(path, ":", 2)
	switch len(parts) {
	case 1:
		return os.Open(path)

	case 2:
		gf, err := ggpk.Open(parts[0])
		if err != nil {
			return nil, err
		}
		return gf.Open(parts[1])

	default:
		panic("how?")
	}
}
