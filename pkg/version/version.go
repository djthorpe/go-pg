package version

import (
	"os"
	"path/filepath"
	"runtime"
)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	GitSource   string
	GitTag      string
	GitBranch   string
	GitHash     string
	GoBuildTime string
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func ExecName() string {
	name, err := os.Executable()
	if err != nil {
		return "unknown"
	}
	return filepath.Base(name)
}

func Version() string {
	if GitTag != "" {
		return GitTag
	}
	if GitBranch != "" {
		return GitBranch
	}
	if GitHash != "" {
		return GitHash
	}
	return "dev"
}

func Compiler() string {
	return runtime.Compiler + "/" + runtime.Version()
}
