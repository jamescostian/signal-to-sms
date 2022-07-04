//*
// Copyright (C) 2018 Open Whisper Systems
//
// Licensed according to the LICENSE file in this repository.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: backups.proto

package signal

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

type SqlStatement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Statement  *string                      `protobuf:"bytes,1,opt,name=statement" json:"statement,omitempty"`
	Parameters []*SqlStatement_SqlParameter `protobuf:"bytes,2,rep,name=parameters" json:"parameters,omitempty"`
}

func (x *SqlStatement) Reset() {
	*x = SqlStatement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backups_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SqlStatement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SqlStatement) ProtoMessage() {}

func (x *SqlStatement) ProtoReflect() protoreflect.Message {
	mi := &file_backups_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SqlStatement.ProtoReflect.Descriptor instead.
func (*SqlStatement) Descriptor() ([]byte, []int) {
	return file_backups_proto_rawDescGZIP(), []int{0}
}

func (x *SqlStatement) GetStatement() string {
	if x != nil && x.Statement != nil {
		return *x.Statement
	}
	return ""
}

func (x *SqlStatement) GetParameters() []*SqlStatement_SqlParameter {
	if x != nil {
		return x.Parameters
	}
	return nil
}

type SharedPreference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File             *string  `protobuf:"bytes,1,opt,name=file" json:"file,omitempty"`
	Key              *string  `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	Value            *string  `protobuf:"bytes,3,opt,name=value" json:"value,omitempty"`
	BooleanValue     *bool    `protobuf:"varint,4,opt,name=booleanValue" json:"booleanValue,omitempty"`
	StringSetValue   []string `protobuf:"bytes,5,rep,name=stringSetValue" json:"stringSetValue,omitempty"`
	IsStringSetValue *bool    `protobuf:"varint,6,opt,name=isStringSetValue" json:"isStringSetValue,omitempty"`
}

func (x *SharedPreference) Reset() {
	*x = SharedPreference{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backups_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SharedPreference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SharedPreference) ProtoMessage() {}

func (x *SharedPreference) ProtoReflect() protoreflect.Message {
	mi := &file_backups_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SharedPreference.ProtoReflect.Descriptor instead.
func (*SharedPreference) Descriptor() ([]byte, []int) {
	return file_backups_proto_rawDescGZIP(), []int{1}
}

func (x *SharedPreference) GetFile() string {
	if x != nil && x.File != nil {
		return *x.File
	}
	return ""
}

func (x *SharedPreference) GetKey() string {
	if x != nil && x.Key != nil {
		return *x.Key
	}
	return ""
}

func (x *SharedPreference) GetValue() string {
	if x != nil && x.Value != nil {
		return *x.Value
	}
	return ""
}

func (x *SharedPreference) GetBooleanValue() bool {
	if x != nil && x.BooleanValue != nil {
		return *x.BooleanValue
	}
	return false
}

func (x *SharedPreference) GetStringSetValue() []string {
	if x != nil {
		return x.StringSetValue
	}
	return nil
}

func (x *SharedPreference) GetIsStringSetValue() bool {
	if x != nil && x.IsStringSetValue != nil {
		return *x.IsStringSetValue
	}
	return false
}

type Attachment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RowId        *uint64 `protobuf:"varint,1,opt,name=rowId" json:"rowId,omitempty"`
	AttachmentId *uint64 `protobuf:"varint,2,opt,name=attachmentId" json:"attachmentId,omitempty"`
	Length       *uint32 `protobuf:"varint,3,opt,name=length" json:"length,omitempty"`
}

func (x *Attachment) Reset() {
	*x = Attachment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backups_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attachment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attachment) ProtoMessage() {}

func (x *Attachment) ProtoReflect() protoreflect.Message {
	mi := &file_backups_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attachment.ProtoReflect.Descriptor instead.
func (*Attachment) Descriptor() ([]byte, []int) {
	return file_backups_proto_rawDescGZIP(), []int{2}
}

func (x *Attachment) GetRowId() uint64 {
	if x != nil && x.RowId != nil {
		return *x.RowId
	}
	return 0
}

func (x *Attachment) GetAttachmentId() uint64 {
	if x != nil && x.AttachmentId != nil {
		return *x.AttachmentId
	}
	return 0
}

func (x *Attachment) GetLength() uint32 {
	if x != nil && x.Length != nil {
		return *x.Length
	}
	return 0
}

type Sticker struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RowId  *uint64 `protobuf:"varint,1,opt,name=rowId" json:"rowId,omitempty"`
	Length *uint32 `protobuf:"varint,2,opt,name=length" json:"length,omitempty"`
}

func (x *Sticker) Reset() {
	*x = Sticker{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backups_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sticker) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sticker) ProtoMessage() {}

func (x *Sticker) ProtoReflect() protoreflect.Message {
	mi := &file_backups_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sticker.ProtoReflect.Descriptor instead.
func (*Sticker) Descriptor() ([]byte, []int) {
	return file_backups_proto_rawDescGZIP(), []int{3}
}

func (x *Sticker) GetRowId() uint64 {
	if x != nil && x.RowId != nil {
		return *x.RowId
	}
	return 0
}

func (x *Sticker) GetLength() uint32 {
	if x != nil && x.Length != nil {
		return *x.Length
	}
	return 0
}

type Avatar struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        *string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	RecipientId *string `protobuf:"bytes,3,opt,name=recipientId" json:"recipientId,omitempty"`
	Length      *uint32 `protobuf:"varint,2,opt,name=length" json:"length,omitempty"`
}

func (x *Avatar) Reset() {
	*x = Avatar{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backups_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Avatar) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Avatar) ProtoMessage() {}

func (x *Avatar) ProtoReflect() protoreflect.Message {
	mi := &file_backups_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Avatar.ProtoReflect.Descriptor instead.
func (*Avatar) Descriptor() ([]byte, []int) {
	return file_backups_proto_rawDescGZIP(), []int{4}
}

func (x *Avatar) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *Avatar) GetRecipientId() string {
	if x != nil && x.RecipientId != nil {
		return *x.RecipientId
	}
	return ""
}

func (x *Avatar) GetLength() uint32 {
	if x != nil && x.Length != nil {
		return *x.Length
	}
	return 0
}

type DatabaseVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version *uint32 `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
}

