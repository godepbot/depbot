package jspackages_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/godepbot/depbot"
	"github.com/godepbot/depbot/internal/jspackages"
)

var packageLockFile = `
{
	"name": "buffalo",
	"version": "1.0.0",
	"lockfileVersion": 2,
	"requires": true,
	"packages": {
		"": {
			"name": "buffalo",
			"version": "1.0.0",
			"license": "MIT",
			"dependencies": {
				"@fortawesome/fontawesome-free": "^5.15.4"
			},
			"devDependencies": {
				"@babel/cli": "^7.16.0"
			}
		}
	},
	"dependencies": {
		"@babel/cli": {
			"version": "7.16.0",
			"resolved": "https://registry.yarnpkg.com/@babel/cli/-/cli-7.16.0.tgz",
			"integrity": "sha512-WLrM42vKX/4atIoQB+eb0ovUof53UUvecb4qGjU2PDDWRiZr50ZpiV8NpcLo7iSxeGYrRG0Mqembsa+UrTAV6Q==",
			"dev": true
		},
		"@babel/code-frame": {
			"version": "7.16.0",
			"resolved": "https://registry.yarnpkg.com/@babel/code-frame/-/code-frame-7.16.0.tgz",
			"integrity": "sha512-IF4EOMEV+bfYwOmNxGzSnjR2EmQod7f1UXOpZM3l4i4o4QNwzjtJAu/HxdjHq0aYBvdqMuQEY1eg0nqW9ZPORA==",
			"dev": true
		},
		"@babel/compat-data": {
			"version": "7.16.4",
			"resolved": "https://registry.yarnpkg.com/@babel/compat-data/-/compat-data-7.16.4.tgz",
			"integrity": "sha512-1o/jo7D+kC9ZjHX5v+EHrdjl3PhxMrLSOTGsOdHJ+KL8HCaEK6ehrVL2RS6oHDZp+L7xLirLrPmQtEng769J/Q==",
			"dev": true
		}
	}
}
`

var secondPackagelockFile = `
{
	"name": "buffalo",
	"version": "1.0.0",
	"lockfileVersion": 2,
	"requires": true,
	"packages": {
		"": {
			"name": "buffalo",
			"version": "1.0.0",
			"license": "MIT",
			"dependencies": {
				"@fortawesome/fontawesome-free": "^5.15.4"
			},
			"devDependencies": {
				"@babel/cli": "^7.16.0"
			}
		}
	},
	"dependencies": {
		"milo": {
		  "version": "7.16.0",
		  "resolved": "https://registry.yarnpkg.com/@babel/cli/-/cli-7.16.0.tgz",
		  "integrity": "sha512-WLrM42vKX/4atIoQB+eb0ovUof53UUvecb4qGjU2PDDWRiZr50ZpiV8NpcLo7iSxeGYrRG0Mqembsa+UrTAV6Q==",
		  "dev": true,
		  "license": "MIT"
		},
		"depbot": {
		  "version": "4.1.1",
		  "resolved": "https://registry.yarnpkg.com/commander/-/commander-4.1.1.tgz",
		  "integrity": "sha512-NOKm8xhkzAjzFx8B2v5OAHT+u5pRQc2UCa2Vq9jYL/31o2wi9mxBA7LIFs3sV5VSC49z6pEhfbMULvShKj26WA==",
		  "dev": true,
		  "license": "MIT"
		}
	}
}
`

func Test_PackageLock_SingleDependency(t *testing.T) {

	file, err := os.Create(t.TempDir() + "/" + "package-lock.json")

	if err != nil {
		t.Fatalf("got an error but should be nil, error: %v ", err)
		return
	}

	errWriteFile := ioutil.WriteFile(file.Name(), []byte(packageLockFile), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	dependencies, err := jspackages.FindPackageLockDependencies(file.Name())
	if err != nil {
		t.Fatalf("got an error but should be nil, error : %v ", err.Error())
		return
	}

	if len(dependencies) != 5 {
		t.Fatalf("got %v, but was expected %v", len(dependencies), 5)
		return
	}

	for _, dependencie := range dependencies {
		if dependencie.Language != depbot.DependencyLanguageJs {
			t.Fatalf("got %v, but was expected %v", dependencie.Language, depbot.DependencyLanguageJs)
		}
	}

	tcases := []struct {
		name    string
		version string
		exist   bool
	}{
		{
			name:    "@babel/cli",
			version: "7.16.0",
			exist:   true,
		},
		{
			name:    "@babel/code-frame",
			version: "7.16.0",
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

func Test_Pakcagelock_MultipleDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	errWriteFile := ioutil.WriteFile(tmpDir+"/package-lock.json", []byte(packageLockFile), 0644)
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

	errWriteFile = ioutil.WriteFile(packagePath+"/package-lock.json", []byte(packageLockFile), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	errWriteFile = ioutil.WriteFile(newDirectoriesPath+"/package-lock.json", []byte(secondPackagelockFile), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	dependencies, errFindDep := jspackages.FindPackageLockDependencies(tmpDir)
	if errFindDep != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	tcases := []struct {
		name    string
		version string
		exist   bool
	}{
		{
			name:    "@babel/cli",
			version: "7.16.0",
			exist:   true,
		},
		{
			name:    "@babel/code-frame",
			version: "7.16.0",
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
			name:    "depbot",
			version: "4.1.1",
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

func Test_PackageLock_NoDependency(t *testing.T) {
	tmp := t.TempDir()

	dependencies, err := jspackages.FindPackageLockDependencies(tmp)

	if err != nil {
		t.Fatalf("error finding the dependencies : %v ", err.Error())
		return
	}

	if len(dependencies) > 0 {
		t.Fatalf("got %v, but was expected %v", len(dependencies), 0)
	}

}
