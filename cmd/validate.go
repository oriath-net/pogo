package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/oriath-net/pogo/dat"
	"github.com/oriath-net/pogo/poefs"

	cli "github.com/urfave/cli/v2"
)

var Validate = cli.Command{
	Name:      "validate",
	Usage:     "Perform cross-version validation on data formats.",
	UsageText: "pogo validate <source path pattern> <data file names...>",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "report",
			Usage: "output status to a Markdown file",
		},
		&cli.StringSliceFlag{
			Name:  "version",
			Usage: "specify versions to test against",
		},
		&cli.StringFlag{
			Name:  "fmt",
			Usage: "path to a directory containing formats",
		},
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "display details about failures",
		},
	},

	Action: do_validate,
}

type validateStatus byte

const (
	statusNA      validateStatus = iota
	statusMissing                // before OK to support missing->OK transitions
	statusOk
	statusWarn
	statusFail
)

func (s validateStatus) String() string {
	switch s {
	case statusNA:
		return "n/a"
	case statusMissing:
		return "missing"
	case statusOk:
		return "ok"
	case statusWarn:
		return "WARN"
	case statusFail:
		return "FAIL"
	default:
		return "???"
	}
}

func (s validateStatus) Glyph() string {
	switch s {
	case statusNA:
		return "n/a"
	case statusMissing:
		return "\u2753" // question mark
	case statusOk:
		return "\u2705" // white check in green box
	case statusWarn:
		return "\u26a0\ufe0f" // warning sign + emoji variation selector
	case statusFail:
		return "\u274c" // red X
	default:
		return "???"
	}
}

func (s *validateStatus) update(to validateStatus) {
	if to > *s {
		*s = to
	}
}

type validateRunner struct {
	opencache poefs.OpenCache
	pathspec  string
	fmtDir    string
	versions  []string
	verbose   bool
}

type validateOutput struct {
	oldestSeen, newestSeen  string
	current, dat64, history validateStatus
}

func (ro validateOutput) VersionRange(r validateRunner) string {
	if ro.newestSeen != r.versions[len(r.versions)-1] {
		return ro.oldestSeen + "-" + ro.newestSeen
	} else if ro.oldestSeen != r.versions[0] {
		return ro.oldestSeen + "-"
	} else {
		return "all"
	}
}

func (ro validateOutput) Log(r validateRunner, filename string) {
	log.Printf(
		"%-32s [%-8s] - structure %-8s dat64 %-8s history %-8s\n",
		filename,
		ro.VersionRange(r),
		ro.current.String(),
		ro.dat64.String(),
		ro.history.String(),
	)
}

//go:embed validate.header.md
var reportHeader string

func do_validate(c *cli.Context) error {
	if c.NArg() < 2 {
		return errNotEnoughArguments
	}

	r := validateRunner{
		opencache: poefs.NewOpenCache(),
		pathspec:  c.Args().Get(0),
		fmtDir:    c.String("fmt"),
		versions:  c.StringSlice("version"),
		verbose:   c.Bool("verbose"),
	}

	if !strings.Contains(r.pathspec, "%v") {
		return fmt.Errorf("pathspec must contain a %%v version placeholder")
	}

	if len(r.versions) == 0 {
		r.versions = []string{
			"0.11", "1.0", "1.1", "1.2", "1.3",
			"2.0", "2.1", "2.2", "2.3", "2.4", "2.5", "2.6",
			"3.0", "3.1", "3.2", "3.3", "3.4", "3.5", "3.6", "3.7", "3.8", "3.9",
			"3.10", "3.11", "3.12", "3.13", "3.14", "3.15", "3.16", "3.17",
		}
	}

	reportLetter := byte(0)
	reportFile := io.Writer(nil)
	if reportTo := c.String("report"); reportTo != "" {
		fd, err := os.Create(reportTo)
		if err != nil {
			return err
		}
		reportFile = fd
		fmt.Fprint(reportFile, reportHeader)
	}

	for i := 1; i < c.NArg(); i++ {
		filename := c.Args().Get(i)

		result, err := r.Validate(filename)
		if err != nil {
			return err
		}

		result.Log(r, filename)

		if reportFile != nil {
			if filename[0] != reportLetter {
				reportLetter = filename[0]
				fmt.Fprintf(reportFile, "\n## %c\n\n", reportLetter)
				fmt.Fprintf(reportFile, "| File                             | Releases | current  | dat64 | history\n")
				fmt.Fprintf(reportFile, "| -------------------------------- | -------- | -------- | ----- | --------\n")
			}
			fmt.Fprintf(reportFile,
				"| %-32s | %-8s | %s | %s | %s\n",
				filename,
				result.VersionRange(r),
				result.current.Glyph(),
				result.dat64.Glyph(),
				result.history.Glyph(),
			)
		}
	}

	return nil
}

