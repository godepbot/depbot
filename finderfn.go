package depbot

type FinderFn func(wd string) (Dependencies, error)
