package flags

import "github.com/spf13/pflag"

type Flags struct {
	Verbose bool
	Root    string
}

func NewFlags() *Flags {
	verbose := pflag.BoolP("verbose", "v", false, "verbose output")
	root := pflag.StringP("root", "r", "/", "path to root file system")
	pflag.Parse()

	return &Flags{
		Verbose: *verbose,
		Root:    *root,
	}
}
