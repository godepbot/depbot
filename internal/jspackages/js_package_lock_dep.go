package jspackages

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/godepbot/depbot"
)

type PackageLockJson struct {
	License      string                 `json:"license" db:"-"`
	Name         string                 `json:"name" db:"-"`
	Version      string                 `json:"version" db:"-"`
	Dependencies map[string]interface{} `json:"dependencies" db:"-"`
}

func FindPackageLockDependencies(wd string) (depbot.Dependencies, error) {
	pths := []string{}

	filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
		if filepath.Base(path) == jsPackageLockFile {
			pths = append(pths, path)
		}

		return nil
	})

	dependencies := depbot.Dependencies{}
	for _, p := range pths {
		relPath, _ := filepath.Rel(wd, p)
		if relPath == "" {
			relPath = jsPackageFile
		}

		openFile, err := ioutil.ReadFile(p)
		if err != nil {
			return dependencies, fmt.Errorf("error reading dependency file '%v': %w", p, err)
		}

		packageJson := PackageLockJson{}
		errU := json.Unmarshal(openFile, &packageJson)
		if errU != nil {
			return dependencies, fmt.Errorf("error parsing dependency file '%v': %w", p, errU)
		}
		dependencies = append(dependencies, packageLockDependencies(packageJson, relPath)...)
	}

	return dependencies, nil
}

func packageLockDependencies(p PackageLockJson, filePath string) depbot.Dependencies {
	dependencies := depbot.Dependencies{
		{
			File:     filePath,
			Kind:     depbot.DependencyKindTool,
			Language: depbot.DependencyLanguageJs,
			Name:     jsDependencyNameNPM,
		},
		{
			File:     filePath,
			Kind:     depbot.DependencyKindLanguage,
			Language: depbot.DependencyLanguageJs,
			Version:  p.Version,
			Name:     jsDependencyNameJs,
		},
	}

	for d := range p.Dependencies {
		version := p.Dependencies[d].(map[string]interface{})["version"]
		transitive := p.Dependencies[d].(map[string]interface{})["dev"]

		dependencies = append(dependencies, depbot.Dependency{
			Direct:   transitive != true,
			File:     filePath,
			Kind:     depbot.DependencyKindLibrary,
			Language: depbot.DependencyLanguageJs,
			Name:     d,
			Version:  version.(string),
		})
	}
	return dependencies
}
