//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoIngress) DeepCopyInto(out *AutoIngress) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoIngress.
func (in *AutoIngress) DeepCopy() *AutoIngress {
	if in == nil {
		return nil
	}
	out := new(AutoIngress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AutoIngress) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoIngressList) DeepCopyInto(out *AutoIngressList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AutoIngress, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoIngressList.
func (in *AutoIngressList) DeepCopy() *AutoIngressList {
	if in == nil {
		return nil
	}
	out := new(AutoIngressList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AutoIngressList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoIngressSpec) DeepCopyInto(out *AutoIngressSpec) {
	*out = *in
	if in.TlsSecretName != nil {
		in, out := &in.TlsSecretName, &out.TlsSecretName
		*out = new(string)
		**out = **in
	}
	if in.ServicePrefixes != nil {
		in, out := &in.ServicePrefixes, &out.ServicePrefixes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.NameSpaces != nil {
		in, out := &in.NameSpaces, &out.NameSpaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.NameSpaceBlackList != nil {
		in, out := &in.NameSpaceBlackList, &out.NameSpaceBlackList
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoIngressSpec.
func (in *AutoIngressSpec) DeepCopy() *AutoIngressSpec {
	if in == nil {
		return nil
	}
	out := new(AutoIngressSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoIngressStatus) DeepCopyInto(out *AutoIngressStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoIngressStatus.
func (in *AutoIngressStatus) DeepCopy() *AutoIngressStatus {
	if in == nil {
		return nil
	}
	out := new(AutoIngressStatus)
	in.DeepCopyInto(out)
	return out
}
