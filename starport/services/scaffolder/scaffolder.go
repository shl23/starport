// Package scaffolder initializes Starport apps and modifies existing ones
// to add more features in a later time.
package scaffolder

import (
	"strings"

	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

// Scaffolder is Starport app scaffolder.
type Scaffolder struct {
	// path is app's path on the filesystem.
	path string

	// options to configure scaffolding.
	options *scaffoldingOptions
}

// New initializes a new Scaffolder for app at path.
func New(path string, options ...Option) *Scaffolder {
	return &Scaffolder{
		path:    path,
		options: newOptions(options...),
	}
}

func (s *Scaffolder) version() (cosmosver.MajorVersion, error) {
	v, err := cosmosver.Detect(s.path)
	if err != nil {
		return "", err
	}
	return v.Major(), nil
}

func owner(modulePath string) string {
	return strings.Split(modulePath, "/")[1]
}
