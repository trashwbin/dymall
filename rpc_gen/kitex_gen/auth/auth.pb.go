// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.29.3
// source: auth/auth.proto

package auth

import (
	context "context"
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

// 错误码枚举
type ErrorCode int32

const (
	ErrorCode_Success            ErrorCode = 0
	ErrorCode_GenerateTokenError ErrorCode = 1
	ErrorCode_TokenExpired       ErrorCode = 2
	ErrorCode_TokenInvalid       ErrorCode = 3
	ErrorCode_PermissionDenied   ErrorCode = 4
)

// Enum value maps for ErrorCode.
var (
	ErrorCode_name = map[int32]string{
		0: "Success",
		1: "GenerateTokenError",
		2: "TokenExpired",
		3: "TokenInvalid",
		4: "PermissionDenied",
	}
	ErrorCode_value = map[string]int32{
		"Success":            0,
		"GenerateTokenError": 1,
		"TokenExpired":       2,
		"TokenInvalid":       3,
		"PermissionDenied":   4,
	}
)

func (x ErrorCode) Enum() *ErrorCode {
	p := new(ErrorCode)
	*p = x
	return p
}

func (x ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_auth_auth_proto_enumTypes[0].Descriptor()
}

func (ErrorCode) Type() protoreflect.EnumType {
	return &file_auth_auth_proto_enumTypes[0]
}

func (x ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorCode.Descriptor instead.
func (ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_auth_auth_proto_rawDescGZIP(), []int{0}
}

// 令牌相关消息
type DeliverTokenReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // 只需要用户ID，角色信息从认证中心获取
}

func (x *DeliverTokenReq) Reset() {
	*x = DeliverTokenReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeliverTokenReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeliverTokenReq) ProtoMessage() {}

func (x *DeliverTokenReq) ProtoReflect() protoreflect.Message {
	mi := &file_auth_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeliverTokenReq.ProtoReflect.Descriptor instead.
func (*DeliverTokenReq) Descriptor() ([]byte, []int) {
	return file_auth_auth_proto_rawDescGZIP(), []int{0}
}

func (x *DeliverTokenReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type VerifyTokenReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *VerifyTokenReq) Reset() {
	*x = VerifyTokenReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_auth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyTokenReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyTokenReq) ProtoMessage() {}

func (x *VerifyTokenReq) ProtoReflect() protoreflect.Message {
	mi := &file_auth_auth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyTokenReq.ProtoReflect.Descriptor instead.
func (*VerifyTokenReq) Descriptor() ([]byte, []int) {
	return file_auth_auth_proto_rawDescGZIP(), []int{1}
}

func (x *VerifyTokenReq) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type DeliveryResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    ErrorCode `protobuf:"varint,1,opt,name=code,proto3,enum=auth.ErrorCode" json:"code,omitempty"`
	Message string    `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Token   string    `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	Role    string    `protobuf:"bytes,4,opt,name=role,proto3" json:"role,omitempty"` // 返回用户角色信息
}

func (x *DeliveryResp) Reset() {
	*x = DeliveryResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_auth_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeliveryResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeliveryResp) ProtoMessage() {}

func (x *DeliveryResp) ProtoReflect() protoreflect.Message {
	mi := &file_auth_auth_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeliveryResp.ProtoReflect.Descriptor instead.
func (*DeliveryResp) Descriptor() ([]byte, []int) {
	return file_auth_auth_proto_rawDescGZIP(), []int{2}
}

func (x *DeliveryResp) GetCode() ErrorCode {
	if x != nil {
		return x.Code
	}
	return ErrorCode_Success
}

func (x *DeliveryResp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *DeliveryResp) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *DeliveryResp) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

type VerifyResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    ErrorCode `protobuf:"varint,1,opt,name=code,proto3,enum=auth.ErrorCode" json:"code,omitempty"`
	Message string    `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	IsValid bool      `protobuf:"varint,3,opt,name=is_valid,json=isValid,proto3" json:"is_valid,omitempty"`
	UserId  int64     `protobuf:"varint,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Role    string    `protobuf:"bytes,5,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *VerifyResp) Reset() {
	*x = VerifyResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_auth_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyResp) ProtoMessage() {}

func (x *VerifyResp) ProtoReflect() protoreflect.Message {
	mi := &file_auth_auth_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyResp.ProtoReflect.Descriptor instead.
func (*VerifyResp) Descriptor() ([]byte, []int) {
	return file_auth_auth_proto_rawDescGZIP(), []int{3}
}

func (x *VerifyResp) GetCode() ErrorCode {
	if x != nil {
		return x.Code
	}
	return ErrorCode_Success
}

func (x *VerifyResp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *VerifyResp) GetIsValid() bool {
	if x != nil {
		return x.IsValid
	}
	return false
}

func (x *VerifyResp) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *VerifyResp) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

// 权限相关消息
type PolicyReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role     string `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`         // 角色
	Resource string `protobuf:"bytes,2,opt,name=resource,proto3" json:"resource,omitempty"` // 资源
	Action   string `protobuf:"bytes,3,opt,name=action,proto3" json:"action,omitempty"`     // 操作
}