func (x *DatabaseVersion) Reset() {
	*x = DatabaseVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backups_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatabaseVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatabaseVersion) ProtoMessage() {}

func (x *DatabaseVersion) ProtoReflect() protoreflect.Message {
	mi := &file_backups_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatabaseVersion.ProtoReflect.Descriptor instead.
func (*DatabaseVersion) Descriptor() ([]byte, []int) {
	return file_backups_proto_rawDescGZIP(), []int{5}
}

func (x *DatabaseVersion) GetVersion() uint32 {
	if x != nil && x.Version != nil {
		return *x.Version
	}
	return 0
}

type Header struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Iv   []byte `protobuf:"bytes,1,opt,name=iv" json:"iv,omitempty"`
	Salt []byte `protobuf:"bytes,2,opt,name=salt" json:"salt,omitempty"`
}

func (x *Header) Reset() {
	*x = Header{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backups_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Header) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Header) ProtoMessage() {}

func (x *Header) ProtoReflect() protoreflect.Message {
	mi := &file_backups_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Header.ProtoReflect.Descriptor instead.
func (*Header) Descriptor() ([]byte, []int) {
	return file_backups_proto_rawDescGZIP(), []int{6}
}

func (x *Header) GetIv() []byte {
	if x != nil {
		return x.Iv
	}
	return nil
}

func (x *Header) GetSalt() []byte {
	if x != nil {
		return x.Salt
	}
	return nil
}

type KeyValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key          *string  `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	BlobValue    []byte   `protobuf:"bytes,2,opt,name=blobValue" json:"blobValue,omitempty"`
	BooleanValue *bool    `protobuf:"varint,3,opt,name=booleanValue" json:"booleanValue,omitempty"`
	FloatValue   *float32 `protobuf:"fixed32,4,opt,name=floatValue" json:"floatValue,omitempty"`
	IntegerValue *int32   `protobuf:"varint,5,opt,name=integerValue" json:"integerValue,omitempty"`
	LongValue    *int64   `protobuf:"varint,6,opt,name=longValue" json:"longValue,omitempty"`
	StringValue  *string  `protobuf:"bytes,7,opt,name=stringValue" json:"stringValue,omitempty"`
}

func (x *KeyValue) Reset() {
	*x = KeyValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backups_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValue) ProtoMessage() {}

func (x *KeyValue) ProtoReflect() protoreflect.Message {
	mi := &file_backups_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValue.ProtoReflect.Descriptor instead.
func (*KeyValue) Descriptor() ([]byte, []int) {
	return file_backups_proto_rawDescGZIP(), []int{7}
}

func (x *KeyValue) GetKey() string {
	if x != nil && x.Key != nil {
		return *x.Key
	}
	return ""
}

func (x *KeyValue) GetBlobValue() []byte {
	if x != nil {
		return x.BlobValue
	}
	return nil
}

func (x *KeyValue) GetBooleanValue() bool {
	if x != nil && x.BooleanValue != nil {
		return *x.BooleanValue
	}
	return false
}

func (x *KeyValue) GetFloatValue() float32 {
	if x != nil && x.FloatValue != nil {
		return *x.FloatValue
	}
	return 0
}

func (x *KeyValue) GetIntegerValue() int32 {
	if x != nil && x.IntegerValue != nil {
		return *x.IntegerValue
	}
	return 0
}

func (x *KeyValue) GetLongValue() int64 {
	if x != nil && x.LongValue != nil {
		return *x.LongValue
	}
	return 0
}

func (x *KeyValue) GetStringValue() string {
	if x != nil && x.StringValue != nil {
		return *x.StringValue
	}
	return ""
}

type BackupFrame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header     *Header           `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Statement  *SqlStatement     `protobuf:"bytes,2,opt,name=statement" json:"statement,omitempty"`
	Preference *SharedPreference `protobuf:"bytes,3,opt,name=preference" json:"preference,omitempty"`
	Attachment *Attachment       `protobuf:"bytes,4,opt,name=attachment" json:"attachment,omitempty"`
	Version    *DatabaseVersion  `protobuf:"bytes,5,opt,name=version" json:"version,omitempty"`
	End        *bool             `protobuf:"varint,6,opt,name=end" json:"end,omitempty"`
	Avatar     *Avatar           `protobuf:"bytes,7,opt,name=avatar" json:"avatar,omitempty"`
	Sticker    *Sticker          `protobuf:"bytes,8,opt,name=sticker" json:"sticker,omitempty"`
	KeyValue   *KeyValue         `protobuf:"bytes,9,opt,name=keyValue" json:"keyValue,omitempty"`
}

