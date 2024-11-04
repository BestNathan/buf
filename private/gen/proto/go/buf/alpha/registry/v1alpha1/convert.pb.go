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
// source: buf/alpha/registry/v1alpha1/convert.proto

package registryv1alpha1

import (
	v1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/image/v1"
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

// The supported formats for the serialized message conversion.
type ConvertFormat int32

const (
	ConvertFormat_CONVERT_FORMAT_UNSPECIFIED ConvertFormat = 0
	ConvertFormat_CONVERT_FORMAT_BIN         ConvertFormat = 1
	ConvertFormat_CONVERT_FORMAT_JSON        ConvertFormat = 2
)

// Enum value maps for ConvertFormat.
var (
	ConvertFormat_name = map[int32]string{
		0: "CONVERT_FORMAT_UNSPECIFIED",
		1: "CONVERT_FORMAT_BIN",
		2: "CONVERT_FORMAT_JSON",
	}
	ConvertFormat_value = map[string]int32{
		"CONVERT_FORMAT_UNSPECIFIED": 0,
		"CONVERT_FORMAT_BIN":         1,
		"CONVERT_FORMAT_JSON":        2,
	}
)

func (x ConvertFormat) Enum() *ConvertFormat {
	p := new(ConvertFormat)
	*p = x
	return p
}

func (x ConvertFormat) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConvertFormat) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_alpha_registry_v1alpha1_convert_proto_enumTypes[0].Descriptor()
}

func (ConvertFormat) Type() protoreflect.EnumType {
	return &file_buf_alpha_registry_v1alpha1_convert_proto_enumTypes[0]
}

func (x ConvertFormat) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConvertFormat.Descriptor instead.
func (ConvertFormat) EnumDescriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_convert_proto_rawDescGZIP(), []int{0}
}

type ConvertRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// type_name is the full type name of the serialized message (like acme.weather.v1.Units).
	TypeName string `protobuf:"bytes,1,opt,name=type_name,json=typeName,proto3" json:"type_name,omitempty"`
	// image is the image source that defines the serialized message.
	Image *v1.Image `protobuf:"bytes,2,opt,name=image,proto3" json:"image,omitempty"`
	// payload is the serialized Protobuf message.
	Payload []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	// request_format is the format of the payload.
	RequestFormat ConvertFormat `protobuf:"varint,4,opt,name=request_format,json=requestFormat,proto3,enum=buf.alpha.registry.v1alpha1.ConvertFormat" json:"request_format,omitempty"`
	// response_format is the desired format of the output result.
	ResponseFormat ConvertFormat `protobuf:"varint,5,opt,name=response_format,json=responseFormat,proto3,enum=buf.alpha.registry.v1alpha1.ConvertFormat" json:"response_format,omitempty"`
}

func (x *ConvertRequest) Reset() {
	*x = ConvertRequest{}
	mi := &file_buf_alpha_registry_v1alpha1_convert_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConvertRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConvertRequest) ProtoMessage() {}

func (x *ConvertRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_convert_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConvertRequest.ProtoReflect.Descriptor instead.
func (*ConvertRequest) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_convert_proto_rawDescGZIP(), []int{0}
}

func (x *ConvertRequest) GetTypeName() string {
	if x != nil {
		return x.TypeName
	}
	return ""
}

func (x *ConvertRequest) GetImage() *v1.Image {
	if x != nil {
		return x.Image
	}
	return nil
}

func (x *ConvertRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *ConvertRequest) GetRequestFormat() ConvertFormat {
	if x != nil {
		return x.RequestFormat
	}
	return ConvertFormat_CONVERT_FORMAT_UNSPECIFIED
}

func (x *ConvertRequest) GetResponseFormat() ConvertFormat {
	if x != nil {
		return x.ResponseFormat
	}
	return ConvertFormat_CONVERT_FORMAT_UNSPECIFIED
}

type ConvertResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// payload is the converted serialized message in one of the supported formats.
	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *ConvertResponse) Reset() {
	*x = ConvertResponse{}
	mi := &file_buf_alpha_registry_v1alpha1_convert_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConvertResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConvertResponse) ProtoMessage() {}

func (x *ConvertResponse) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_convert_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConvertResponse.ProtoReflect.Descriptor instead.
func (*ConvertResponse) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_convert_proto_rawDescGZIP(), []int{1}
}

func (x *ConvertResponse) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

var File_buf_alpha_registry_v1alpha1_convert_proto protoreflect.FileDescriptor

