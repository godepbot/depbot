package gomodules_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/godepbot/depbot/internal/gomodules"
)

var fileContent = `
	module github.com/godepbot/depbot

	go 1.18
	
	require (
		github.com/gobuffalo/plush/v4 v4.1.9
		golang.org/x/mod v0.5.1 // indirect
		golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	)
`

func Test_SingleDependency(t *testing.T) {

	file, err := os.Create(t.TempDir() + "/" + "go.mod")

	if err != nil {
		t.Logf("got an error but should be nil, error: %v ", err)
		t.Fail()
		return
	}

	errWriteFile := ioutil.WriteFile(file.Name(), []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Logf("got an error but should be nil, error : %v ", errWriteFile.Error())
		t.Fail()
		return
	}

	dependencies, err := gomodules.FindDependencies(file.Name())
	if err != nil {
		t.Logf("got an error but should be nil, error : %v ", err.Error())
		t.Fail()
		return
	}

	if len(dependencies) != 4 {
		t.Logf("got %v, but was expected %v", len(dependencies), 4)
		t.Fail()
		return
	}

	if dependencies[0].Name != "go.mod" {
		t.Logf("Got %v, but was expected %v", dependencies[0].Name, "go.mod")
		t.Fail()
	}

	if dependencies[1].Name != "github.com/gobuffalo/plush/v4" {
		t.Logf("Got %v, but was expected %v", dependencies[1].Name, "github.com/gobuffalo/plush/v4")
		t.Fail()
	}

	if dependencies[3].Version != "v0.0.0-20191011141410-1b5146add898" {
		t.Logf("Got %v, but was expected %v", dependencies[1].Version, "v0.0.0-20191011141410-1b5146add898")
		t.Fail()
	}

}

func Test_MultipleDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	errWriteFile := ioutil.WriteFile(tmpDir+"/go.mod", []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Logf("got an error but should be nil, error : %v ", errWriteFile.Error())
		t.Fail()
		return
	}

	newDirectoriesPath := tmpDir + "/package/v2"
	packagePath := tmpDir + "/package"

	err := os.MkdirAll(newDirectoriesPath, os.ModePerm)

	if err != nil {
		t.Fatalf("got an error but should be nil, error : %v ", err)
		return
	}

	errWriteFile = ioutil.WriteFile(packagePath+"/go.mod", []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	errWriteFile = ioutil.WriteFile(newDirectoriesPath+"/go.mod", []byte(fileContent), 0644)
	if errWriteFile != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	dependencies, errFindDep := gomodules.FindDependencies(tmpDir)
	if errFindDep != nil {
		t.Fatalf("got an error but should be nil, error : %v ", errWriteFile.Error())
		return
	}

	if len(dependencies) != 12 {
		t.Fatalf("got %v, but was expected %v", len(dependencies), 12)
		return
	}

}

func Test_NoDependency(t *testing.T) {
	tmp := t.TempDir()

	dependencies, err := gomodules.FindDependencies(tmp)

	if err != nil {
		t.Logf("Error finding the dependencies : %v ", err.Error())
		t.Fail()
		return
	}

	if len(dependencies) > 0 {
		t.Logf("got %v, but was expected %v", len(dependencies), 0)
		t.Fail()
	}

}