func (x *BackupFrame) Reset() {
	*x = BackupFrame{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backups_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BackupFrame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BackupFrame) ProtoMessage() {}

func (x *BackupFrame) ProtoReflect() protoreflect.Message {
	mi := &file_backups_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BackupFrame.ProtoReflect.Descriptor instead.
func (*BackupFrame) Descriptor() ([]byte, []int) {
	return file_backups_proto_rawDescGZIP(), []int{8}
}

func (x *BackupFrame) GetHeader() *Header {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *BackupFrame) GetStatement() *SqlStatement {
	if x != nil {
		return x.Statement
	}
	return nil
}

func (x *BackupFrame) GetPreference() *SharedPreference {
	if x != nil {
		return x.Preference
	}
	return nil
}

func (x *BackupFrame) GetAttachment() *Attachment {
	if x != nil {
		return x.Attachment
	}
	return nil
}

func (x *BackupFrame) GetVersion() *DatabaseVersion {
	if x != nil {
		return x.Version
	}
	return nil
}

func (x *BackupFrame) GetEnd() bool {
	if x != nil && x.End != nil {
		return *x.End
	}
	return false
}

func (x *BackupFrame) GetAvatar() *Avatar {
	if x != nil {
		return x.Avatar
	}
	return nil
}

func (x *BackupFrame) GetSticker() *Sticker {
	if x != nil {
		return x.Sticker
	}
	return nil
}

func (x *BackupFrame) GetKeyValue() *KeyValue {
	if x != nil {
		return x.KeyValue
	}
	return nil
}

type SqlStatement_SqlParameter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StringParamter   *string  `protobuf:"bytes,1,opt,name=stringParamter" json:"stringParamter,omitempty"`
	IntegerParameter *uint64  `protobuf:"varint,2,opt,name=integerParameter" json:"integerParameter,omitempty"`
	DoubleParameter  *float64 `protobuf:"fixed64,3,opt,name=doubleParameter" json:"doubleParameter,omitempty"`
	BlobParameter    []byte   `protobuf:"bytes,4,opt,name=blobParameter" json:"blobParameter,omitempty"`
	Nullparameter    *bool    `protobuf:"varint,5,opt,name=nullparameter" json:"nullparameter,omitempty"`
}

func (x *SqlStatement_SqlParameter) Reset() {
	*x = SqlStatement_SqlParameter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backups_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SqlStatement_SqlParameter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SqlStatement_SqlParameter) ProtoMessage() {}

func (x *SqlStatement_SqlParameter) ProtoReflect() protoreflect.Message {
	mi := &file_backups_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SqlStatement_SqlParameter.ProtoReflect.Descriptor instead.
func (*SqlStatement_SqlParameter) Descriptor() ([]byte, []int) {
	return file_backups_proto_rawDescGZIP(), []int{0, 0}
}

func (x *SqlStatement_SqlParameter) GetStringParamter() string {
	if x != nil && x.StringParamter != nil {
		return *x.StringParamter
	}
	return ""
}

func (x *SqlStatement_SqlParameter) GetIntegerParameter() uint64 {
	if x != nil && x.IntegerParameter != nil {
		return *x.IntegerParameter
	}
	return 0
}

func (x *SqlStatement_SqlParameter) GetDoubleParameter() float64 {
	if x != nil && x.DoubleParameter != nil {
		return *x.DoubleParameter
	}
	return 0
}

func (x *SqlStatement_SqlParameter) GetBlobParameter() []byte {
	if x != nil {
		return x.BlobParameter
	}
	return nil
}

func (x *SqlStatement_SqlParameter) GetNullparameter() bool {
	if x != nil && x.Nullparameter != nil {
		return *x.Nullparameter
	}
	return false
}

var File_backups_proto protoreflect.FileDescriptor

var file_backups_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x62, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x06, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x22, 0xca, 0x02, 0x0a, 0x0c, 0x53, 0x71, 0x6c, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x41, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65,
	0x74, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x6c, 0x2e, 0x53, 0x71, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x53, 0x71, 0x6c, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x52, 0x0a, 0x70,
	0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x1a, 0xd8, 0x01, 0x0a, 0x0c, 0x53, 0x71,
	0x6c, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x12, 0x26, 0x0a, 0x0e, 0x73, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x74,
	0x65, 0x72, 0x12, 0x2a, 0x0a, 0x10, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x65, 0x72, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x69, 0x6e,
	0x74, 0x65, 0x67, 0x65, 0x72, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x12, 0x28,
	0x0a, 0x0f, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x0d, 0x62, 0x6c, 0x6f, 0x62,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0d, 0x62, 0x6c, 0x6f, 0x62, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x12, 0x24,
	0x0a, 0x0d, 0x6e, 0x75, 0x6c, 0x6c, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x6e, 0x75, 0x6c, 0x6c, 0x70, 0x61, 0x72, 0x61, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x22, 0xc6, 0x01, 0x0a, 0x10, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x50,
	0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x69, 0x6c,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x62, 0x6f, 0x6f,
	0x6c, 0x65, 0x61, 0x6e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x73, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x53, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x0e, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x2a, 0x0a, 0x10, 0x69, 0x73, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x74,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x10, 0x69, 0x73, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x5e, 0x0a,
	0x0a, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x72,
	0x6f, 0x77, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x77, 0x49,
	0x64, 0x12, 0x22, 0x0a, 0x0c, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x22, 0x37, 0x0a,
	0x07, 0x53, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x77, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x77, 0x49, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06,
	0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x22, 0x56, 0x0a, 0x06, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x63, 0x69, 0x70,
	0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x22, 0x2b,
	0x0a, 0x0f, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x2c, 0x0a, 0x06, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x76, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x02, 0x69, 0x76, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x61, 0x6c, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x73, 0x61, 0x6c, 0x74, 0x22, 0xe2, 0x01, 0x0a, 0x08, 0x4b, 0x65,
	0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x62,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x62, 0x6c, 0x6f,
	0x62, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61,
	0x6e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x62, 0x6f,
	0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x66, 0x6c,
	0x6f, 0x61, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0a,
	0x66, 0x6c, 0x6f, 0x61, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x69, 0x6e,
	0x74, 0x65, 0x67, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0c, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x9d,
	0x03, 0x0a, 0x0b, 0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x26,
	0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e,
	0x2e, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x32, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x74, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x6c, 0x2e, 0x53, 0x71, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x09, 0x73, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x38, 0x0a, 0x0a, 0x70, 0x72,
	0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x50, 0x72,
	0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0a, 0x70, 0x72, 0x65, 0x66, 0x65, 0x72,
	0x65, 0x6e, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65,
	0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x6c, 0x2e, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x61, 0x74,
	0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x31, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x6c, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x65,
	0x6e, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x12, 0x26, 0x0a,
	0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x2e, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x52, 0x06, 0x61,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x12, 0x29, 0x0a, 0x07, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x2e,
	0x53, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x52, 0x07, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72,
	0x12, 0x2c, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x2e, 0x4b, 0x65, 0x79, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x31,
	0x0a, 0x21, 0x6f, 0x72, 0x67, 0x2e, 0x74, 0x68, 0x6f, 0x75, 0x67, 0x68, 0x74, 0x63, 0x72, 0x69,
	0x6d, 0x65, 0x2e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x65, 0x73, 0x6d, 0x73, 0x2e, 0x62, 0x61, 0x63,
	0x6b, 0x75, 0x70, 0x42, 0x0c, 0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x73,
}

var (
	file_backups_proto_rawDescOnce sync.Once
	file_backups_proto_rawDescData = file_backups_proto_rawDesc
)

func file_backups_proto_rawDescGZIP() []byte {
	file_backups_proto_rawDescOnce.Do(func() {
		file_backups_proto_rawDescData = protoimpl.X.CompressGZIP(file_backups_proto_rawDescData)
	})
	return file_backups_proto_rawDescData
}

var file_backups_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_backups_proto_goTypes = []interface{}{
	(*SqlStatement)(nil),              // 0: signal.SqlStatement
	(*SharedPreference)(nil),          // 1: signal.SharedPreference
	(*Attachment)(nil),                // 2: signal.Attachment
	(*Sticker)(nil),                   // 3: signal.Sticker
	(*Avatar)(nil),                    // 4: signal.Avatar
	(*DatabaseVersion)(nil),           // 5: signal.DatabaseVersion
	(*Header)(nil),                    // 6: signal.Header
	(*KeyValue)(nil),                  // 7: signal.KeyValue
	(*BackupFrame)(nil),               // 8: signal.BackupFrame
	(*SqlStatement_SqlParameter)(nil), // 9: signal.SqlStatement.SqlParameter
}
var file_backups_proto_depIdxs = []int32{
	9, // 0: signal.SqlStatement.parameters:type_name -> signal.SqlStatement.SqlParameter
	6, // 1: signal.BackupFrame.header:type_name -> signal.Header
	0, // 2: signal.BackupFrame.statement:type_name -> signal.SqlStatement
	1, // 3: signal.BackupFrame.preference:type_name -> signal.SharedPreference
	2, // 4: signal.BackupFrame.attachment:type_name -> signal.Attachment
	5, // 5: signal.BackupFrame.version:type_name -> signal.DatabaseVersion
	4, // 6: signal.BackupFrame.avatar:type_name -> signal.Avatar
	3, // 7: signal.BackupFrame.sticker:type_name -> signal.Sticker
	7, // 8: signal.BackupFrame.keyValue:type_name -> signal.KeyValue
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_backups_proto_init() }
func file_backups_proto_init() {
	if File_backups_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_backups_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SqlStatement); i {
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
		file_backups_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SharedPreference); i {
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
		file_backups_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Attachment); i {
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
		file_backups_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sticker); i {
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
		file_backups_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Avatar); i {
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
		file_backups_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatabaseVersion); i {
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
		file_backups_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Header); i {
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
		file_backups_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValue); i {
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
		file_backups_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BackupFrame); i {
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
		file_backups_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SqlStatement_SqlParameter); i {
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
			RawDescriptor: file_backups_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_backups_proto_goTypes,
		DependencyIndexes: file_backups_proto_depIdxs,
		MessageInfos:      file_backups_proto_msgTypes,
	}.Build()
	File_backups_proto = out.File
	file_backups_proto_rawDesc = nil
	file_backups_proto_goTypes = nil
	file_backups_proto_depIdxs = nil
}
