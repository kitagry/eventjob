//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventJob) DeepCopyInto(out *EventJob) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventJob.
func (in *EventJob) DeepCopy() *EventJob {
	if in == nil {
		return nil
	}
	out := new(EventJob)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EventJob) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventJobList) DeepCopyInto(out *EventJobList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EventJob, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventJobList.
func (in *EventJobList) DeepCopy() *EventJobList {
	if in == nil {
		return nil
	}
	out := new(EventJobList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EventJobList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventJobSpec) DeepCopyInto(out *EventJobSpec) {
	*out = *in
	out.Trigger = in.Trigger
	in.JobTemplate.DeepCopyInto(&out.JobTemplate)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventJobSpec.
func (in *EventJobSpec) DeepCopy() *EventJobSpec {
	if in == nil {
		return nil
	}
	out := new(EventJobSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventJobStatus) DeepCopyInto(out *EventJobStatus) {
	*out = *in
	in.LastScheduledTime.DeepCopyInto(&out.LastScheduledTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventJobStatus.
func (in *EventJobStatus) DeepCopy() *EventJobStatus {
	if in == nil {
		return nil
	}
	out := new(EventJobStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventJobTrigger) DeepCopyInto(out *EventJobTrigger) {
	*out = *in
	out.TypeMeta = in.TypeMeta
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventJobTrigger.
func (in *EventJobTrigger) DeepCopy() *EventJobTrigger {
	if in == nil {
		return nil
	}
	out := new(EventJobTrigger)
	in.DeepCopyInto(out)
	return out
}
