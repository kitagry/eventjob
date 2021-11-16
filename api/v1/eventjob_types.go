/*
Copyright 2021.

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

package v1

import (
	"fmt"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EventJobTriggerType string

const (
	CompleteTrigger EventJobTriggerType = "Complete"
	FailedTrigger   EventJobTriggerType = "Failed"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type EventJobTrigger struct {
	metav1.TypeMeta `json:",inline"`

	Name string `json:"name"`

	Type EventJobTriggerType `json:"type"`
}

func (t EventJobTrigger) Match(e corev1.Event) bool {
	return strings.HasSuffix(e.Message, fmt.Sprintf("status: %s", t.Type)) && e.InvolvedObject.Name == t.Name
}

// EventJobSpec defines the desired state of EventJob
type EventJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Trigger EventJobTrigger `json:"trigger"`

	// Job
	JobTemplate batchv1.JobTemplateSpec `json:"jobTemplate"`
}

// EventJobStatus defines the observed state of EventJob
type EventJobStatus struct { // INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	LastScheduledTime metav1.Time `json:"last_scheduled_time"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EventJob is the Schema for the eventjobs API
type EventJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventJobSpec   `json:"spec,omitempty"`
	Status EventJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EventJobList contains a list of EventJob
type EventJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EventJob{}, &EventJobList{})
}
