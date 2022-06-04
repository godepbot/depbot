package depbot

import "fmt"

const (
	EnvVariable_ApiKey     = "DEPBOT_API_KEY"
	EnvVariable_ServerADDR = "DEPBOT_SERVER_ADDR"
)

var (
	// ErrorMissingApiKey is an error that will be returned if the API key is missing
	ErrorMissingApiKey error = fmt.Errorf("missing api key")

	// ErrNoDendenciesFound is an error that will be returned after the finders have run
	ErrorNoDependenciesFound error = fmt.Errorf("no dependendies found")

	// ErrorNoSyncDep is an error that will be returned if the dependencies could not be synchronized
	ErrorNoSyncDep error = fmt.Errorf("could not sync the dependencies")
)
