package replicator

type Registry struct {
	Url                          string `json:"url"`
	Unauthenticated              bool   `json:"unauthenticated"`
	Username                     string `json:"username"`
	Password                     string `json:"password"`
	UseCredentialsFromAzureVault bool   `json:"useCredentialsFromAzureVault"`
	VaultURI                     string `json:"vaultURI"`
	UseExistingSecret            bool   `json:"useExistingSecret"`
	UsernameKey                  bool   `json:"usernameKey"`
	PasswordKey                  bool   `json:"passwordKey"`
}

type Artifact struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Repo    string `json:"repo"`
}

type Replication struct {
	Identifier string     `json:"identifier"`
	Source     Registry   `json:"source"`
	Target     Registry   `json:"target"`
	Images     []Artifact `json:"images"`
	Charts     []Artifact `json:"charts"`
}
