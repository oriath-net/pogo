package validate

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/oriath-net/pogo/cmd"
	"github.com/oriath-net/pogo/dat"
	"github.com/oriath-net/pogo/poefs"
	"github.com/spf13/pflag"
)

//go:embed header.md
var reportHeader string

//go:embed formats.txt
var allFormats string

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

type validateOutput struct {
	oldestSeen, newestSeen  string
	current, dat64, history validateStatus
}

func (ro validateOutput) VersionRange(r *validateRunner) string {
	if ro.newestSeen != r.versions[len(r.versions)-1] {
		return ro.oldestSeen + "-" + ro.newestSeen
	} else if ro.oldestSeen != r.versions[0] {
		return ro.oldestSeen + "-"
	} else {
		return "all"
	}
}

func (ro validateOutput) Log(r *validateRunner, filename string) {
	log.Printf(
		"%-32s [%-8s] - structure %-8s dat64 %-8s history %-8s\n",
		filename,
		ro.VersionRange(r),
		ro.current.String(),
		ro.dat64.String(),
		ro.history.String(),
	)
}

type validateRunner struct {
	opencache poefs.OpenCache
	pathspec  string
	fmtDir    string
	reportTo  string
	versions  []string
	verbose   bool

	looseParsers, strictParsers map[string]*dat.DataParser

	reportFh     io.WriteCloser
	reportLetter byte
}

func (vr *validateRunner) Init() {
	if vr.reportTo != "" {
		fh, err := os.Create(vr.reportTo)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(fh, reportHeader)
		vr.reportFh = fh
	}

	vr.looseParsers = make(map[string]*dat.DataParser)
	vr.strictParsers = make(map[string]*dat.DataParser)

	for _, v := range vr.versions {
		lp := dat.InitParser(v)
		sp := dat.InitParser(v)

		if vr.fmtDir != "" {
			lp.SetFormatDir(vr.fmtDir)
			sp.SetFormatDir(vr.fmtDir)
		}
		sp.SetStrict(1)

		vr.looseParsers[v], vr.strictParsers[v] = lp, sp
	}
}

func (vr *validateRunner) Validate(filename string) {
	result := validateOutput{
		current: statusMissing,
		dat64:   statusNA,
		history: statusNA,
	}

	if vr.reportFh != nil {
		if filename[0] != vr.reportLetter {
			vr.reportLetter = filename[0]
			fmt.Fprintf(vr.reportFh, "\n## %c\n\n", vr.reportLetter)
			fmt.Fprintf(vr.reportFh, "| File                             | Releases | current  | dat64 | history\n")
			fmt.Fprintf(vr.reportFh, "| -------------------------------- | -------- | -------- | ----- | --------\n")
		}
	}

	for i := len(vr.versions) - 1; i >= 0; i-- {
		v := vr.versions[i]
		basename := strings.Replace(vr.pathspec, "%v", v, -1) + filename
		err := vr.validateVersion(&result, v, filename, basename)
		if err != nil {
			log.Fatal(err)
		}
	}

	result.Log(vr, filename)

	if vr.reportFh != nil {
		fmt.Fprintf(vr.reportFh,
			"| %-32s | %-8s | %s | %s | %s\n",
			filename,
			result.VersionRange(vr),
			result.current.Glyph(),
			result.dat64.Glyph(),
			result.history.Glyph(),
		)
	}
}

func (r *validateRunner) validateVersion(out *validateOutput, version string, filename, basename string) error {
	p := r.looseParsers[version]
	q := r.strictParsers[version]

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
		}
		out.history.update(statusOk)
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

func (vr *validateRunner) Close() {
	if vr.reportFh != nil {
		vr.reportFh.Close()
	}
}

func init() {
	flags := pflag.NewFlagSet("analyze", pflag.ExitOnError)
	flagReport := flags.String("report", "", "output status to a Markdown file")
	flagVersions := flags.StringArray("version", []string{}, "only test against specific versions")

	cmd.AddCommand(&cmd.Command{
		Name:        "validate",
		Description: "Perform cross-version validation on data formats",
		Usage:       "pogo validate <source path pattern> <data file names...>",

		MinArgs: 1,
		MaxArgs: -1,

		Flags: flags,

		Action: func(args []string) {
			versions := *flagVersions
			if len(versions) == 0 {
				versions = []string{
					"0.11", "1.0", "1.1", "1.2", "1.3",
					"2.0", "2.1", "2.2", "2.3", "2.4", "2.5", "2.6",
					"3.0", "3.1", "3.2", "3.3", "3.4", "3.5", "3.6", "3.7", "3.8", "3.9",
					"3.10", "3.11", "3.12", "3.13", "3.14", "3.15", "3.16", "3.17", "3.18", "3.19",
				}
			}

			r := validateRunner{
				opencache: poefs.NewOpenCache(),
				pathspec:  args[0],
				fmtDir:    *cmd.GlobalFmt,
				reportTo:  *flagReport,
				verbose:   *cmd.GlobalVerbose > 0,
				versions:  versions,
			}

			formats := args[1:]
			if len(formats) == 0 {
				for _, f := range strings.Split(allFormats, "\n") {
					if f != "" {
						formats = append(formats, f)
					}
				}
			}

			sort.Slice(formats, func(a, b int) bool {
				return strings.ToLower(formats[a]) < strings.ToLower(formats[b])
			})

			r.Init()
			for _, filename := range formats {
				r.Validate(filename)
			}
			r.Close()
		},
	})
}
