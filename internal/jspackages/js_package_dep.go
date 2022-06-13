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
		if strings.Contains(path, "node_modules") {
			return nil
		}

		if filepath.Base(path) == jsPackageLockFile {
			hasPackageLockDeps = true
			return nil
		}
		if filepath.Base(path) == jsPackageFile {
			pths = append(pths, path)
		}

		return nil
	})

	dependencies := depbot.Dependencies{}

	if hasPackageLockDeps {
		return dependencies, nil
	}

	for _, p := range pths {
		relPath, _ := filepath.Rel(wd, p)
		if relPath == "" {
			relPath = jsPackageFile
		}

		openFile, err := ioutil.ReadFile(p)
		if err != nil {
			return dependencies, fmt.Errorf("error reading dependency file '%v': %w", p, err)
		}

		packageJson := PackageJson{}
		errU := json.Unmarshal(openFile, &packageJson)
		if errU != nil {
			return dependencies, fmt.Errorf("error parsing dependency file '%v': %w", p, errU)
		}
		dependencies = append(dependencies, packageDependencies(packageJson, relPath)...)
	}

	return dependencies, nil
}

func packageDependencies(p PackageJson, file string) depbot.Dependencies {
	dependencies := depbot.Dependencies{
		{
			File:     file,
			Kind:     depbot.DependencyKindTool,
			Language: depbot.DependencyLanguageJs,
			Name:     jsDependencyNameNPM,
		},
		{
			File:     file,
			Kind:     depbot.DependencyKindLanguage,
			Language: depbot.DependencyLanguageJs,
			Name:     jsDependencyNameJs,
			Version:  p.Version,
		},
	}

	for d := range p.Dependencies {
		dependencies = append(dependencies, depbot.Dependency{
			Direct:   true,
			Kind:     depbot.DependencyKindLibrary,
			Language: depbot.DependencyLanguageJs,
			File:     file,
			Name:     d,
			Version:  p.Dependencies[d],
		})
	}

	for d := range p.DevDependencies {
		dependencies = append(dependencies, depbot.Dependency{
			File:     file,
			Kind:     depbot.DependencyKindLibrary,
			Language: depbot.DependencyLanguageJs,
			Name:     d,
			Version:  p.DevDependencies[d],
		})
	}

	return dependencies
}
