package version

import (
	"fmt"
	"runtime/debug"
)

var (
	// These will be set by ldflags during build
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// GetVersion returns the application version
func GetVersion() string {
	if Version != "dev" {
		return Version
	}

	// Fallback to build info for go install users
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}

	return "dev"
}

// GetBuildInfo returns detailed build information
func GetBuildInfo() BuildInfo {
	buildInfo := BuildInfo{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
	}

	// If built with go install, get info from debug.BuildInfo
	if info, ok := debug.ReadBuildInfo(); ok {
		if buildInfo.Version == "dev" && info.Main.Version != "(devel)" {
			buildInfo.Version = info.Main.Version
		}

		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				if buildInfo.GitCommit == "unknown" {
					buildInfo.GitCommit = setting.Value
				}
			case "vcs.time":
				if buildInfo.BuildDate == "unknown" {
					buildInfo.BuildDate = setting.Value
				}
			}
		}
	}

	return buildInfo
}

// BuildInfo contains version and build information
type BuildInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildDate string `json:"build_date"`
}

// String returns a formatted version string
func (b BuildInfo) String() string {
	return fmt.Sprintf("gcp-iam %s (commit: %s, built: %s)",
		b.Version, b.GitCommit[:8], b.BuildDate)
}
