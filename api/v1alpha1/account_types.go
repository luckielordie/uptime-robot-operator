/*
Copyright 2024.

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

// AccountSpec defines the desired state of Account
type AccountSpec struct {
}

// AccountStatus defines the observed state of Account
type AccountStatus struct {
	Email           string `json:"email,omitempty"`
	MonitorLimit    int    `json:"monitorLimit,omitempty"`
	MonitorInterval int    `json:"monitorInterval,omitempty"`
	UpMonitors      int    `json:"upMonitors,omitempty"`
	DownMonitors    int    `json:"downMonitors,omitempty"`
	PausedMonitors  int    `json:"pausedMonitors,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Account is the Schema for the accounts API
type Account struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AccountSpec   `json:"spec,omitempty"`
	Status AccountStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AccountList contains a list of Account
type AccountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Account `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Account{}, &AccountList{})
}
