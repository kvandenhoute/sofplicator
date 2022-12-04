/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

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

// ReplicationSpec defines the desired state of Replication
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

// ReplicationStatus defines the observed state of Replication
type ReplicationStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Replication is the Schema for the replications API
type Replication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReplicationSpec   `json:"spec,omitempty"`
	Status ReplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ReplicationList contains a list of Replication
type ReplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Replication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Replication{}, &ReplicationList{})
}
