// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.11.4
// source: github.com/moby/buildkit/session/exporter/exporter.proto

package exporter

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

type FindExportersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata map[string][]byte `protobuf:"bytes,1,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Refs     []string          `protobuf:"bytes,2,rep,name=refs,proto3" json:"refs,omitempty"`
}

func (x *FindExportersRequest) Reset() {
	*x = FindExportersRequest{}
	mi := &file_github_com_moby_buildkit_session_exporter_exporter_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FindExportersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindExportersRequest) ProtoMessage() {}

func (x *FindExportersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_moby_buildkit_session_exporter_exporter_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindExportersRequest.ProtoReflect.Descriptor instead.
func (*FindExportersRequest) Descriptor() ([]byte, []int) {
	return file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDescGZIP(), []int{0}
}

func (x *FindExportersRequest) GetMetadata() map[string][]byte {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *FindExportersRequest) GetRefs() []string {
	if x != nil {
		return x.Refs
	}
	return nil
}

type FindExportersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Exporters []*ExporterRequest `protobuf:"bytes,1,rep,name=exporters,proto3" json:"exporters,omitempty"`
}

func (x *FindExportersResponse) Reset() {
	*x = FindExportersResponse{}
	mi := &file_github_com_moby_buildkit_session_exporter_exporter_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FindExportersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindExportersResponse) ProtoMessage() {}

func (x *FindExportersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_moby_buildkit_session_exporter_exporter_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindExportersResponse.ProtoReflect.Descriptor instead.
func (*FindExportersResponse) Descriptor() ([]byte, []int) {
	return file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDescGZIP(), []int{1}
}

func (x *FindExportersResponse) GetExporters() []*ExporterRequest {
	if x != nil {
		return x.Exporters
	}
	return nil
}

type ExporterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  string            `protobuf:"bytes,1,opt,name=Type,proto3" json:"Type,omitempty"`
	Attrs map[string]string `protobuf:"bytes,2,rep,name=Attrs,proto3" json:"Attrs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ExporterRequest) Reset() {
	*x = ExporterRequest{}
	mi := &file_github_com_moby_buildkit_session_exporter_exporter_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExporterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExporterRequest) ProtoMessage() {}

func (x *ExporterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_moby_buildkit_session_exporter_exporter_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExporterRequest.ProtoReflect.Descriptor instead.
func (*ExporterRequest) Descriptor() ([]byte, []int) {
	return file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDescGZIP(), []int{2}
}

func (x *ExporterRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *ExporterRequest) GetAttrs() map[string]string {
	if x != nil {
		return x.Attrs
	}
	return nil
}

var File_github_com_moby_buildkit_session_exporter_exporter_proto protoreflect.FileDescriptor

var file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDesc = []byte{
	0x0a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x62,
	0x79, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x6b, 0x69, 0x74, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x2f, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x65, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x6d, 0x6f, 0x62, 0x79,
	0x2e, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x22, 0xb9, 0x01, 0x0a,
	0x14, 0x46, 0x69, 0x6e, 0x64, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x50, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x6d, 0x6f, 0x62, 0x79, 0x2e, 0x65,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x45,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x65, 0x66, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x72, 0x65, 0x66, 0x73, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x58, 0x0a, 0x15, 0x46, 0x69, 0x6e, 0x64,
	0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x3f, 0x0a, 0x09, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6d, 0x6f, 0x62, 0x79, 0x2e, 0x65, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x09, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65,
	0x72, 0x73, 0x22, 0xa3, 0x01, 0x0a, 0x0f, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x42, 0x0a, 0x05, 0x41, 0x74,
	0x74, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x6d, 0x6f, 0x62, 0x79,
	0x2e, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x74, 0x74,
	0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x41, 0x74, 0x74, 0x72, 0x73, 0x1a, 0x38,
	0x0a, 0x0a, 0x41, 0x74, 0x74, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x6c, 0x0a, 0x08, 0x45, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x65, 0x72, 0x12, 0x60, 0x0a, 0x0d, 0x46, 0x69, 0x6e, 0x64, 0x45, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x65, 0x72, 0x73, 0x12, 0x26, 0x2e, 0x6d, 0x6f, 0x62, 0x79, 0x2e, 0x65, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x45, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e,
	0x6d, 0x6f, 0x62, 0x79, 0x2e, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x46, 0x69, 0x6e, 0x64, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x62, 0x79, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x6b,
	0x69, 0x74, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x65, 0x78, 0x70, 0x6f, 0x72,
	0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDescOnce sync.Once
	file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDescData = file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDesc
)

func file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDescGZIP() []byte {
	file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDescOnce.Do(func() {
		file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDescData)
	})
	return file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDescData
}

var file_github_com_moby_buildkit_session_exporter_exporter_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_github_com_moby_buildkit_session_exporter_exporter_proto_goTypes = []any{
	(*FindExportersRequest)(nil),  // 0: moby.exporter.v1.FindExportersRequest
	(*FindExportersResponse)(nil), // 1: moby.exporter.v1.FindExportersResponse
	(*ExporterRequest)(nil),       // 2: moby.exporter.v1.ExporterRequest
	nil,                           // 3: moby.exporter.v1.FindExportersRequest.MetadataEntry
	nil,                           // 4: moby.exporter.v1.ExporterRequest.AttrsEntry
}
var file_github_com_moby_buildkit_session_exporter_exporter_proto_depIdxs = []int32{
	3, // 0: moby.exporter.v1.FindExportersRequest.metadata:type_name -> moby.exporter.v1.FindExportersRequest.MetadataEntry
	2, // 1: moby.exporter.v1.FindExportersResponse.exporters:type_name -> moby.exporter.v1.ExporterRequest
	4, // 2: moby.exporter.v1.ExporterRequest.Attrs:type_name -> moby.exporter.v1.ExporterRequest.AttrsEntry
	0, // 3: moby.exporter.v1.Exporter.FindExporters:input_type -> moby.exporter.v1.FindExportersRequest
	1, // 4: moby.exporter.v1.Exporter.FindExporters:output_type -> moby.exporter.v1.FindExportersResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_github_com_moby_buildkit_session_exporter_exporter_proto_init() }
func file_github_com_moby_buildkit_session_exporter_exporter_proto_init() {
	if File_github_com_moby_buildkit_session_exporter_exporter_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_moby_buildkit_session_exporter_exporter_proto_goTypes,
		DependencyIndexes: file_github_com_moby_buildkit_session_exporter_exporter_proto_depIdxs,
		MessageInfos:      file_github_com_moby_buildkit_session_exporter_exporter_proto_msgTypes,
	}.Build()
	File_github_com_moby_buildkit_session_exporter_exporter_proto = out.File
	file_github_com_moby_buildkit_session_exporter_exporter_proto_rawDesc = nil
	file_github_com_moby_buildkit_session_exporter_exporter_proto_goTypes = nil
	file_github_com_moby_buildkit_session_exporter_exporter_proto_depIdxs = nil
}
