package mevcommit

import (
	"runtime/debug"
	"strings"
)

var (
	// version will be the version git tag if the binary is built with
	// the Makefile and the tag is set. If the tag does not exist the
	// version will be set to "(devel)". If the binary is built some
	// other way, it will be set to "unknown".
	version = "unknown"

	// revision is set from the vcs.revision tag in Go 1.18+.
	revision = "unknown"

	// dirtyBuild is set from the vcs.modified tag in Go 1.18+.
	dirtyBuild = true
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	for _, kv := range info.Settings {
		if kv.Value == "" {
			continue
		}
		switch kv.Key {
		case "vcs.revision":
			revision = kv.Value
		case "vcs.modified":
			dirtyBuild = kv.Value == "true"
		}
	}

}

// Version returns the build version string of the binary in format:
//
//	<tag>-<commit-hash>[-dirty]
//
// If the tag is not set during build (the tag does not exist or the binary is
// not build with the Makefile), the tag element will remain empty. If the
// binary is built outside the vcs repository, the returned result will be
// "devel". If the binary is built with the Makefile and the repository
// is dirty, the -dirty element will be appended to the result.
//
// Examples:
//
//	"devel"
//	"rev-333ab74"
//	"rev-333ab74-dirty"
//	"v1.0.0-rev-333ab74"
//	"v1.0.0-rev-333ab74-dirty"
func Version() string {
	parts := make([]string, 0, 3)
	if version != "unknown" && version != "" {
		parts = append(parts, version)
	}
	if revision != "unknown" && revision != "" {
		parts = append(parts, "rev")
		commit := revision
		if len(commit) > 7 {
			commit = commit[:7]
		}
		parts = append(parts, commit)
		if dirtyBuild {
			parts = append(parts, "dirty")
		}
	}
	if len(parts) == 0 {
		return "devel"
	}
	return strings.Join(parts, "-")
}
