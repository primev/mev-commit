// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: handshake/v1/handshake.proto

package handshakev1

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

type SerializedKeys struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PKEPublicKey  []byte `protobuf:"bytes,1,opt,name=PKEPublicKey,proto3" json:"PKEPublicKey,omitempty"`
	NIKEPublicKey []byte `protobuf:"bytes,2,opt,name=NIKEPublicKey,proto3" json:"NIKEPublicKey,omitempty"`
}

func (x *SerializedKeys) Reset() {
	*x = SerializedKeys{}
	if protoimpl.UnsafeEnabled {
		mi := &file_handshake_v1_handshake_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SerializedKeys) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SerializedKeys) ProtoMessage() {}

func (x *SerializedKeys) ProtoReflect() protoreflect.Message {
	mi := &file_handshake_v1_handshake_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SerializedKeys.ProtoReflect.Descriptor instead.
func (*SerializedKeys) Descriptor() ([]byte, []int) {
	return file_handshake_v1_handshake_proto_rawDescGZIP(), []int{0}
}

func (x *SerializedKeys) GetPKEPublicKey() []byte {
	if x != nil {
		return x.PKEPublicKey
	}
	return nil
}

func (x *SerializedKeys) GetNIKEPublicKey() []byte {
	if x != nil {
		return x.NIKEPublicKey
	}
	return nil
}

type HandshakeReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PeerType string          `protobuf:"bytes,1,opt,name=peer_type,json=peerType,proto3" json:"peer_type,omitempty"`
	Token    string          `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	Sig      []byte          `protobuf:"bytes,3,opt,name=sig,proto3" json:"sig,omitempty"`
	Keys     *SerializedKeys `protobuf:"bytes,4,opt,name=keys,proto3" json:"keys,omitempty"`
}

func (x *HandshakeReq) Reset() {
	*x = HandshakeReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_handshake_v1_handshake_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandshakeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandshakeReq) ProtoMessage() {}

func (x *HandshakeReq) ProtoReflect() protoreflect.Message {
	mi := &file_handshake_v1_handshake_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandshakeReq.ProtoReflect.Descriptor instead.
func (*HandshakeReq) Descriptor() ([]byte, []int) {
	return file_handshake_v1_handshake_proto_rawDescGZIP(), []int{1}
}

func (x *HandshakeReq) GetPeerType() string {
	if x != nil {
		return x.PeerType
	}
	return ""
}

func (x *HandshakeReq) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *HandshakeReq) GetSig() []byte {
	if x != nil {
		return x.Sig
	}
	return nil
}

func (x *HandshakeReq) GetKeys() *SerializedKeys {
	if x != nil {
		return x.Keys
	}
	return nil
}

type HandshakeResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ObservedAddress []byte `protobuf:"bytes,1,opt,name=observed_address,json=observedAddress,proto3" json:"observed_address,omitempty"`
	PeerType        string `protobuf:"bytes,2,opt,name=peer_type,json=peerType,proto3" json:"peer_type,omitempty"`
}

func (x *HandshakeResp) Reset() {
	*x = HandshakeResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_handshake_v1_handshake_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandshakeResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandshakeResp) ProtoMessage() {}

func (x *HandshakeResp) ProtoReflect() protoreflect.Message {
	mi := &file_handshake_v1_handshake_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandshakeResp.ProtoReflect.Descriptor instead.
func (*HandshakeResp) Descriptor() ([]byte, []int) {
	return file_handshake_v1_handshake_proto_rawDescGZIP(), []int{2}
}

func (x *HandshakeResp) GetObservedAddress() []byte {
	if x != nil {
		return x.ObservedAddress
	}
	return nil
}

func (x *HandshakeResp) GetPeerType() string {
	if x != nil {
		return x.PeerType
	}
	return ""
}

var File_handshake_v1_handshake_proto protoreflect.FileDescriptor

var file_handshake_v1_handshake_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x68,
	0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c,
	0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x2e, 0x76, 0x31, 0x22, 0x5a, 0x0a, 0x0e,
	0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x22,
	0x0a, 0x0c, 0x50, 0x4b, 0x45, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x50, 0x4b, 0x45, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b,
	0x65, 0x79, 0x12, 0x24, 0x0a, 0x0d, 0x4e, 0x49, 0x4b, 0x45, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x4e, 0x49, 0x4b, 0x45, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x22, 0x85, 0x01, 0x0a, 0x0c, 0x48, 0x61, 0x6e,
	0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x65, 0x65,
	0x72, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x65,
	0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x10, 0x0a, 0x03,
	0x73, 0x69, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x73, 0x69, 0x67, 0x12, 0x30,
	0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x68,
	0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x72, 0x69,
	0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x4b, 0x65, 0x79, 0x73, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73,
	0x22, 0x57, 0x0a, 0x0d, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x12, 0x29, 0x0a, 0x10, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x5f, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x6f, 0x62, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1b, 0x0a, 0x09,
	0x70, 0x65, 0x65, 0x72, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x65, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x42, 0xb5, 0x01, 0x0a, 0x10, 0x63, 0x6f,
	0x6d, 0x2e, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x0e,
	0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x69,
	0x6d, 0x65, 0x76, 0x2f, 0x6d, 0x65, 0x76, 0x2d, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x2f, 0x70,
	0x32, 0x70, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68,
	0x61, 0x6b, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65,
	0x76, 0x31, 0xa2, 0x02, 0x03, 0x48, 0x58, 0x58, 0xaa, 0x02, 0x0c, 0x48, 0x61, 0x6e, 0x64, 0x73,
	0x68, 0x61, 0x6b, 0x65, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0c, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68,
	0x61, 0x6b, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x18, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61,
	0x6b, 0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x0d, 0x48, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_handshake_v1_handshake_proto_rawDescOnce sync.Once
	file_handshake_v1_handshake_proto_rawDescData = file_handshake_v1_handshake_proto_rawDesc
)

func file_handshake_v1_handshake_proto_rawDescGZIP() []byte {
	file_handshake_v1_handshake_proto_rawDescOnce.Do(func() {
		file_handshake_v1_handshake_proto_rawDescData = protoimpl.X.CompressGZIP(file_handshake_v1_handshake_proto_rawDescData)
	})
	return file_handshake_v1_handshake_proto_rawDescData
}

var file_handshake_v1_handshake_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_handshake_v1_handshake_proto_goTypes = []interface{}{
	(*SerializedKeys)(nil), // 0: handshake.v1.SerializedKeys
	(*HandshakeReq)(nil),   // 1: handshake.v1.HandshakeReq
	(*HandshakeResp)(nil),  // 2: handshake.v1.HandshakeResp
}
var file_handshake_v1_handshake_proto_depIdxs = []int32{
	0, // 0: handshake.v1.HandshakeReq.keys:type_name -> handshake.v1.SerializedKeys
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_handshake_v1_handshake_proto_init() }
func file_handshake_v1_handshake_proto_init() {
	if File_handshake_v1_handshake_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_handshake_v1_handshake_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SerializedKeys); i {
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
		file_handshake_v1_handshake_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandshakeReq); i {
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
		file_handshake_v1_handshake_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandshakeResp); i {
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
			RawDescriptor: file_handshake_v1_handshake_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_handshake_v1_handshake_proto_goTypes,
		DependencyIndexes: file_handshake_v1_handshake_proto_depIdxs,
		MessageInfos:      file_handshake_v1_handshake_proto_msgTypes,
	}.Build()
	File_handshake_v1_handshake_proto = out.File
	file_handshake_v1_handshake_proto_rawDesc = nil
	file_handshake_v1_handshake_proto_goTypes = nil
	file_handshake_v1_handshake_proto_depIdxs = nil
}
