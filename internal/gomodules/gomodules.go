// package gomodules does the analysis of Go modules and returns its findings.
package gomodules

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/godepbot/depbot"
	"golang.org/x/mod/modfile"
)

const (
	goDependencyFile string = "go.mod"
	dependencyNameGo string = "Go"
)

// Find walks the directory three and looks for go.mod files
// to then parse dependencies and return them.
func FindDependencies(wd string) (depbot.Dependencies, error) {
	pths := []string{}

	filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
		if filepath.Base(path) == goDependencyFile {
			pths = append(pths, path)
		}

		return nil
	})

	dependencies := depbot.Dependencies{}

	for _, p := range pths {
		d, err := ioutil.ReadFile(p)
		if err != nil {
			return dependencies, fmt.Errorf("error reading dependency file '%v': %w", p, err)
		}

		f, err := modfile.Parse(p, d, nil)
		if err != nil {
			return dependencies, fmt.Errorf("error parsing dependencies on file '%v': %w", p, err)
		}

		dependencies = append(dependencies, depbot.Dependency{
			File:    goDependencyFile,
			Version: f.Go.Version,
			Name:    dependencyNameGo,
			Kind:    depbot.DependencyKindLanguage,
			Direct:  true,
		})

		for _, r := range f.Require {
			dependencies = append(dependencies, depbot.Dependency{
				File:    goDependencyFile,
				Name:    r.Mod.Path,
				Version: r.Mod.Version,
				Kind:    depbot.DependencyKindLibrary,
				Direct:  !r.Indirect,
			})
		}
	}

	return dependencies, nil
}
