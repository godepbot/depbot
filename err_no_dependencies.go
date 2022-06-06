package depbot

import "fmt"

var (
	// ErrNoDendenciesFound is an error that will be returned after the finders have run
	ErrorNoDependenciesFound error = fmt.Errorf("no dependendies found")
)
