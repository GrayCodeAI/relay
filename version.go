// Package mux provides LLM provider clients and configuration.
//
// The Version variable is sourced from the VERSION file at the repo root
// and propagated to sub-packages at init time.
package mux

import (
	_ "embed"
	"strings"

	"github.com/GrayCodeAI/mux/provider"
)

//go:embed VERSION
var versionFile string

// Version of the mux library. Single source of truth: VERSION file.
var Version = strings.TrimSpace(versionFile)

func init() {
	provider.SetVersion(Version)
}
