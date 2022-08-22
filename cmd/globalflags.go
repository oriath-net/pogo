package cmd

import (
	"github.com/spf13/pflag"
)

var (
	GlobalFlags = pflag.NewFlagSet("pogo", pflag.ExitOnError)

	GlobalFmt     = GlobalFlags.String("fmt", "", "Path to a directory containing custom formats")
	GlobalVersion = GlobalFlags.String("version", "9.99", "Path of Exile version to assume for formats")
	GlobalVerbose = GlobalFlags.CountP("verbose", "v", "display more details (repeatable)")
)
