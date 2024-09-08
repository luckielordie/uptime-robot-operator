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

// +kubebuilder:validation:Enum=sms;email;twitter;webhook;pushbullet;zapier;pro-sms;pushover;slack;voice-call;splunk;pagerduty;opsgenie;ms-teams;google-chat;discord
type AlertContactType string

const (
	SMS        AlertContactType = "sms"
	EMAIL      AlertContactType = "email"
	TWITTER    AlertContactType = "twitter"
	WEBHOOK    AlertContactType = "webhook"
	PUSHBULLET AlertContactType = "pushbullet"
	ZAPIER     AlertContactType = "zapier"
	PROSMS     AlertContactType = "pro-sms"
	PUSHOVER   AlertContactType = "pushover"
	SLACK      AlertContactType = "slack"
	VOICECALL  AlertContactType = "voice-call"
	SPLUNK     AlertContactType = "splunk"
	PAGERDUTY  AlertContactType = "pagerduty"
	OPSGENIE   AlertContactType = "opsgenie"
	TEAMS      AlertContactType = "ms-teams"
	GOOGLECHAT AlertContactType = "google-chat"
	DISCORD    AlertContactType = "discord"
)

// AlertContactSpec defines the desired state of AlertContact
type AlertContactSpec struct {
	// Name is a friendly name for your AlertContact
	Name  string           `json:"name"`
	Type  AlertContactType `json:"type"`
	Value string           `json:"value"`
}

// AlertContactStatus defines the observed state of AlertContact
type AlertContactStatus struct {
	Id     string           `json:"id"`
	Name   string           `json:"name"`
	Type   AlertContactType `json:"type"`
	Value  string           `json:"value"`
	Status int              `json:"status"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AlertContact is the Schema for the alertcontacts API
type AlertContact struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlertContactSpec   `json:"spec,omitempty"`
	Status AlertContactStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AlertContactList contains a list of AlertContact
type AlertContactList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlertContact `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlertContact{}, &AlertContactList{})
}
