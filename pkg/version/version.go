package version

import "runtime"

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
