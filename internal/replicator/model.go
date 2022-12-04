package replicator

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SourceSpec struct {
	Name            string `json:"name"`
	ReplicationType string `json:"replicationType"`
}

type TargetSpec struct {
	Name            string `json:"name"`
	ReplicationType string `json:"replicationType"`
}

type ArtifactSpec struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Repo    string `json:"repo"`
}

type Replication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ReplicationSpec `json:"spec,omitempty"`
}

type ReplicationSpec struct {
	Customer        string         `json:"customer"`
	ReplicationType string         `json:"replicationType"`
	VaultURI        string         `json:"vaultURI"`
	TargetRegistry  string         `json:"targetRegistry"`
	SourceRegistry  string         `json:"sourceResgistry"`
	Source          SourceSpec     `json:"source"`
	Target          TargetSpec     `json:"target"`
	Images          []ArtifactSpec `json:"images"`
	Charts          []ArtifactSpec `json:"charts"`
}
