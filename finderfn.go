package depbot

type FinderFn func(wd string) ([]Dependency, error)