func (x *PolicyReq) Reset() {
	*x = PolicyReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_auth_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PolicyReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PolicyReq) ProtoMessage() {}

func (x *PolicyReq) ProtoReflect() protoreflect.Message {
	mi := &file_auth_auth_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PolicyReq.ProtoReflect.Descriptor instead.
func (*PolicyReq) Descriptor() ([]byte, []int) {
	return file_auth_auth_proto_rawDescGZIP(), []int{4}
}

func (x *PolicyReq) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *PolicyReq) GetResource() string {
	if x != nil {
		return x.Resource
	}
	return ""
}

func (x *PolicyReq) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

type PolicyResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    ErrorCode `protobuf:"varint,1,opt,name=code,proto3,enum=auth.ErrorCode" json:"code,omitempty"`
	Message string    `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Success bool      `protobuf:"varint,3,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *PolicyResp) Reset() {
	*x = PolicyResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_auth_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PolicyResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PolicyResp) ProtoMessage() {}

func (x *PolicyResp) ProtoReflect() protoreflect.Message {
	mi := &file_auth_auth_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PolicyResp.ProtoReflect.Descriptor instead.
func (*PolicyResp) Descriptor() ([]byte, []int) {
	return file_auth_auth_proto_rawDescGZIP(), []int{5}
}

func (x *PolicyResp) GetCode() ErrorCode {
	if x != nil {
		return x.Code
	}
	return ErrorCode_Success
}

func (x *PolicyResp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *PolicyResp) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type RoleBindingReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Role   string `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *RoleBindingReq) Reset() {
	*x = RoleBindingReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_auth_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoleBindingReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoleBindingReq) ProtoMessage() {}

func (x *RoleBindingReq) ProtoReflect() protoreflect.Message {
	mi := &file_auth_auth_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoleBindingReq.ProtoReflect.Descriptor instead.
func (*RoleBindingReq) Descriptor() ([]byte, []int) {
	return file_auth_auth_proto_rawDescGZIP(), []int{6}
}

func (x *RoleBindingReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *RoleBindingReq) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

type RoleQueryReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *RoleQueryReq) Reset() {
	*x = RoleQueryReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_auth_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoleQueryReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoleQueryReq) ProtoMessage() {}

func (x *RoleQueryReq) ProtoReflect() protoreflect.Message {
	mi := &file_auth_auth_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoleQueryReq.ProtoReflect.Descriptor instead.
func (*RoleQueryReq) Descriptor() ([]byte, []int) {
	return file_auth_auth_proto_rawDescGZIP(), []int{7}
}

func (x *RoleQueryReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type RoleQueryResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    ErrorCode `protobuf:"varint,1,opt,name=code,proto3,enum=auth.ErrorCode" json:"code,omitempty"`
	Message string    `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Roles   []string  `protobuf:"bytes,3,rep,name=roles,proto3" json:"roles,omitempty"`
}

func (x *RoleQueryResp) Reset() {
	*x = RoleQueryResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_auth_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoleQueryResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoleQueryResp) ProtoMessage() {}

