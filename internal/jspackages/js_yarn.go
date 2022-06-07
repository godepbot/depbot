package jspackages

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/godepbot/depbot"
)

func FindYarnDependencies(wd string) (depbot.Dependencies, error) {
	pths := []string{}

	filepath.WalkDir(wd, func(path string, d fs.DirEntry, err error) error {
		if strings.Contains(path, jsYarnLockFile) {
			pths = append(pths, path)
		}

		return nil
	})

	dependencies := depbot.Dependencies{}

	depRegex, err := regexp.Compile(`(?:"?([^\s]+?)@)`)
	if err != nil {
		return dependencies, fmt.Errorf("error compiling regexp: %w", err)
	}

	versionRegex, err := regexp.Compile(`"(\d.+?)"`)
	if err != nil {
		// Continue it cannot parse the version
		return dependencies, fmt.Errorf("error compiling regexp: %w", err)
	}

	for _, p := range pths {
		openFile, err := ioutil.ReadFile(p)
		if err != nil {
			return dependencies, fmt.Errorf("error reading dependency file '%v': %w", p, err)
		}

		dependencies = append(dependencies, depbot.Dependency{
			File: jsYarnLockFile,
			Name: jsDependencyNameJs,
			Kind: depbot.DependencyKindLanguage,
		})

		rawFile := strings.Split(string(openFile), "\n")
		version := ""
		name := ""

		for _, line := range rawFile {
			if strings.Contains(line, "lockfile") {
				dependencies = append(dependencies, depbot.Dependency{
					File:    jsYarnLockFile,
					Name:    jsDependencyNameYARN,
					Version: strings.TrimSpace(strings.ReplaceAll(line, "#", "")),
					Kind:    depbot.DependencyKindTool,
				})
				continue
			}

			result := depRegex.FindStringSubmatch(line)
			// position 0 is the whole match for the regexp
			// position 1 is the clean string
			if len(result) >= 2 {
				name = result[1]
			}

			if strings.Contains(line, "version") {
				version = versionRegex.FindString(line)
			}

			if version != "" && name != "" {
				dependencies = append(dependencies, depbot.Dependency{
					File:    jsYarnLockFile,
					Name:    strings.ReplaceAll(name, "\"", ""),
					Version: strings.ReplaceAll(version, "\"", ""),
					Kind:    depbot.DependencyKindLibrary,
				})
				name = ""
				version = ""
			}
		}
	}

	return dependencies, nil
}
