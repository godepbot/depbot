package gomodules

import (
	"fmt"
	"os"
	"testing"
)

// d, err := os.TempDir()
// Write go.mod en directory
// Write package/go.mod
// Write package/p2/go.mod
// Check if no go.mod

func Test_SingleDependency(t *testing.T) {
	// r := require.New(t)
	fmt.Println("Hello moto")
	tmp := os.TempDir()
	fmt.Println("Tmp is:", tmp)

	// gomodules.FindDependencies()

}