func (x *RoleQueryResp) ProtoReflect() protoreflect.Message {
	mi := &file_auth_auth_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoleQueryResp.ProtoReflect.Descriptor instead.
func (*RoleQueryResp) Descriptor() ([]byte, []int) {
	return file_auth_auth_proto_rawDescGZIP(), []int{8}
}

func (x *RoleQueryResp) GetCode() ErrorCode {
	if x != nil {
		return x.Code
	}
	return ErrorCode_Success
}

func (x *RoleQueryResp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *RoleQueryResp) GetRoles() []string {
	if x != nil {
		return x.Roles
	}
	return nil
}

var File_auth_auth_proto protoreflect.FileDescriptor

var file_auth_auth_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x61, 0x75, 0x74, 0x68, 0x22, 0x2a, 0x0a, 0x0f, 0x44, 0x65, 0x6c, 0x69, 0x76,
	0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x22, 0x26, 0x0a, 0x0e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x77, 0x0a, 0x0c, 0x44,
	0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x12, 0x23, 0x0a, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x72, 0x6f, 0x6c, 0x65, 0x22, 0x93, 0x01, 0x0a, 0x0a, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x12, 0x23, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x0f, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f,
	0x64, 0x65, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x5f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0x53, 0x0a, 0x09, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0x65, 0x0a, 0x0a, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x73, 0x70, 0x12, 0x23, 0x0a,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x61, 0x75,
	0x74, 0x68, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x3d, 0x0a, 0x0e, 0x52, 0x6f, 0x6c, 0x65, 0x42, 0x69,
	0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0x27, 0x0a, 0x0c, 0x52, 0x6f, 0x6c, 0x65, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x64,
	0x0a, 0x0d, 0x52, 0x6f, 0x6c, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x12,
	0x23, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e,
	0x61, 0x75, 0x74, 0x68, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x72,
	0x6f, 0x6c, 0x65, 0x73, 0x2a, 0x6a, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x10, 0x00, 0x12, 0x16,
	0x0a, 0x12, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x45,
	0x78, 0x70, 0x69, 0x72, 0x65, 0x64, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x10, 0x03, 0x12, 0x14, 0x0a, 0x10, 0x50, 0x65,
	0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x6e, 0x69, 0x65, 0x64, 0x10, 0x04,
	0x32, 0xad, 0x03, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x40, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x42, 0x79, 0x52, 0x50, 0x43, 0x12, 0x15, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x44, 0x65, 0x6c,
	0x69, 0x76, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x22, 0x00, 0x12, 0x3c, 0x0a, 0x10, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x42, 0x79, 0x52, 0x50, 0x43, 0x12, 0x14, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x10, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00,
	0x12, 0x30, 0x0a, 0x09, 0x41, 0x64, 0x64, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x0f, 0x2e,
	0x61, 0x75, 0x74, 0x68, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x71, 0x1a, 0x10,
	0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x22, 0x00, 0x12, 0x33, 0x0a, 0x0c, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x50, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x12, 0x0f, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x52, 0x65, 0x71, 0x1a, 0x10, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x0e, 0x41, 0x64, 0x64, 0x52, 0x6f,
	0x6c, 0x65, 0x46, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x12, 0x14, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x42, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x1a,
	0x10, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x11, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x52, 0x6f, 0x6c,
	0x65, 0x46, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x12, 0x14, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e,
	0x52, 0x6f, 0x6c, 0x65, 0x42, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x1a, 0x10,
	0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x22, 0x00, 0x12, 0x3c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x46, 0x6f,
	0x72, 0x55, 0x73, 0x65, 0x72, 0x12, 0x12, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x52, 0x6f, 0x6c,
	0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x1a, 0x13, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00,
	0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74,
	0x72, 0x61, 0x73, 0x68, 0x77, 0x62, 0x69, 0x6e, 0x2f, 0x64, 0x79, 0x6d, 0x61, 0x6c, 0x6c, 0x2f,
	0x72, 0x70, 0x63, 0x5f, 0x67, 0x65, 0x6e, 0x2f, 0x6b, 0x69, 0x74, 0x65, 0x78, 0x5f, 0x67, 0x65,
	0x6e, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_auth_proto_rawDescOnce sync.Once
	file_auth_auth_proto_rawDescData = file_auth_auth_proto_rawDesc
)

func file_auth_auth_proto_rawDescGZIP() []byte {
	file_auth_auth_proto_rawDescOnce.Do(func() {
		file_auth_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_auth_proto_rawDescData)
	})
	return file_auth_auth_proto_rawDescData
}

