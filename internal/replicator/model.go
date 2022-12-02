package replicator

type Source struct {
	Name            string `json:"name"`
	ReplicationType string `json:"replicationType"`
}

type Target struct {
	Name            string `json:"name"`
	ReplicationType string `json:"replicationType"`
}

type Artifact struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Repo    string `json:"repo"`
}

type Replication struct {
	Customer       string     `json:"customer"`
	VaultURI       string     `json:"vaultURI"`
	TargetRegistry string     `json:"targetRegistry"`
	SourceRegistry string     `json:"sourceResgistry"`
	Source         Source     `json:"source"`
	Target         Target     `json:"target"`
	Images         []Artifact `json:"images"`
	Charts         []Artifact `json:"charts"`
}
