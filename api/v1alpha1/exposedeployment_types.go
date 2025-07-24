/*
Copyright 2025.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ExposeDeploymentSpec defines the desired state of ExposeDeployment.
type ExposeDeploymentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	Image string `json:"image"`
	// +kubebuilder:validation:Maximum=10
	Replicas            int32          `json:"replicas"`
	MinAvailableTimeSec int32          `json:"minavailabletimesec"`
	CustomEnv           []string       `json:"customenv,omitempty"`
	PortDefinition      PortDefinition `json:"portdefinition,omitempty"`
}

// ExposeDeploymentStatus defines the observed state of ExposeDeployment.
type ExposeDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ReadyPods     int32             `json:"readyPods,omitempty"`
	AvailablePods int32             `json:"availablePods,omitempty"`
	Conditions    []CustomCondition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ExposeDeployment is the Schema for the exposedeployments API.
type ExposeDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ExposeDeploymentSpec   `json:"spec,omitempty"`
	Status            ExposeDeploymentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ExposeDeploymentList contains a list of ExposeDeployment.
type ExposeDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExposeDeployment `json:"items"`
}

// Custom condition types for ExposeDeployment
type MyPodConditionType string

const (
	LastReconcileSucceeded MyPodConditionType = "LastReconcileSucceeded"
	Ready                  MyPodConditionType = "Ready"
	Available              MyPodConditionType = "Available"
)

type CustomCondition struct {
	Type               MyPodConditionType `json:"type,omitempty"`
	Status             bool               `json:"status,omitempty"`
	Message            string             `json:"message,omitempty"`
	Reason             string             `json:"reason,omitempty"`
	LastTransitionTime metav1.Time        `json:"lastTransitionTime,omitempty"`
}

type PortDefinition struct {
	Port       int32 `json:"port,omitempty"`
	TargetPort int32 `json:"targetport,omitempty"`
}

func init() {
	SchemeBuilder.Register(&ExposeDeployment{}, &ExposeDeploymentList{})
}