var file_auth_auth_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_auth_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_auth_auth_proto_goTypes = []interface{}{
	(ErrorCode)(0),          // 0: auth.ErrorCode
	(*DeliverTokenReq)(nil), // 1: auth.DeliverTokenReq
	(*VerifyTokenReq)(nil),  // 2: auth.VerifyTokenReq
	(*DeliveryResp)(nil),    // 3: auth.DeliveryResp
	(*VerifyResp)(nil),      // 4: auth.VerifyResp
	(*PolicyReq)(nil),       // 5: auth.PolicyReq
	(*PolicyResp)(nil),      // 6: auth.PolicyResp
	(*RoleBindingReq)(nil),  // 7: auth.RoleBindingReq
	(*RoleQueryReq)(nil),    // 8: auth.RoleQueryReq
	(*RoleQueryResp)(nil),   // 9: auth.RoleQueryResp
}
var file_auth_auth_proto_depIdxs = []int32{
	0,  // 0: auth.DeliveryResp.code:type_name -> auth.ErrorCode
	0,  // 1: auth.VerifyResp.code:type_name -> auth.ErrorCode
	0,  // 2: auth.PolicyResp.code:type_name -> auth.ErrorCode
	0,  // 3: auth.RoleQueryResp.code:type_name -> auth.ErrorCode
	1,  // 4: auth.AuthService.DeliverTokenByRPC:input_type -> auth.DeliverTokenReq
	2,  // 5: auth.AuthService.VerifyTokenByRPC:input_type -> auth.VerifyTokenReq
	5,  // 6: auth.AuthService.AddPolicy:input_type -> auth.PolicyReq
	5,  // 7: auth.AuthService.RemovePolicy:input_type -> auth.PolicyReq
	7,  // 8: auth.AuthService.AddRoleForUser:input_type -> auth.RoleBindingReq
	7,  // 9: auth.AuthService.RemoveRoleForUser:input_type -> auth.RoleBindingReq
	8,  // 10: auth.AuthService.GetRolesForUser:input_type -> auth.RoleQueryReq
	3,  // 11: auth.AuthService.DeliverTokenByRPC:output_type -> auth.DeliveryResp
	4,  // 12: auth.AuthService.VerifyTokenByRPC:output_type -> auth.VerifyResp
	6,  // 13: auth.AuthService.AddPolicy:output_type -> auth.PolicyResp
	6,  // 14: auth.AuthService.RemovePolicy:output_type -> auth.PolicyResp
	6,  // 15: auth.AuthService.AddRoleForUser:output_type -> auth.PolicyResp
	6,  // 16: auth.AuthService.RemoveRoleForUser:output_type -> auth.PolicyResp
	9,  // 17: auth.AuthService.GetRolesForUser:output_type -> auth.RoleQueryResp
	11, // [11:18] is the sub-list for method output_type
	4,  // [4:11] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_auth_auth_proto_init() }
func file_auth_auth_proto_init() {
	if File_auth_auth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_auth_auth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeliverTokenReq); i {
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
		file_auth_auth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyTokenReq); i {
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
		file_auth_auth_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeliveryResp); i {
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
		file_auth_auth_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyResp); i {
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
		file_auth_auth_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PolicyReq); i {
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
		file_auth_auth_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PolicyResp); i {
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
		file_auth_auth_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoleBindingReq); i {
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
		file_auth_auth_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoleQueryReq); i {
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
		file_auth_auth_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoleQueryResp); i {
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
			RawDescriptor: file_auth_auth_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_auth_auth_proto_goTypes,
		DependencyIndexes: file_auth_auth_proto_depIdxs,
		EnumInfos:         file_auth_auth_proto_enumTypes,
		MessageInfos:      file_auth_auth_proto_msgTypes,
	}.Build()
	File_auth_auth_proto = out.File
	file_auth_auth_proto_rawDesc = nil
	file_auth_auth_proto_goTypes = nil
	file_auth_auth_proto_depIdxs = nil
}

var _ context.Context

// Code generated by Kitex v0.9.1. DO NOT EDIT.

type AuthService interface {
	DeliverTokenByRPC(ctx context.Context, req *DeliverTokenReq) (res *DeliveryResp, err error)
	VerifyTokenByRPC(ctx context.Context, req *VerifyTokenReq) (res *VerifyResp, err error)
	AddPolicy(ctx context.Context, req *PolicyReq) (res *PolicyResp, err error)
	RemovePolicy(ctx context.Context, req *PolicyReq) (res *PolicyResp, err error)
	AddRoleForUser(ctx context.Context, req *RoleBindingReq) (res *PolicyResp, err error)
	RemoveRoleForUser(ctx context.Context, req *RoleBindingReq) (res *PolicyResp, err error)
	GetRolesForUser(ctx context.Context, req *RoleQueryReq) (res *RoleQueryResp, err error)
}
