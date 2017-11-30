package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the binary version
	Version string

	// Revision is the short form of git commit
	Revision string

	// Branch is the name of git branch
	Branch string

	// BuildTime is the build time of binary
	BuildTime string

	// GoVersion is the golang version
	GoVersion = runtime.Version()

	// GoOS is the binary operating system
	GoOS = runtime.GOOS

	// GoArch is the binary architecture
	GoArch = runtime.GOARCH
)

// GetSemVer returns the semantic versioning format
func GetSemVer() string {
	return Version + "+" + Revision
}

// GetFullSpec returns the full specifiction
func GetFullSpec() string {
	return fmt.Sprintf("%s+%s %s %s %s %s/%s", Version, Revision, Branch, BuildTime, GoVersion, GoOS, GoArch)
}
