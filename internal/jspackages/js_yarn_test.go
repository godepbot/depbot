package jspackages_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/godepbot/depbot/internal/jspackages"
)

var yarnFile = `
# THIS IS AN AUTOGENERATED FILE. DO NOT EDIT THIS FILE DIRECTLY.
# yarn lockfile v1

"@babel/cli@^7.0.0":
  version "7.12.10"
  resolved "https://registry.yarnpkg.com/@babel/cli/-/cli-7.12.10.tgz"
  integrity sha512-+y4ZnePpvWs1fc/LhZRTHkTesbXkyBYuOB+5CyodZqrEuETXi3zOVfpAQIdgC3lXbHLTDG9dQosxR9BhvLKDLQ==
  dependencies:
    commander "^4.0.1"
    convert-source-map "^1.1.0"
    fs-readdir-recursive "^1.1.0"
    glob "^7.0.0"
    lodash "^4.17.19"
    make-dir "^2.1.0"
    slash "^2.0.0"
    source-map "^0.5.0"
  optionalDependencies:
    "@nicolo-ribaudo/chokidar-2" "2.1.8-no-fsevents"
    chokidar "^3.4.0"

"@babel/compat-data@^7.12.5", "@babel/compat-data@^7.12.7":
  version "7.12.7"
  resolved "https://registry.yarnpkg.com/@babel/compat-data/-/compat-data-7.12.7.tgz"
  integrity sha512-YaxPMGs/XIWtYqrdEOZOCPsVWfEoriXopnsz3/i7apYPXQ3698UFhS6dVT1KN5qOsWmVgw/FOrmQgpRaZayGsw==

makeerror@1.0.x:
  version "1.0.11"
  resolved "https://registry.yarnpkg.com/makeerror/-/makeerror-1.0.11.tgz"
  integrity sha1-4BpckQnyr3lmDk6LlYd5AYT1qWw=
  dependencies:
    tmpl "1.0.x"
`

var secondYarnFile = `
# THIS IS AN AUTOGENERATED FILE. DO NOT EDIT THIS FILE DIRECTLY.
# yarn lockfile v1

"@milo@^7.0.0":
  version "7.12.10"
  resolved "https://registry.yarnpkg.com/@babel/cli/-/cli-7.12.10.tgz"
  integrity sha512-+y4ZnePpvWs1fc/LhZRTHkTesbXkyBYuOB+5CyodZqrEuETXi3zOVfpAQIdgC3lXbHLTDG9dQosxR9BhvLKDLQ==

  "@babel/compat-data@^7.12.5", "@babel/compat-data@^7.12.7":
  version "7.12.7"
  resolved "https://registry.yarnpkg.com/@babel/compat-data/-/compat-data-7.12.7.tgz"
  integrity sha512-YaxPMGs/XIWtYqrdEOZOCPsVWfEoriXopnsz3/i7apYPXQ3698UFhS6dVT1KN5qOsWmVgw/FOrmQgpRaZayGsw==

makeerror@1.0.x:
  version "1.0.11"
  resolved "https://registry.yarnpkg.com/makeerror/-/makeerror-1.0.11.tgz"
  integrity sha1-4BpckQnyr3lmDk6LlYd5AYT1qWw=
  dependencies:
    tmpl "1.0.x"
`

func Test_Yarn_SingleDependency(t *testing.T) {

	file, err := os.Create(t.TempDir() + "/" + "yarn.lock")

	if err != nil {
		t.Fatalf("got an error but should be nil, error: %v ", err)
		return
	}

	errWriteFile := ioutil.WriteFile(file.Name(), []byte(yarnFile), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	dependencies, err := jspackages.FindYarnDependencies(file.Name())
	if err != nil {
		t.Fatalf("got an error but should be nil, error : %v ", err.Error())
		return
	}

	if len(dependencies) != 5 {
		t.Fatalf("got %v, but was expected %v", len(dependencies), 5)
		return
	}

	tcases := []struct {
		name    string
		version string
		exist   bool
	}{
		{
			name:    "@babel/cli",
			version: "7.12.10",
			exist:   true,
		},
		{
			name:    "makeerror",
			version: "1.0.11",
			exist:   true,
		},
		{
			name:    "cli",
			version: "^7.16.0",
			exist:   false,
		},
	}

	for _, tcase := range tcases {
		var exist bool
		for _, d := range dependencies {
			if d.Name == tcase.name && d.Version == tcase.version {
				exist = true
			}
		}

		if exist != tcase.exist {
			complement := "exist"
			if !tcase.exist {
				complement = "no exist"
			}
			t.Fatalf("expected %v with version %v to %v", tcase.name, tcase.version, complement)
		}
	}
}

func Test_Yarn_MultipleDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	errWriteFile := ioutil.WriteFile(tmpDir+"/yarn.lock", []byte(yarnFile), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	newDirectoriesPath := tmpDir + "/package/v2"
	packagePath := tmpDir + "/package"

	err := os.MkdirAll(newDirectoriesPath, os.ModePerm)

	if err != nil {
		t.Fatalf("got an error but should be nil, error : %v ", err)
		return
	}

	errWriteFile = ioutil.WriteFile(packagePath+"/yarn.lock", []byte(yarnFile), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	errWriteFile = ioutil.WriteFile(newDirectoriesPath+"/yarn.lock", []byte(secondYarnFile), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	dependencies, errFindDep := jspackages.FindYarnDependencies(tmpDir)
	if errFindDep != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	fmt.Println("----> ", dependencies)

	tcases := []struct {
		name    string
		version string
		exist   bool
	}{
		{
			name:    "@babel/cli",
			version: "7.12.10",
			exist:   true,
		},
		{
			name:    "makeerror",
			version: "1.0.11",
			exist:   true,
		},
		{
			name:    "cli",
			version: "^7.16.0",
			exist:   false,
		},
		{
			name:    "depbot",
			version: "^7.16.0",
			exist:   false,
		},
		{
			name:    "@milo",
			version: "7.12.10",
			exist:   true,
		},
	}

	for _, tcase := range tcases {
		var exist bool
		for _, d := range dependencies {
			if d.Name == tcase.name && d.Version == tcase.version {
				exist = true
			}
		}

		if exist != tcase.exist {
			complement := "exist"
			if !tcase.exist {
				complement = "no exist"
			}
			t.Fatalf("expected %v with version %v to %v", tcase.name, tcase.version, complement)
		}
	}

}

func Test_Yarn_NoDependency(t *testing.T) {
	tmp := t.TempDir()

	dependencies, err := jspackages.FindYarnDependencies(tmp)

	if err != nil {
		t.Fatalf("error finding the dependencies : %v ", err.Error())
		return
	}

	if len(dependencies) > 0 {
		t.Fatalf("got %v, but was expected %v", len(dependencies), 0)
	}

}