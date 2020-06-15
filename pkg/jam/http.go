package jam

type StackCreateConfig struct {
	Repository      StackRepo
	CustomSubDomain string
	IsThirdParty    bool
}

type StackDestroyConfig struct {
	ID           string
	IsThirdParty bool
}

// RoutePrefix is the JAM endpoint
const RoutePrefix = "/jam"
