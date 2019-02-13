package version

import (
	"fmt"
)

const (
	template = `
	version:    %s
	revision:   %s
	branch:     %s
	goVersion:  %s
	buildTool:  %s
	buildTime:  %s` + "\n"
)

var (
	// Version is the semantic version
	Version string

	// Revision is the SHA-1 of the git revision
	Revision string

	// Branch is the name of the git branch
	Branch string

	// GoVersion is the go compiler version
	GoVersion string

	// BuildTool contains the name and version of build tool
	BuildTool string

	// BuildTime is the time binary built
	BuildTime string
)

// String returns a string describing the version information in details
func String() string {
	return fmt.Sprintf(template, Version, Revision, Branch, GoVersion, BuildTool, BuildTime)
}
