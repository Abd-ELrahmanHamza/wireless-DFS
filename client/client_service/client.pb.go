// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.3
// source: protobufs/client.proto

package client_service

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

type UploadingCompletionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UploadingCompletionRequest) Reset() {
	*x = UploadingCompletionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobufs_client_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadingCompletionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadingCompletionRequest) ProtoMessage() {}

func (x *UploadingCompletionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protobufs_client_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadingCompletionRequest.ProtoReflect.Descriptor instead.
func (*UploadingCompletionRequest) Descriptor() ([]byte, []int) {
	return file_protobufs_client_proto_rawDescGZIP(), []int{0}
}

type UploadingCompletionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UploadingCompletionResponse) Reset() {
	*x = UploadingCompletionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobufs_client_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadingCompletionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadingCompletionResponse) ProtoMessage() {}

func (x *UploadingCompletionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protobufs_client_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadingCompletionResponse.ProtoReflect.Descriptor instead.
func (*UploadingCompletionResponse) Descriptor() ([]byte, []int) {
	return file_protobufs_client_proto_rawDescGZIP(), []int{1}
}

var File_protobufs_client_proto protoreflect.FileDescriptor

var file_protobufs_client_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2f, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1c, 0x0a, 0x1a, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x1d, 0x0a, 0x1b, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x61, 0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x50, 0x0a, 0x13, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x2e,
	0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x11, 0x5a, 0x0f, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_protobufs_client_proto_rawDescOnce sync.Once
	file_protobufs_client_proto_rawDescData = file_protobufs_client_proto_rawDesc
)

func file_protobufs_client_proto_rawDescGZIP() []byte {
	file_protobufs_client_proto_rawDescOnce.Do(func() {
		file_protobufs_client_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobufs_client_proto_rawDescData)
	})
	return file_protobufs_client_proto_rawDescData
}

var file_protobufs_client_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_protobufs_client_proto_goTypes = []interface{}{
	(*UploadingCompletionRequest)(nil),  // 0: UploadingCompletionRequest
	(*UploadingCompletionResponse)(nil), // 1: UploadingCompletionResponse
}
var file_protobufs_client_proto_depIdxs = []int32{
	0, // 0: ClientService.UploadingCompletion:input_type -> UploadingCompletionRequest
	1, // 1: ClientService.UploadingCompletion:output_type -> UploadingCompletionResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protobufs_client_proto_init() }
func file_protobufs_client_proto_init() {
	if File_protobufs_client_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protobufs_client_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadingCompletionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protobufs_client_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadingCompletionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protobufs_client_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protobufs_client_proto_goTypes,
		DependencyIndexes: file_protobufs_client_proto_depIdxs,
		MessageInfos:      file_protobufs_client_proto_msgTypes,
	}.Build()
	File_protobufs_client_proto = out.File
	file_protobufs_client_proto_rawDesc = nil
	file_protobufs_client_proto_goTypes = nil
	file_protobufs_client_proto_depIdxs = nil
}