var file_buf_alpha_registry_v1alpha1_convert_proto_rawDesc = []byte{
	0x0a, 0x29, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2f, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x63, 0x6f,
	0x6e, 0x76, 0x65, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x62, 0x75, 0x66,
	0x2e, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x1a, 0x1e, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x2f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa0, 0x02, 0x0a, 0x0e, 0x43, 0x6f, 0x6e,
	0x76, 0x65, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x74,
	0x79, 0x70, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x74, 0x79, 0x70, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2f, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x2e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x12, 0x51, 0x0a, 0x0e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x66,
	0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2a, 0x2e, 0x62, 0x75,
	0x66, 0x2e, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72,
	0x74, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x52, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x53, 0x0a, 0x0f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x5f, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x2a, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2e, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x43, 0x6f,
	0x6e, 0x76, 0x65, 0x72, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x52, 0x0e, 0x72, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x22, 0x2b, 0x0a, 0x0f, 0x43,
	0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2a, 0x60, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x76,
	0x65, 0x72, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x1e, 0x0a, 0x1a, 0x43, 0x4f, 0x4e,
	0x56, 0x45, 0x52, 0x54, 0x5f, 0x46, 0x4f, 0x52, 0x4d, 0x41, 0x54, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x4f, 0x4e,
	0x56, 0x45, 0x52, 0x54, 0x5f, 0x46, 0x4f, 0x52, 0x4d, 0x41, 0x54, 0x5f, 0x42, 0x49, 0x4e, 0x10,
	0x01, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x4f, 0x4e, 0x56, 0x45, 0x52, 0x54, 0x5f, 0x46, 0x4f, 0x52,
	0x4d, 0x41, 0x54, 0x5f, 0x4a, 0x53, 0x4f, 0x4e, 0x10, 0x02, 0x32, 0x76, 0x0a, 0x0e, 0x43, 0x6f,
	0x6e, 0x76, 0x65, 0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x64, 0x0a, 0x07,
	0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x12, 0x2b, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x99, 0x02, 0x0a, 0x1f, 0x63, 0x6f, 0x6d, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x42, 0x0c, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x74, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x59, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x66, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x62, 0x75, 0x66, 0x2f,
	0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2f, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x3b, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0xa2, 0x02, 0x03, 0x42, 0x41, 0x52, 0xaa, 0x02, 0x1b, 0x42, 0x75, 0x66, 0x2e, 0x41, 0x6c,
	0x70, 0x68, 0x61, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x56, 0x31, 0x61,
	0x6c, 0x70, 0x68, 0x61, 0x31, 0xca, 0x02, 0x1b, 0x42, 0x75, 0x66, 0x5c, 0x41, 0x6c, 0x70, 0x68,
	0x61, 0x5c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0xe2, 0x02, 0x27, 0x42, 0x75, 0x66, 0x5c, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x5c,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1e,
	0x42, 0x75, 0x66, 0x3a, 0x3a, 0x41, 0x6c, 0x70, 0x68, 0x61, 0x3a, 0x3a, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x79, 0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_buf_alpha_registry_v1alpha1_convert_proto_rawDescOnce sync.Once
	file_buf_alpha_registry_v1alpha1_convert_proto_rawDescData = file_buf_alpha_registry_v1alpha1_convert_proto_rawDesc
)

func file_buf_alpha_registry_v1alpha1_convert_proto_rawDescGZIP() []byte {
	file_buf_alpha_registry_v1alpha1_convert_proto_rawDescOnce.Do(func() {
		file_buf_alpha_registry_v1alpha1_convert_proto_rawDescData = protoimpl.X.CompressGZIP(file_buf_alpha_registry_v1alpha1_convert_proto_rawDescData)
	})
	return file_buf_alpha_registry_v1alpha1_convert_proto_rawDescData
}

var file_buf_alpha_registry_v1alpha1_convert_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_buf_alpha_registry_v1alpha1_convert_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_buf_alpha_registry_v1alpha1_convert_proto_goTypes = []any{
	(ConvertFormat)(0),      // 0: buf.alpha.registry.v1alpha1.ConvertFormat
	(*ConvertRequest)(nil),  // 1: buf.alpha.registry.v1alpha1.ConvertRequest
	(*ConvertResponse)(nil), // 2: buf.alpha.registry.v1alpha1.ConvertResponse
	(*v1.Image)(nil),        // 3: buf.alpha.image.v1.Image
}
var file_buf_alpha_registry_v1alpha1_convert_proto_depIdxs = []int32{
	3, // 0: buf.alpha.registry.v1alpha1.ConvertRequest.image:type_name -> buf.alpha.image.v1.Image
	0, // 1: buf.alpha.registry.v1alpha1.ConvertRequest.request_format:type_name -> buf.alpha.registry.v1alpha1.ConvertFormat
	0, // 2: buf.alpha.registry.v1alpha1.ConvertRequest.response_format:type_name -> buf.alpha.registry.v1alpha1.ConvertFormat
	1, // 3: buf.alpha.registry.v1alpha1.ConvertService.Convert:input_type -> buf.alpha.registry.v1alpha1.ConvertRequest
	2, // 4: buf.alpha.registry.v1alpha1.ConvertService.Convert:output_type -> buf.alpha.registry.v1alpha1.ConvertResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_buf_alpha_registry_v1alpha1_convert_proto_init() }
func file_buf_alpha_registry_v1alpha1_convert_proto_init() {
	if File_buf_alpha_registry_v1alpha1_convert_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_buf_alpha_registry_v1alpha1_convert_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_buf_alpha_registry_v1alpha1_convert_proto_goTypes,
		DependencyIndexes: file_buf_alpha_registry_v1alpha1_convert_proto_depIdxs,
		EnumInfos:         file_buf_alpha_registry_v1alpha1_convert_proto_enumTypes,
		MessageInfos:      file_buf_alpha_registry_v1alpha1_convert_proto_msgTypes,
	}.Build()
	File_buf_alpha_registry_v1alpha1_convert_proto = out.File
	file_buf_alpha_registry_v1alpha1_convert_proto_rawDesc = nil
	file_buf_alpha_registry_v1alpha1_convert_proto_goTypes = nil
	file_buf_alpha_registry_v1alpha1_convert_proto_depIdxs = nil
}
