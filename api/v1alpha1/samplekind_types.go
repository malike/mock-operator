/*
Copyright 2023.

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
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SampleKindSpec defines the desired state of SampleKind
type SampleKindSpec struct {
	//+kubebuilder:validation:Type:=object
	// Image defines image configuration
	Image ImageSpec `json:"image,omitempty"`
	//+kubebuilder:validation:Type:=number
	//+kubebuilder:default:=2
	// Nodes defines number of instance
	Nodes int32 `json:"nodes,omitempty"`
	//+kubebuilder:validation:Type:=number
	//+kubebuilder:default:=80
	// ContainerPort defines port for container
	ContainerPort int32 `json:"containerPort,omitempty"`
	//+kubebuilder:validation:Type:=number
	//+kubebuilder:default:=80
	// ServicePort defines port for service
	ServicePort int32 `json:"servicePort,omitempty"`
}

// ImageSpec defines Image details
type ImageSpec struct {
	//+kubebuilder:validation:Type:=string
	//+kubebuilder:default:=ghcr.io/malike/sample-mock-service
	// Defines the container image repo for the service
	Repository string `json:"repository,omitempty"`
	//+kubebuilder:validation:Type:=string
	//+kubebuilder:default:=latest
	// Specifies the tag of the container image to be used.
	Tag string `json:"tag,omitempty"`
	//+kubebuilder:validation:Type:=string
	//+kubebuilder:default:=IfNotPresent
	// Specifies ImagePullPolicy of the container image.
	ImagePullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty"`
	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// +optional
	//+kubebuilder:validation:Type:=array
	ImagePullSecrets []corev1.LocalObjectReference `json:"pullSecretName,omitempty"`
}

// SampleKindStatus defines the observed state of SampleKind
type SampleKindStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=samplekind,shortName=smk

// SampleKind is the Schema for the samplekind API
type SampleKind struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SampleKindSpec   `json:"spec,omitempty"`
	Status SampleKindStatus `json:"status,omitempty"`
}

func (spec SampleKindSpec) String() string {
	specString, _ := json.Marshal(spec)
	return string(specString)
}

//+kubebuilder:object:root=true

// SampleKindList contains a list of SampleKind
type SampleKindList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SampleKind `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SampleKind{}, &SampleKindList{})
}
