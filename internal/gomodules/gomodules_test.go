package gomodules

import (
	"fmt"
	"io/ioutil"
	"testing"
)

// d, err := os.TempDir()
// Write go.mod en directory
// Write package/go.mod
// Write package/p2/go.mod
// Check if no go.mod

var FileContent = `
	module github.com/godepbot/depbot

	go 1.18
	
	require (
		golang.org/x/mod v0.5.1 // indirect
		golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	)
`

func Test_SingleDependency(t *testing.T) {

	tmp := t.TempDir()

	file, errTempFile := ioutil.TempFile(tmp, "go*.mod")
	if errTempFile != nil {
		t.Logf("got an error but should be nil, error: %v ", errTempFile)
		t.Fail()
		return
	}

	fmt.Println("TempFile is: ", file.Name())
	errWriteFile := ioutil.WriteFile(file.Name(), []byte(FileContent), 0644)
	if errWriteFile != nil {
		t.Logf("got an error but should be nil, error : %v ", errWriteFile.Error())
		t.Fail()
		return
	}

	dependecies, err := FindDependencies(tmp)
	if err != nil {
		t.Logf("got an error but should be nil, error : %v ", err.Error())
		t.Fail()
		return
	}

	if len(dependecies) != 3 {
		t.Logf("got %v, but was expected %v", len(dependecies), 3)
		t.Fail()
		return
	}
}
