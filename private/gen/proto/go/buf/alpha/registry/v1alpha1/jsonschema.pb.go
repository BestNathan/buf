// Copyright 2020-2024 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: buf/alpha/registry/v1alpha1/jsonschema.proto

package registryv1alpha1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetJSONSchemaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Owner      string `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	Repository string `protobuf:"bytes,2,opt,name=repository,proto3" json:"repository,omitempty"`
	// Optional reference (if unspecified, will use the repository's default_branch).
	Reference string `protobuf:"bytes,3,opt,name=reference,proto3" json:"reference,omitempty"`
	// A fully qualified name of the type to generate a JSONSchema for, e.g.
	// "pkg.foo.Bar". The type needs to resolve in the referenced module or any of
	// its dependencies. Currently only messages types are supported.
	TypeName string `protobuf:"bytes,4,opt,name=type_name,json=typeName,proto3" json:"type_name,omitempty"`
}

func (x *GetJSONSchemaRequest) Reset() {
	*x = GetJSONSchemaRequest{}
	mi := &file_buf_alpha_registry_v1alpha1_jsonschema_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetJSONSchemaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetJSONSchemaRequest) ProtoMessage() {}

func (x *GetJSONSchemaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_jsonschema_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetJSONSchemaRequest.ProtoReflect.Descriptor instead.
func (*GetJSONSchemaRequest) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_jsonschema_proto_rawDescGZIP(), []int{0}
}

func (x *GetJSONSchemaRequest) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *GetJSONSchemaRequest) GetRepository() string {
	if x != nil {
		return x.Repository
	}
	return ""
}

func (x *GetJSONSchemaRequest) GetReference() string {
	if x != nil {
		return x.Reference
	}
	return ""
}

func (x *GetJSONSchemaRequest) GetTypeName() string {
	if x != nil {
		return x.TypeName
	}
	return ""
}

type GetJSONSchemaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A json schema representing what the json encoded payload for type_name
	// should conform to. This schema is an approximation to be used by editors
	// for validation and autocompletion, not a lossless representation of the
	// type's descriptor.
	JsonSchema []byte `protobuf:"bytes,1,opt,name=json_schema,json=jsonSchema,proto3" json:"json_schema,omitempty"`
}

func (x *GetJSONSchemaResponse) Reset() {
	*x = GetJSONSchemaResponse{}
	mi := &file_buf_alpha_registry_v1alpha1_jsonschema_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetJSONSchemaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetJSONSchemaResponse) ProtoMessage() {}

func (x *GetJSONSchemaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_jsonschema_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetJSONSchemaResponse.ProtoReflect.Descriptor instead.
func (*GetJSONSchemaResponse) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_jsonschema_proto_rawDescGZIP(), []int{1}
}

func (x *GetJSONSchemaResponse) GetJsonSchema() []byte {
	if x != nil {
		return x.JsonSchema
	}
	return nil
}

var File_buf_alpha_registry_v1alpha1_jsonschema_proto protoreflect.FileDescriptor

var file_buf_alpha_registry_v1alpha1_jsonschema_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2f, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x6a, 0x73,
	0x6f, 0x6e, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b,
	0x62, 0x75, 0x66, 0x2e, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x79, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x22, 0x87, 0x01, 0x0a, 0x14,
	0x47, 0x65, 0x74, 0x4a, 0x53, 0x4f, 0x4e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72,
	0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x79, 0x70, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x79, 0x70,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x38, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x4a, 0x53, 0x4f, 0x4e,
	0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f,
	0x0a, 0x0b, 0x6a, 0x73, 0x6f, 0x6e, 0x5f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0a, 0x6a, 0x73, 0x6f, 0x6e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x32,
	0x90, 0x01, 0x0a, 0x11, 0x4a, 0x53, 0x4f, 0x4e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x7b, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4a, 0x53, 0x4f, 0x4e,
	0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x31, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4a, 0x53, 0x4f, 0x4e, 0x53, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x32, 0x2e, 0x62, 0x75, 0x66, 0x2e,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4a, 0x53, 0x4f, 0x4e, 0x53,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x03, 0x90,
	0x02, 0x01, 0x42, 0x9c, 0x02, 0x0a, 0x1f, 0x63, 0x6f, 0x6d, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x42, 0x0f, 0x4a, 0x73, 0x6f, 0x6e, 0x73, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x59, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x66, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x62,
	0x75, 0x66, 0x2f, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x3b, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0xa2, 0x02, 0x03, 0x42, 0x41, 0x52, 0xaa, 0x02, 0x1b, 0x42, 0x75, 0x66,
	0x2e, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e,
	0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xca, 0x02, 0x1b, 0x42, 0x75, 0x66, 0x5c, 0x41,
	0x6c, 0x70, 0x68, 0x61, 0x5c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x5c, 0x56, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xe2, 0x02, 0x27, 0x42, 0x75, 0x66, 0x5c, 0x41, 0x6c, 0x70,
	0x68, 0x61, 0x5c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x5c, 0x56, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x1e, 0x42, 0x75, 0x66, 0x3a, 0x3a, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x3a, 0x3a, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_buf_alpha_registry_v1alpha1_jsonschema_proto_rawDescOnce sync.Once
	file_buf_alpha_registry_v1alpha1_jsonschema_proto_rawDescData = file_buf_alpha_registry_v1alpha1_jsonschema_proto_rawDesc
)

func file_buf_alpha_registry_v1alpha1_jsonschema_proto_rawDescGZIP() []byte {
	file_buf_alpha_registry_v1alpha1_jsonschema_proto_rawDescOnce.Do(func() {
		file_buf_alpha_registry_v1alpha1_jsonschema_proto_rawDescData = protoimpl.X.CompressGZIP(file_buf_alpha_registry_v1alpha1_jsonschema_proto_rawDescData)
	})
	return file_buf_alpha_registry_v1alpha1_jsonschema_proto_rawDescData
}

var file_buf_alpha_registry_v1alpha1_jsonschema_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_buf_alpha_registry_v1alpha1_jsonschema_proto_goTypes = []any{
	(*GetJSONSchemaRequest)(nil),  // 0: buf.alpha.registry.v1alpha1.GetJSONSchemaRequest
	(*GetJSONSchemaResponse)(nil), // 1: buf.alpha.registry.v1alpha1.GetJSONSchemaResponse
}
var file_buf_alpha_registry_v1alpha1_jsonschema_proto_depIdxs = []int32{
	0, // 0: buf.alpha.registry.v1alpha1.JSONSchemaService.GetJSONSchema:input_type -> buf.alpha.registry.v1alpha1.GetJSONSchemaRequest
	1, // 1: buf.alpha.registry.v1alpha1.JSONSchemaService.GetJSONSchema:output_type -> buf.alpha.registry.v1alpha1.GetJSONSchemaResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_buf_alpha_registry_v1alpha1_jsonschema_proto_init() }
func file_buf_alpha_registry_v1alpha1_jsonschema_proto_init() {
	if File_buf_alpha_registry_v1alpha1_jsonschema_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_buf_alpha_registry_v1alpha1_jsonschema_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_buf_alpha_registry_v1alpha1_jsonschema_proto_goTypes,
		DependencyIndexes: file_buf_alpha_registry_v1alpha1_jsonschema_proto_depIdxs,
		MessageInfos:      file_buf_alpha_registry_v1alpha1_jsonschema_proto_msgTypes,
	}.Build()
	File_buf_alpha_registry_v1alpha1_jsonschema_proto = out.File
	file_buf_alpha_registry_v1alpha1_jsonschema_proto_rawDesc = nil
	file_buf_alpha_registry_v1alpha1_jsonschema_proto_goTypes = nil
	file_buf_alpha_registry_v1alpha1_jsonschema_proto_depIdxs = nil
}
