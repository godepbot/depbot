package depbot

// A FinderFn is a function that can find dependencies for multiple means
// the main one would be looking within the passed working directory for
// specific files.
type FinderFn func(wd string) (DependencyAnalisys, error)
