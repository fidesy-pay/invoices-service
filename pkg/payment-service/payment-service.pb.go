// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.3
// source: api/payment-service/payment-service.proto

package payment_service

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

type Chain int32

const (
	Chain_UNKNOWN_CHAIN Chain = 0
	Chain_POLYGON       Chain = 1
)

// Enum value maps for Chain.
var (
	Chain_name = map[int32]string{
		0: "UNKNOWN_CHAIN",
		1: "POLYGON",
	}
	Chain_value = map[string]int32{
		"UNKNOWN_CHAIN": 0,
		"POLYGON":       1,
	}
)

func (x Chain) Enum() *Chain {
	p := new(Chain)
	*p = x
	return p
}

func (x Chain) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Chain) Descriptor() protoreflect.EnumDescriptor {
	return file_api_payment_service_payment_service_proto_enumTypes[0].Descriptor()
}

func (Chain) Type() protoreflect.EnumType {
	return &file_api_payment_service_payment_service_proto_enumTypes[0]
}

func (x Chain) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Chain.Descriptor instead.
func (Chain) EnumDescriptor() ([]byte, []int) {
	return file_api_payment_service_payment_service_proto_rawDescGZIP(), []int{0}
}

type Token int32

const (
	Token_UNKNOWN_TOKEN Token = 0
	Token_MATIC         Token = 1
)

// Enum value maps for Token.
var (
	Token_name = map[int32]string{
		0: "UNKNOWN_TOKEN",
		1: "MATIC",
	}
	Token_value = map[string]int32{
		"UNKNOWN_TOKEN": 0,
		"MATIC":         1,
	}
)

func (x Token) Enum() *Token {
	p := new(Token)
	*p = x
	return p
}

func (x Token) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Token) Descriptor() protoreflect.EnumDescriptor {
	return file_api_payment_service_payment_service_proto_enumTypes[1].Descriptor()
}

func (Token) Type() protoreflect.EnumType {
	return &file_api_payment_service_payment_service_proto_enumTypes[1]
}

func (x Token) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Token.Descriptor instead.
func (Token) EnumDescriptor() ([]byte, []int) {
	return file_api_payment_service_payment_service_proto_rawDescGZIP(), []int{1}
}

type PaymentStatus int32

const (
	PaymentStatus_UNKNOWN_STATUS PaymentStatus = 0
	PaymentStatus_PENDING        PaymentStatus = 1
	PaymentStatus_FAILED         PaymentStatus = 2
	PaymentStatus_SUCCESS        PaymentStatus = 3
)

// Enum value maps for PaymentStatus.
var (
	PaymentStatus_name = map[int32]string{
		0: "UNKNOWN_STATUS",
		1: "PENDING",
		2: "FAILED",
		3: "SUCCESS",
	}
	PaymentStatus_value = map[string]int32{
		"UNKNOWN_STATUS": 0,
		"PENDING":        1,
		"FAILED":         2,
		"SUCCESS":        3,
	}
)

func (x PaymentStatus) Enum() *PaymentStatus {
	p := new(PaymentStatus)
	*p = x
	return p
}

func (x PaymentStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PaymentStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_payment_service_payment_service_proto_enumTypes[2].Descriptor()
}

func (PaymentStatus) Type() protoreflect.EnumType {
	return &file_api_payment_service_payment_service_proto_enumTypes[2]
}

func (x PaymentStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PaymentStatus.Descriptor instead.
func (PaymentStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_payment_service_payment_service_proto_rawDescGZIP(), []int{2}
}

type CreatePaymentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Amount float64 `protobuf:"fixed64,1,opt,name=amount,proto3" json:"amount,omitempty"`
	Chain  Chain   `protobuf:"varint,2,opt,name=chain,proto3,enum=payment_service.Chain" json:"chain,omitempty"`
	Token  Token   `protobuf:"varint,3,opt,name=token,proto3,enum=payment_service.Token" json:"token,omitempty"`
}

func (x *CreatePaymentRequest) Reset() {
	*x = CreatePaymentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_payment_service_payment_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePaymentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePaymentRequest) ProtoMessage() {}

func (x *CreatePaymentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_payment_service_payment_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePaymentRequest.ProtoReflect.Descriptor instead.
func (*CreatePaymentRequest) Descriptor() ([]byte, []int) {
	return file_api_payment_service_payment_service_proto_rawDescGZIP(), []int{0}
}

func (x *CreatePaymentRequest) GetAmount() float64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *CreatePaymentRequest) GetChain() Chain {
	if x != nil {
		return x.Chain
	}
	return Chain_UNKNOWN_CHAIN
}

func (x *CreatePaymentRequest) GetToken() Token {
	if x != nil {
		return x.Token
	}
	return Token_UNKNOWN_TOKEN
}

type CreatePaymentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// payment identifier
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// address
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *CreatePaymentResponse) Reset() {
	*x = CreatePaymentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_payment_service_payment_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePaymentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePaymentResponse) ProtoMessage() {}

func (x *CreatePaymentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_payment_service_payment_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePaymentResponse.ProtoReflect.Descriptor instead.
func (*CreatePaymentResponse) Descriptor() ([]byte, []int) {
	return file_api_payment_service_payment_service_proto_rawDescGZIP(), []int{1}
}

func (x *CreatePaymentResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CreatePaymentResponse) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type CheckPaymentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// payment identifier
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CheckPaymentRequest) Reset() {
	*x = CheckPaymentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_payment_service_payment_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckPaymentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckPaymentRequest) ProtoMessage() {}

