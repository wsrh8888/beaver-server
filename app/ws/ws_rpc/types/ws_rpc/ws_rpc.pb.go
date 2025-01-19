// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.19.4
// source: ws_rpc.proto

package ws_rpc

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

type HandleWebSocketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID string `protobuf:"bytes,1,opt,name=user_id,json=userID,proto3" json:"user_id,omitempty"`
	Source string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
}

func (x *HandleWebSocketRequest) Reset() {
	*x = HandleWebSocketRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ws_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandleWebSocketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandleWebSocketRequest) ProtoMessage() {}

func (x *HandleWebSocketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ws_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandleWebSocketRequest.ProtoReflect.Descriptor instead.
func (*HandleWebSocketRequest) Descriptor() ([]byte, []int) {
	return file_ws_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *HandleWebSocketRequest) GetUserId() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *HandleWebSocketRequest) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

type HandleWebSocketResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *HandleWebSocketResponse) Reset() {
	*x = HandleWebSocketResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ws_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandleWebSocketResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandleWebSocketResponse) ProtoMessage() {}

func (x *HandleWebSocketResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ws_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandleWebSocketResponse.ProtoReflect.Descriptor instead.
func (*HandleWebSocketResponse) Descriptor() ([]byte, []int) {
	return file_ws_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *HandleWebSocketResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type SendProxyMessageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID      string            `protobuf:"bytes,1,opt,name=user_id,json=userID,proto3" json:"user_id,omitempty"`
	Command     string            `protobuf:"bytes,2,opt,name=command,proto3" json:"command,omitempty"`
	TargetID    string            `protobuf:"bytes,3,opt,name=target_id,json=targetID,proto3" json:"target_id,omitempty"`
	MessageType string            `protobuf:"bytes,4,opt,name=message_type,json=messageType,proto3" json:"message_type,omitempty"`
	Body        map[string]string `protobuf:"bytes,5,rep,name=body,proto3" json:"body,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *SendProxyMessageRequest) Reset() {
	*x = SendProxyMessageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ws_rpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendProxyMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendProxyMessageRequest) ProtoMessage() {}

func (x *SendProxyMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ws_rpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendProxyMessageRequest.ProtoReflect.Descriptor instead.
func (*SendProxyMessageRequest) Descriptor() ([]byte, []int) {
	return file_ws_rpc_proto_rawDescGZIP(), []int{2}
}

func (x *SendProxyMessageRequest) GetUserId() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *SendProxyMessageRequest) GetCommand() string {
	if x != nil {
		return x.Command
	}
	return ""
}

func (x *SendProxyMessageRequest) GetTargetId() string {
	if x != nil {
		return x.TargetID
	}
	return ""
}

func (x *SendProxyMessageRequest) GetMessageType() string {
	if x != nil {
		return x.MessageType
	}
	return ""
}

func (x *SendProxyMessageRequest) GetBody() map[string]string {
	if x != nil {
		return x.Body
	}
	return nil
}

type SendProxyMessageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *SendProxyMessageResponse) Reset() {
	*x = SendProxyMessageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ws_rpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendProxyMessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendProxyMessageResponse) ProtoMessage() {}

func (x *SendProxyMessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ws_rpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendProxyMessageResponse.ProtoReflect.Descriptor instead.
func (*SendProxyMessageResponse) Descriptor() ([]byte, []int) {
	return file_ws_rpc_proto_rawDescGZIP(), []int{3}
}

func (x *SendProxyMessageResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_ws_rpc_proto protoreflect.FileDescriptor

var file_ws_rpc_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x77, 0x73, 0x5f, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x77, 0x73, 0x5f, 0x72, 0x70, 0x63, 0x22, 0x49, 0x0a, 0x16, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x57, 0x65, 0x62, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x22, 0x31, 0x0a, 0x17, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x57, 0x65, 0x62, 0x53, 0x6f,
	0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x22, 0x84, 0x02, 0x0a, 0x17, 0x53, 0x65, 0x6e, 0x64, 0x50, 0x72, 0x6f,
	0x78, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x49, 0x64,
	0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x3d, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x29, 0x2e, 0x77, 0x73, 0x5f, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x50,
	0x72, 0x6f, 0x78, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x2e, 0x42, 0x6f, 0x64, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x62, 0x6f,
	0x64, 0x79, 0x1a, 0x37, 0x0a, 0x09, 0x42, 0x6f, 0x64, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x32, 0x0a, 0x18, 0x53,
	0x65, 0x6e, 0x64, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32,
	0xaf, 0x01, 0x0a, 0x02, 0x57, 0x73, 0x12, 0x52, 0x0a, 0x0f, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x57, 0x65, 0x62, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x1e, 0x2e, 0x77, 0x73, 0x5f, 0x72,
	0x70, 0x63, 0x2e, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x57, 0x65, 0x62, 0x53, 0x6f, 0x63, 0x6b,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x77, 0x73, 0x5f, 0x72,
	0x70, 0x63, 0x2e, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x57, 0x65, 0x62, 0x53, 0x6f, 0x63, 0x6b,
	0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55, 0x0a, 0x10, 0x53, 0x65,
	0x6e, 0x64, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f,
	0x2e, 0x77, 0x73, 0x5f, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x50, 0x72, 0x6f, 0x78,
	0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x20, 0x2e, 0x77, 0x73, 0x5f, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x50, 0x72, 0x6f,
	0x78, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x77, 0x73, 0x5f, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ws_rpc_proto_rawDescOnce sync.Once
	file_ws_rpc_proto_rawDescData = file_ws_rpc_proto_rawDesc
)

func file_ws_rpc_proto_rawDescGZIP() []byte {
	file_ws_rpc_proto_rawDescOnce.Do(func() {
		file_ws_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_ws_rpc_proto_rawDescData)
	})
	return file_ws_rpc_proto_rawDescData
}

var file_ws_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_ws_rpc_proto_goTypes = []interface{}{
	(*HandleWebSocketRequest)(nil),   // 0: ws_rpc.HandleWebSocketRequest
	(*HandleWebSocketResponse)(nil),  // 1: ws_rpc.HandleWebSocketResponse
	(*SendProxyMessageRequest)(nil),  // 2: ws_rpc.SendProxyMessageRequest
	(*SendProxyMessageResponse)(nil), // 3: ws_rpc.SendProxyMessageResponse
	nil,                              // 4: ws_rpc.SendProxyMessageRequest.BodyEntry
}
var file_ws_rpc_proto_depIdxs = []int32{
	4, // 0: ws_rpc.SendProxyMessageRequest.body:type_name -> ws_rpc.SendProxyMessageRequest.BodyEntry
	0, // 1: ws_rpc.Ws.HandleWebSocket:input_type -> ws_rpc.HandleWebSocketRequest
	2, // 2: ws_rpc.Ws.SendProxyMessage:input_type -> ws_rpc.SendProxyMessageRequest
	1, // 3: ws_rpc.Ws.HandleWebSocket:output_type -> ws_rpc.HandleWebSocketResponse
	3, // 4: ws_rpc.Ws.SendProxyMessage:output_type -> ws_rpc.SendProxyMessageResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_ws_rpc_proto_init() }
func file_ws_rpc_proto_init() {
	if File_ws_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ws_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandleWebSocketRequest); i {
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
		file_ws_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandleWebSocketResponse); i {
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
		file_ws_rpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendProxyMessageRequest); i {
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
		file_ws_rpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendProxyMessageResponse); i {
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
			RawDescriptor: file_ws_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ws_rpc_proto_goTypes,
		DependencyIndexes: file_ws_rpc_proto_depIdxs,
		MessageInfos:      file_ws_rpc_proto_msgTypes,
	}.Build()
	File_ws_rpc_proto = out.File
	file_ws_rpc_proto_rawDesc = nil
	file_ws_rpc_proto_goTypes = nil
	file_ws_rpc_proto_depIdxs = nil
}
