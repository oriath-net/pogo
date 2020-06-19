package cmd

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	cli "github.com/urfave/cli/v2"
)

var RiverDl = cli.Command{
	Name:      "riverdl",
	Usage:     "Download stash tab data for analysis",
	UsageText: "pogo riverdl [options]",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "id",
			Usage:    "Change ID to start at",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "into",
			Aliases: []string{"o"},
			Usage:   "download into this directory",
			Value:   "out",
		},
		&cli.BoolFlag{
			Name:  "continuous",
			Usage: "Follow next_change_id",
		},
	},

	Action: do_riverdl,
}

func do_riverdl(c *cli.Context) error {
	id := c.String("id")
	continuous := c.Bool("continuous")
	out := c.String("into")

	err := os.MkdirAll(out, 0777)
	if err != nil {
		return fmt.Errorf("mkdir %s: %s", out, err)
	}

	maxRate := 1000 * time.Millisecond

	for {
		reqStart := time.Now()

		nextId, err := riverDownloadId(id, out)
		if err != nil {
			return err
		}
		if !continuous {
			break
		}
		id = nextId

		// Rate control
		reqLength := time.Since(reqStart)
		delay := maxRate - reqLength
		if delay > 0 {
			fmt.Printf(
				"request took %d ms - sleeping %d ms\n",
				reqLength/time.Millisecond,
				delay/time.Millisecond,
			)
			time.Sleep(delay)
		}
	}

	return nil
}

func riverDownloadId(id string, outpath string) (string, error) {
	filename := filepath.Join(outpath, fmt.Sprintf("river-%s.json.gz", id))

	fmt.Printf("%s: waiting", filename)
	resp, err := http.Get(fmt.Sprintf(
		"http://api.pathofexile.com/public-stash-tabs?id=%s",
		id,
	))
	if err == nil && resp.StatusCode != 200 {
		err = errors.New("Request returned HTTP " + resp.Status)
	}
	if err != nil {
		fmt.Printf("\r\x1b[K%s: failed\n", filename)
		return "", err
	}

	progress := int(0)
	data, err := ioutil.ReadAll(
		&progressReader{resp.Body, func(n int) {
			progress += n
			fmt.Printf("\r\x1b[K%s: %d", filename, progress)
		}},
	)
	if err != nil {
		fmt.Printf("\r\x1b[K%s: failed\n", filename)
		return "", err
	}
	fmt.Printf("\r\x1b[K%s: complete\n", filename)

	cidStruct := struct {
		NextChangeID string `json:"next_change_id"`
	}{}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&cidStruct)
	if err != nil {
		return "", fmt.Errorf("unable to get next change ID: %w", err)
	}

	outputFh, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", fmt.Errorf("unable to open %s: %w", filename, err)
	}
	defer outputFh.Close()

	outputGz, err := gzip.NewWriterLevel(outputFh, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	defer outputGz.Close()

	_, err = outputGz.Write(data)
	if err != nil {
		return "", fmt.Errorf("unable to write to %s: %w", filename, err)
	}
	outputGz.Flush()

	return cidStruct.NextChangeID, nil
}

// as suggested by https://stackoverflow.com/a/26050400
type progressReader struct {
	io.Reader
	Reporter func(n int)
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.Reporter(n)
	return n, err
}