func (x *CheckPaymentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_payment_service_payment_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckPaymentRequest.ProtoReflect.Descriptor instead.
func (*CheckPaymentRequest) Descriptor() ([]byte, []int) {
	return file_api_payment_service_payment_service_proto_rawDescGZIP(), []int{2}
}

func (x *CheckPaymentRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CheckPaymentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status PaymentStatus `protobuf:"varint,1,opt,name=status,proto3,enum=payment_service.PaymentStatus" json:"status,omitempty"`
}

func (x *CheckPaymentResponse) Reset() {
	*x = CheckPaymentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_payment_service_payment_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckPaymentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckPaymentResponse) ProtoMessage() {}

func (x *CheckPaymentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_payment_service_payment_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckPaymentResponse.ProtoReflect.Descriptor instead.
func (*CheckPaymentResponse) Descriptor() ([]byte, []int) {
	return file_api_payment_service_payment_service_proto_rawDescGZIP(), []int{3}
}

func (x *CheckPaymentResponse) GetStatus() PaymentStatus {
	if x != nil {
		return x.Status
	}
	return PaymentStatus_UNKNOWN_STATUS
}

var File_api_payment_service_payment_service_proto protoreflect.FileDescriptor

var file_api_payment_service_payment_service_proto_rawDesc = []byte{
	0x0a, 0x29, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x70, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x8a, 0x01, 0x0a,
	0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2c, 0x0a,
	0x05, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x70,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43,
	0x68, 0x61, 0x69, 0x6e, 0x52, 0x05, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x2c, 0x0a, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x70, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x41, 0x0a, 0x15, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x25, 0x0a, 0x13,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x4e, 0x0a, 0x14, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x50, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x70, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x2a, 0x27, 0x0a, 0x05, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x11, 0x0a, 0x0d,
	0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x43, 0x48, 0x41, 0x49, 0x4e, 0x10, 0x00, 0x12,
	0x0b, 0x0a, 0x07, 0x50, 0x4f, 0x4c, 0x59, 0x47, 0x4f, 0x4e, 0x10, 0x01, 0x2a, 0x25, 0x0a, 0x05,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e,
	0x5f, 0x54, 0x4f, 0x4b, 0x45, 0x4e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x41, 0x54, 0x49,
	0x43, 0x10, 0x01, 0x2a, 0x49, 0x0a, 0x0d, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x0e, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x45, 0x4e, 0x44,
	0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10,
	0x02, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x03, 0x32, 0xcd,
	0x01, 0x0a, 0x0e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x5e, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x12, 0x25, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x61, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x70, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x5b, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x12, 0x24, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x50,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x37,
	0x5a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x64,
	0x65, 0x73, 0x79, 0x2d, 0x70, 0x61, 0x79, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2d,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x3b, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_payment_service_payment_service_proto_rawDescOnce sync.Once
	file_api_payment_service_payment_service_proto_rawDescData = file_api_payment_service_payment_service_proto_rawDesc
)

func file_api_payment_service_payment_service_proto_rawDescGZIP() []byte {
	file_api_payment_service_payment_service_proto_rawDescOnce.Do(func() {
		file_api_payment_service_payment_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_payment_service_payment_service_proto_rawDescData)
	})
	return file_api_payment_service_payment_service_proto_rawDescData
}

var file_api_payment_service_payment_service_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_api_payment_service_payment_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_payment_service_payment_service_proto_goTypes = []interface{}{
	(Chain)(0),                    // 0: payment_service.Chain
	(Token)(0),                    // 1: payment_service.Token
	(PaymentStatus)(0),            // 2: payment_service.PaymentStatus
	(*CreatePaymentRequest)(nil),  // 3: payment_service.CreatePaymentRequest
	(*CreatePaymentResponse)(nil), // 4: payment_service.CreatePaymentResponse
	(*CheckPaymentRequest)(nil),   // 5: payment_service.CheckPaymentRequest
	(*CheckPaymentResponse)(nil),  // 6: payment_service.CheckPaymentResponse
}
var file_api_payment_service_payment_service_proto_depIdxs = []int32{
	0, // 0: payment_service.CreatePaymentRequest.chain:type_name -> payment_service.Chain
	1, // 1: payment_service.CreatePaymentRequest.token:type_name -> payment_service.Token
	2, // 2: payment_service.CheckPaymentResponse.status:type_name -> payment_service.PaymentStatus
	3, // 3: payment_service.PaymentService.CreatePayment:input_type -> payment_service.CreatePaymentRequest
	5, // 4: payment_service.PaymentService.CheckPayment:input_type -> payment_service.CheckPaymentRequest
	4, // 5: payment_service.PaymentService.CreatePayment:output_type -> payment_service.CreatePaymentResponse
	6, // 6: payment_service.PaymentService.CheckPayment:output_type -> payment_service.CheckPaymentResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_payment_service_payment_service_proto_init() }
func file_api_payment_service_payment_service_proto_init() {
	if File_api_payment_service_payment_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_payment_service_payment_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePaymentRequest); i {
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
		file_api_payment_service_payment_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePaymentResponse); i {
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
		file_api_payment_service_payment_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckPaymentRequest); i {
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
		file_api_payment_service_payment_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckPaymentResponse); i {
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
			RawDescriptor: file_api_payment_service_payment_service_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_payment_service_payment_service_proto_goTypes,
		DependencyIndexes: file_api_payment_service_payment_service_proto_depIdxs,
		EnumInfos:         file_api_payment_service_payment_service_proto_enumTypes,
		MessageInfos:      file_api_payment_service_payment_service_proto_msgTypes,
	}.Build()
	File_api_payment_service_payment_service_proto = out.File
	file_api_payment_service_payment_service_proto_rawDesc = nil
	file_api_payment_service_payment_service_proto_goTypes = nil
	file_api_payment_service_payment_service_proto_depIdxs = nil
}