func (r *validateRunner) Validate(filename string) (validateOutput, error) {
	out := validateOutput{
		current: statusMissing,
		dat64:   statusNA,
		history: statusNA,
	}

	for i := len(r.versions) - 1; i >= 0; i-- {
		v := r.versions[i]
		basename := strings.Replace(r.pathspec, "%v", v, -1) + filename
		err := r.validateVersion(&out, v, filename, basename)
		if err != nil {
			return validateOutput{}, err
		}
	}

	return out, nil
}

func (r *validateRunner) validateVersion(out *validateOutput, version string, filename, basename string) error {
	p := dat.InitParser(version) // normal parser
	q := dat.InitParser(version) // strict parser
	if r.fmtDir != "" {
		p.SetFormatDir(r.fmtDir)
		q.SetFormatDir(r.fmtDir)
	}
	q.SetStrict(1)

	f32, err := r.opencache.OpenFile(basename + ".dat")
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist in this version at all
			return nil
		} else {
			return err
		}
	}
	defer f32.Close()

	// since we've successfully opened the .dat, we can mark the file as existing
	// note that this assumes versions are ordered newest -> oldest
	isCurrentVersion := false
	if out.newestSeen == "" {
		out.newestSeen = version
		isCurrentVersion = true
	}
	out.oldestSeen = version

	// open the dat64 if it exists
	f64, err := r.opencache.OpenFile(basename + ".dat64")
	if err != nil {
		if os.IsNotExist(err) {
			f64 = nil
		} else {
			return err
		}
	}
	if f64 != nil {
		defer f64.Close()
		out.dat64.update(statusMissing) // N/A -> missing
	}

	// we'll need to read the data twice, so store it in a buffer
	buf := bytes.Buffer{}
	buf.ReadFrom(f32)

	// dat parsing
	rows32, err := p.Parse(bytes.NewReader(buf.Bytes()), filename+".dat")
	if err != nil {
		if r.verbose {
			log.Printf("%s.dat %s: parsing failed: %s", filename, version, err)
		}
		if os.IsNotExist(err) {
			out.current.update(statusMissing)
			out.history.update(statusMissing)
			return nil
		}
		if isCurrentVersion {
			out.current.update(statusFail)
			if f64 != nil {
				out.dat64.update(statusFail) // dat64 is implicitly bad because we can't compare it
			}
		} else {
			out.history.update(statusFail) // at least one version has now failed
		}
		return nil
	} else {
		if isCurrentVersion {
			out.current.update(statusOk)
		} else {
			out.history.update(statusOk)
		}
		// normal parsing worked, now try strict parsing!
		_, err := q.Parse(bytes.NewReader(buf.Bytes()), filename+".dat")
		if err != nil {
			if r.verbose {
				log.Printf("%s.dat %s: strict parsing failed: %s", filename, version, err)
			}
			out.history.update(statusWarn)
			if isCurrentVersion {
				out.current.update(statusWarn)
			}
		}
	}

	if f64 == nil {
		// no dat64, we're done
		return nil
	}

	buf.Reset()
	buf.ReadFrom(f64)

	rows64, err := p.Parse(bytes.NewReader(buf.Bytes()), filename+".dat64")
	if err != nil {
		if r.verbose {
			log.Printf("%s.dat64 %s: parsing failed: %s", filename, version, err)
		}
		out.dat64.update(statusFail)
		return nil
	} else {
		out.dat64.update(statusOk)
		// try strict parsing
		_, err := q.Parse(bytes.NewReader(buf.Bytes()), filename+".dat64")
		if err != nil {
			if r.verbose {
				log.Printf("%s.dat64 %s: strict parsing failed: %s", filename, version, err)
			}
			out.dat64.update(statusWarn)
		}

		if !reflect.DeepEqual(rows32, rows64) {
			if r.verbose {
				log.Printf("%s.dat64 %s: data not equal to %s.dat", filename, version, filename)
			}
			out.dat64.update(statusFail)
		}
	}

	return nil
}
