package depbot

const (
	// constants for errors messages
	MessageError_NoDependencies = "no dependendies found"
	MessageError_MissingApiKey  = "missing api key"
	MessageError_NoSyncDep      = "could not sync the dependencies"

	// constatns for success messages
	MessageSucces_SyncDep = "dependencies synchronized."

	//constants for env variables
	EnvVariable_ApiKey     = "DEPBOT_API_KEY"
	EnvVariable_ServerADDR = "DEPBOT_SERVER_ADDR"
)
