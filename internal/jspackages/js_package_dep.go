package jspackages

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/godepbot/depbot"
)

type PackageJson struct {
	License         string            `json:"license" db:"-"`
	Main            string            `json:"main" db:"-"`
	Name            string            `json:"name" db:"-"`
	Repository      string            `json:"repository" db:"-"`
	Scripts         map[string]string `json:"scripts" db:"-"`
	Version         string            `json:"version" db:"-"`
	Dependencies    map[string]string `json:"dependencies" db:"-"`
	DevDependencies map[string]string `json:"devDependencies" db:"-"`
}

func FindPackageDependencies(wd string) (depbot.Dependencies, error) {
	pths := []string{}
	var hasPackageLockDeps bool

	filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
		if strings.Contains(path, jsPackageLockFile) {
			hasPackageLockDeps = true
			return nil
		}
		if strings.Contains(path, jsPackageFile) {
			pths = append(pths, path)
		}

		return nil
	})

	dependencies := depbot.Dependencies{}

	if hasPackageLockDeps {
		return dependencies, nil
	}

	for _, p := range pths {
		openFile, err := ioutil.ReadFile(p)
		if err != nil {
			return dependencies, fmt.Errorf("error reading dependency file '%v': %w", p, err)
		}

		packageJson := PackageJson{}
		errU := json.Unmarshal(openFile, &packageJson)
		if errU != nil {
			return dependencies, fmt.Errorf("error parsing dependency file '%v': %w", p, errU)
		}
		dependencies = append(dependencies, packageDependencies(packageJson)...)
	}

	return dependencies, nil
}

func packageDependencies(p PackageJson) depbot.Dependencies {
	dependencies := depbot.Dependencies{
		{
			File: jsPackageFile,
			Name: jsDependencyNameNPM,
			Kind: depbot.DependencyKindTool,
		},
		{
			File:    jsPackageFile,
			Version: p.Version,
			Name:    jsDependencyNameJs,
			Kind:    depbot.DependencyKindLanguage,
		},
	}

	for d := range p.Dependencies {
		dependencies = append(dependencies, depbot.Dependency{
			File:    jsPackageFile,
			Name:    d,
			Version: p.Dependencies[d],
			Kind:    depbot.DependencyKindLibrary,
			Direct:  true,
		})
	}

	for d := range p.DevDependencies {
		dependencies = append(dependencies, depbot.Dependency{
			File:    jsPackageFile,
			Name:    d,
			Version: p.DevDependencies[d],
			Kind:    depbot.DependencyKindLibrary,
		})
	}

	return dependencies
}
