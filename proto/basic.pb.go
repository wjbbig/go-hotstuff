// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: basic.proto

package proto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Msg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Payload:
	//	*Msg_Prepare
	//	*Msg_PrepareVote
	//	*Msg_PreCommit
	//	*Msg_PreCommitVote
	//	*Msg_Commit
	//	*Msg_CommitVote
	//	*Msg_Decide
	//	*Msg_NewView
	//	*Msg_Request
	Payload isMsg_Payload `protobuf_oneof:"Payload"`
}

func (x *Msg) Reset() {
	*x = Msg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Msg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Msg) ProtoMessage() {}

func (x *Msg) ProtoReflect() protoreflect.Message {
	mi := &file_basic_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Msg.ProtoReflect.Descriptor instead.
func (*Msg) Descriptor() ([]byte, []int) {
	return file_basic_proto_rawDescGZIP(), []int{0}
}

func (m *Msg) GetPayload() isMsg_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (x *Msg) GetPrepare() *Prepare {
	if x, ok := x.GetPayload().(*Msg_Prepare); ok {
		return x.Prepare
	}
	return nil
}

func (x *Msg) GetPrepareVote() *PrepareVote {
	if x, ok := x.GetPayload().(*Msg_PrepareVote); ok {
		return x.PrepareVote
	}
	return nil
}

func (x *Msg) GetPreCommit() *PreCommit {
	if x, ok := x.GetPayload().(*Msg_PreCommit); ok {
		return x.PreCommit
	}
	return nil
}

func (x *Msg) GetPreCommitVote() *PreCommitVote {
	if x, ok := x.GetPayload().(*Msg_PreCommitVote); ok {
		return x.PreCommitVote
	}
	return nil
}

func (x *Msg) GetCommit() *Commit {
	if x, ok := x.GetPayload().(*Msg_Commit); ok {
		return x.Commit
	}
	return nil
}

func (x *Msg) GetCommitVote() *CommitVote {
	if x, ok := x.GetPayload().(*Msg_CommitVote); ok {
		return x.CommitVote
	}
	return nil
}

func (x *Msg) GetDecide() *Decide {
	if x, ok := x.GetPayload().(*Msg_Decide); ok {
		return x.Decide
	}
	return nil
}

func (x *Msg) GetNewView() *NewView {
	if x, ok := x.GetPayload().(*Msg_NewView); ok {
		return x.NewView
	}
	return nil
}

func (x *Msg) GetRequest() *Request {
	if x, ok := x.GetPayload().(*Msg_Request); ok {
		return x.Request
	}
	return nil
}

type isMsg_Payload interface {
	isMsg_Payload()
}

type Msg_Prepare struct {
	Prepare *Prepare `protobuf:"bytes,1,opt,name=prepare,proto3,oneof"`
}

type Msg_PrepareVote struct {
	PrepareVote *PrepareVote `protobuf:"bytes,2,opt,name=prepareVote,proto3,oneof"`
}

type Msg_PreCommit struct {
	PreCommit *PreCommit `protobuf:"bytes,3,opt,name=preCommit,proto3,oneof"`
}

type Msg_PreCommitVote struct {
	PreCommitVote *PreCommitVote `protobuf:"bytes,4,opt,name=preCommitVote,proto3,oneof"`
}

type Msg_Commit struct {
	Commit *Commit `protobuf:"bytes,5,opt,name=commit,proto3,oneof"`
}

type Msg_CommitVote struct {
	CommitVote *CommitVote `protobuf:"bytes,6,opt,name=commitVote,proto3,oneof"`
}

type Msg_Decide struct {
	Decide *Decide `protobuf:"bytes,7,opt,name=decide,proto3,oneof"`
}

type Msg_NewView struct {
	NewView *NewView `protobuf:"bytes,8,opt,name=newView,proto3,oneof"`
}

type Msg_Request struct {
	Request *Request `protobuf:"bytes,9,opt,name=request,proto3,oneof"`
}

func (*Msg_Prepare) isMsg_Payload() {}

func (*Msg_PrepareVote) isMsg_Payload() {}

func (*Msg_PreCommit) isMsg_Payload() {}

func (*Msg_PreCommitVote) isMsg_Payload() {}

func (*Msg_Commit) isMsg_Payload() {}

func (*Msg_CommitVote) isMsg_Payload() {}

func (*Msg_Decide) isMsg_Payload() {}

func (*Msg_NewView) isMsg_Payload() {}

func (*Msg_Request) isMsg_Payload() {}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_basic_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_basic_proto_rawDescGZIP(), []int{1}
}

type Prepare struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CurProposal *Block      `protobuf:"bytes,1,opt,name=curProposal,proto3" json:"curProposal,omitempty"`
	HighQC      *QuorumCert `protobuf:"bytes,2,opt,name=highQC,proto3" json:"highQC,omitempty"`
}

func (x *Prepare) Reset() {
	*x = Prepare{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Prepare) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Prepare) ProtoMessage() {}

func (x *Prepare) ProtoReflect() protoreflect.Message {
	mi := &file_basic_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Prepare.ProtoReflect.Descriptor instead.
func (*Prepare) Descriptor() ([]byte, []int) {
	return file_basic_proto_rawDescGZIP(), []int{2}
}

func (x *Prepare) GetCurProposal() *Block {
	if x != nil {
		return x.CurProposal
	}
	return nil
}

func (x *Prepare) GetHighQC() *QuorumCert {
	if x != nil {
		return x.HighQC
	}
	return nil
}

type PrepareVote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockHash  []byte      `protobuf:"bytes,1,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	Qc         *QuorumCert `protobuf:"bytes,2,opt,name=qc,proto3" json:"qc,omitempty"`
	PartialSig []byte      `protobuf:"bytes,3,opt,name=partialSig,proto3" json:"partialSig,omitempty"`
}

func (x *PrepareVote) Reset() {
	*x = PrepareVote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareVote) ProtoMessage() {}

func (x *PrepareVote) ProtoReflect() protoreflect.Message {
	mi := &file_basic_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareVote.ProtoReflect.Descriptor instead.
func (*PrepareVote) Descriptor() ([]byte, []int) {
	return file_basic_proto_rawDescGZIP(), []int{3}
}

func (x *PrepareVote) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *PrepareVote) GetQc() *QuorumCert {
	if x != nil {
		return x.Qc
	}
	return nil
}

func (x *PrepareVote) GetPartialSig() []byte {
	if x != nil {
		return x.PartialSig
	}
	return nil
}

type PreCommit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PrepareQC *QuorumCert `protobuf:"bytes,1,opt,name=prepareQC,proto3" json:"prepareQC,omitempty"`
}

func (x *PreCommit) Reset() {
	*x = PreCommit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreCommit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreCommit) ProtoMessage() {}

func (x *PreCommit) ProtoReflect() protoreflect.Message {
	mi := &file_basic_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreCommit.ProtoReflect.Descriptor instead.
func (*PreCommit) Descriptor() ([]byte, []int) {
	return file_basic_proto_rawDescGZIP(), []int{4}
}

func (x *PreCommit) GetPrepareQC() *QuorumCert {
	if x != nil {
		return x.PrepareQC
	}
	return nil
}

type PreCommitVote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockHash  []byte      `protobuf:"bytes,1,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	Qc         *QuorumCert `protobuf:"bytes,2,opt,name=qc,proto3" json:"qc,omitempty"`
	PartialSig []byte      `protobuf:"bytes,3,opt,name=partialSig,proto3" json:"partialSig,omitempty"`
}

func (x *PreCommitVote) Reset() {
	*x = PreCommitVote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreCommitVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreCommitVote) ProtoMessage() {}

func (x *PreCommitVote) ProtoReflect() protoreflect.Message {
	mi := &file_basic_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreCommitVote.ProtoReflect.Descriptor instead.
func (*PreCommitVote) Descriptor() ([]byte, []int) {
	return file_basic_proto_rawDescGZIP(), []int{5}
}

func (x *PreCommitVote) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *PreCommitVote) GetQc() *QuorumCert {
	if x != nil {
		return x.Qc
	}
	return nil
}

func (x *PreCommitVote) GetPartialSig() []byte {
	if x != nil {
		return x.PartialSig
	}
	return nil
}

type Commit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PreCommitQC *QuorumCert `protobuf:"bytes,1,opt,name=preCommitQC,proto3" json:"preCommitQC,omitempty"`
}

func (x *Commit) Reset() {
	*x = Commit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Commit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Commit) ProtoMessage() {}

func (x *Commit) ProtoReflect() protoreflect.Message {
	mi := &file_basic_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Commit.ProtoReflect.Descriptor instead.
func (*Commit) Descriptor() ([]byte, []int) {
	return file_basic_proto_rawDescGZIP(), []int{6}
}

func (x *Commit) GetPreCommitQC() *QuorumCert {
	if x != nil {
		return x.PreCommitQC
	}
	return nil
}

type CommitVote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockHash  []byte      `protobuf:"bytes,1,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	Qc         *QuorumCert `protobuf:"bytes,2,opt,name=qc,proto3" json:"qc,omitempty"`
	PartialSig []byte      `protobuf:"bytes,3,opt,name=partialSig,proto3" json:"partialSig,omitempty"`
}

func (x *CommitVote) Reset() {
	*x = CommitVote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommitVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitVote) ProtoMessage() {}

func (x *CommitVote) ProtoReflect() protoreflect.Message {
	mi := &file_basic_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitVote.ProtoReflect.Descriptor instead.
func (*CommitVote) Descriptor() ([]byte, []int) {
	return file_basic_proto_rawDescGZIP(), []int{7}
}

func (x *CommitVote) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *CommitVote) GetQc() *QuorumCert {
	if x != nil {
		return x.Qc
	}
	return nil
}

func (x *CommitVote) GetPartialSig() []byte {
	if x != nil {
		return x.PartialSig
	}
	return nil
}

type Decide struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommitQC *QuorumCert `protobuf:"bytes,1,opt,name=commitQC,proto3" json:"commitQC,omitempty"`
}

func (x *Decide) Reset() {
	*x = Decide{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Decide) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Decide) ProtoMessage() {}

func (x *Decide) ProtoReflect() protoreflect.Message {
	mi := &file_basic_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Decide.ProtoReflect.Descriptor instead.
func (*Decide) Descriptor() ([]byte, []int) {
	return file_basic_proto_rawDescGZIP(), []int{8}
}

func (x *Decide) GetCommitQC() *QuorumCert {
	if x != nil {
		return x.CommitQC
	}
	return nil
}

type NewView struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PrepareQC *QuorumCert `protobuf:"bytes,1,opt,name=prepareQC,proto3" json:"prepareQC,omitempty"`
}

func (x *NewView) Reset() {
	*x = NewView{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewView) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewView) ProtoMessage() {}

func (x *NewView) ProtoReflect() protoreflect.Message {
	mi := &file_basic_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewView.ProtoReflect.Descriptor instead.
func (*NewView) Descriptor() ([]byte, []int) {
	return file_basic_proto_rawDescGZIP(), []int{9}
}

func (x *NewView) GetPrepareQC() *QuorumCert {
	if x != nil {
		return x.PrepareQC
	}
	return nil
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd           string `protobuf:"bytes,1,opt,name=cmd,proto3" json:"cmd,omitempty"`
	ClientAddress string `protobuf:"bytes,2,opt,name=clientAddress,proto3" json:"clientAddress,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_basic_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_basic_proto_rawDescGZIP(), []int{10}
}

func (x *Request) GetCmd() string {
	if x != nil {
		return x.Cmd
	}
	return ""
}

func (x *Request) GetClientAddress() string {
	if x != nil {
		return x.ClientAddress
	}
	return ""
}

type Reply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *Reply) Reset() {
	*x = Reply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_basic_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Reply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reply) ProtoMessage() {}

func (x *Reply) ProtoReflect() protoreflect.Message {
	mi := &file_basic_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reply.ProtoReflect.Descriptor instead.
func (*Reply) Descriptor() ([]byte, []int) {
	return file_basic_proto_rawDescGZIP(), []int{11}
}

func (x *Reply) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_basic_proto protoreflect.FileDescriptor

var file_basic_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xc3, 0x03, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x12, 0x2a, 0x0a, 0x07, 0x70, 0x72,
	0x65, 0x70, 0x61, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x48, 0x00, 0x52, 0x07, 0x70,
	0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12, 0x36, 0x0a, 0x0b, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72,
	0x65, 0x56, 0x6f, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x56, 0x6f, 0x74, 0x65, 0x48,
	0x00, 0x52, 0x0b, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x30,
	0x0a, 0x09, 0x70, 0x72, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x65, 0x43, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x48, 0x00, 0x52, 0x09, 0x70, 0x72, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x12, 0x3c, 0x0a, 0x0d, 0x70, 0x72, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x56, 0x6f, 0x74,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x50, 0x72, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x48, 0x00, 0x52,
	0x0d, 0x70, 0x72, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x27,
	0x0a, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x48, 0x00, 0x52,
	0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x33, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x56, 0x6f, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x48, 0x00,
	0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x27, 0x0a, 0x06,
	0x64, 0x65, 0x63, 0x69, 0x64, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x63, 0x69, 0x64, 0x65, 0x48, 0x00, 0x52, 0x06, 0x64,
	0x65, 0x63, 0x69, 0x64, 0x65, 0x12, 0x2a, 0x0a, 0x07, 0x6e, 0x65, 0x77, 0x56, 0x69, 0x65, 0x77,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e,
	0x65, 0x77, 0x56, 0x69, 0x65, 0x77, 0x48, 0x00, 0x52, 0x07, 0x6e, 0x65, 0x77, 0x56, 0x69, 0x65,
	0x77, 0x12, 0x2a, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x48, 0x00, 0x52, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x09, 0x0a,
	0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x64, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12, 0x2e, 0x0a, 0x0b,
	0x63, 0x75, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52,
	0x0b, 0x63, 0x75, 0x72, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x61, 0x6c, 0x12, 0x29, 0x0a, 0x06,
	0x68, 0x69, 0x67, 0x68, 0x51, 0x43, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72, 0x74, 0x52,
	0x06, 0x68, 0x69, 0x67, 0x68, 0x51, 0x43, 0x22, 0x6e, 0x0a, 0x0b, 0x50, 0x72, 0x65, 0x70, 0x61,
	0x72, 0x65, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48,
	0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x21, 0x0a, 0x02, 0x71, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43,
	0x65, 0x72, 0x74, 0x52, 0x02, 0x71, 0x63, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x74, 0x69,
	0x61, 0x6c, 0x53, 0x69, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x70, 0x61, 0x72,
	0x74, 0x69, 0x61, 0x6c, 0x53, 0x69, 0x67, 0x22, 0x3c, 0x0a, 0x09, 0x50, 0x72, 0x65, 0x43, 0x6f,
	0x6d, 0x6d, 0x69, 0x74, 0x12, 0x2f, 0x0a, 0x09, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x51,
	0x43, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72, 0x74, 0x52, 0x09, 0x70, 0x72, 0x65, 0x70,
	0x61, 0x72, 0x65, 0x51, 0x43, 0x22, 0x70, 0x0a, 0x0d, 0x50, 0x72, 0x65, 0x43, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48,
	0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x21, 0x0a, 0x02, 0x71, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43,
	0x65, 0x72, 0x74, 0x52, 0x02, 0x71, 0x63, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x74, 0x69,
	0x61, 0x6c, 0x53, 0x69, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x70, 0x61, 0x72,
	0x74, 0x69, 0x61, 0x6c, 0x53, 0x69, 0x67, 0x22, 0x3d, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x12, 0x33, 0x0a, 0x0b, 0x70, 0x72, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x51, 0x43,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51,
	0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72, 0x74, 0x52, 0x0b, 0x70, 0x72, 0x65, 0x43, 0x6f,
	0x6d, 0x6d, 0x69, 0x74, 0x51, 0x43, 0x22, 0x6d, 0x0a, 0x0a, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x56, 0x6f, 0x74, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73,
	0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61,
	0x73, 0x68, 0x12, 0x21, 0x0a, 0x02, 0x71, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72,
	0x74, 0x52, 0x02, 0x71, 0x63, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c,
	0x53, 0x69, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x74, 0x69,
	0x61, 0x6c, 0x53, 0x69, 0x67, 0x22, 0x37, 0x0a, 0x06, 0x44, 0x65, 0x63, 0x69, 0x64, 0x65, 0x12,
	0x2d, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x51, 0x43, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d,
	0x43, 0x65, 0x72, 0x74, 0x52, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x51, 0x43, 0x22, 0x3a,
	0x0a, 0x07, 0x4e, 0x65, 0x77, 0x56, 0x69, 0x65, 0x77, 0x12, 0x2f, 0x0a, 0x09, 0x70, 0x72, 0x65,
	0x70, 0x61, 0x72, 0x65, 0x51, 0x43, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72, 0x74, 0x52,
	0x09, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x51, 0x43, 0x22, 0x41, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x1f, 0x0a,
	0x05, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0x8a,
	0x01, 0x0a, 0x0d, 0x42, 0x61, 0x73, 0x69, 0x63, 0x48, 0x6f, 0x74, 0x53, 0x74, 0x75, 0x66, 0x66,
	0x12, 0x25, 0x0a, 0x07, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x12, 0x0a, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x73, 0x67, 0x1a, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x29, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d,
	0x73, 0x67, 0x1a, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x27, 0x0a, 0x09, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x73, 0x67, 0x1a, 0x0c, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_basic_proto_rawDescOnce sync.Once
	file_basic_proto_rawDescData = file_basic_proto_rawDesc
)

func file_basic_proto_rawDescGZIP() []byte {
	file_basic_proto_rawDescOnce.Do(func() {
		file_basic_proto_rawDescData = protoimpl.X.CompressGZIP(file_basic_proto_rawDescData)
	})
	return file_basic_proto_rawDescData
}

var file_basic_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_basic_proto_goTypes = []interface{}{
	(*Msg)(nil),           // 0: proto.Msg
	(*Empty)(nil),         // 1: proto.Empty
	(*Prepare)(nil),       // 2: proto.Prepare
	(*PrepareVote)(nil),   // 3: proto.PrepareVote
	(*PreCommit)(nil),     // 4: proto.PreCommit
	(*PreCommitVote)(nil), // 5: proto.PreCommitVote
	(*Commit)(nil),        // 6: proto.Commit
	(*CommitVote)(nil),    // 7: proto.CommitVote
	(*Decide)(nil),        // 8: proto.Decide
	(*NewView)(nil),       // 9: proto.NewView
	(*Request)(nil),       // 10: proto.Request
	(*Reply)(nil),         // 11: proto.Reply
	(*Block)(nil),         // 12: proto.Block
	(*QuorumCert)(nil),    // 13: proto.QuorumCert
}
var file_basic_proto_depIdxs = []int32{
	2,  // 0: proto.Msg.prepare:type_name -> proto.Prepare
	3,  // 1: proto.Msg.prepareVote:type_name -> proto.PrepareVote
	4,  // 2: proto.Msg.preCommit:type_name -> proto.PreCommit
	5,  // 3: proto.Msg.preCommitVote:type_name -> proto.PreCommitVote
	6,  // 4: proto.Msg.commit:type_name -> proto.Commit
	7,  // 5: proto.Msg.commitVote:type_name -> proto.CommitVote
	8,  // 6: proto.Msg.decide:type_name -> proto.Decide
	9,  // 7: proto.Msg.newView:type_name -> proto.NewView
	10, // 8: proto.Msg.request:type_name -> proto.Request
	12, // 9: proto.Prepare.curProposal:type_name -> proto.Block
	13, // 10: proto.Prepare.highQC:type_name -> proto.QuorumCert
	13, // 11: proto.PrepareVote.qc:type_name -> proto.QuorumCert
	13, // 12: proto.PreCommit.prepareQC:type_name -> proto.QuorumCert
	13, // 13: proto.PreCommitVote.qc:type_name -> proto.QuorumCert
	13, // 14: proto.Commit.preCommitQC:type_name -> proto.QuorumCert
	13, // 15: proto.CommitVote.qc:type_name -> proto.QuorumCert
	13, // 16: proto.Decide.commitQC:type_name -> proto.QuorumCert
	13, // 17: proto.NewView.prepareQC:type_name -> proto.QuorumCert
	0,  // 18: proto.BasicHotStuff.SendMsg:input_type -> proto.Msg
	0,  // 19: proto.BasicHotStuff.SendRequest:input_type -> proto.Msg
	0,  // 20: proto.BasicHotStuff.SendReply:input_type -> proto.Msg
	1,  // 21: proto.BasicHotStuff.SendMsg:output_type -> proto.Empty
	1,  // 22: proto.BasicHotStuff.SendRequest:output_type -> proto.Empty
	1,  // 23: proto.BasicHotStuff.SendReply:output_type -> proto.Empty
	21, // [21:24] is the sub-list for method output_type
	18, // [18:21] is the sub-list for method input_type
	18, // [18:18] is the sub-list for extension type_name
	18, // [18:18] is the sub-list for extension extendee
	0,  // [0:18] is the sub-list for field type_name
}

func init() { file_basic_proto_init() }
func file_basic_proto_init() {
	if File_basic_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_basic_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Msg); i {
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
		file_basic_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_basic_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Prepare); i {
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
		file_basic_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrepareVote); i {
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
		file_basic_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreCommit); i {
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
		file_basic_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreCommitVote); i {
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
		file_basic_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Commit); i {
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
		file_basic_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommitVote); i {
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
		file_basic_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Decide); i {
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
		file_basic_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewView); i {
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
		file_basic_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_basic_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Reply); i {
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
	file_basic_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Msg_Prepare)(nil),
		(*Msg_PrepareVote)(nil),
		(*Msg_PreCommit)(nil),
		(*Msg_PreCommitVote)(nil),
		(*Msg_Commit)(nil),
		(*Msg_CommitVote)(nil),
		(*Msg_Decide)(nil),
		(*Msg_NewView)(nil),
		(*Msg_Request)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_basic_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_basic_proto_goTypes,
		DependencyIndexes: file_basic_proto_depIdxs,
		MessageInfos:      file_basic_proto_msgTypes,
	}.Build()
	File_basic_proto = out.File
	file_basic_proto_rawDesc = nil
	file_basic_proto_goTypes = nil
	file_basic_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// BasicHotStuffClient is the client API for BasicHotStuff service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BasicHotStuffClient interface {
	SendMsg(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Empty, error)
	SendRequest(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Empty, error)
	SendReply(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Empty, error)
}

type basicHotStuffClient struct {
	cc grpc.ClientConnInterface
}

func NewBasicHotStuffClient(cc grpc.ClientConnInterface) BasicHotStuffClient {
	return &basicHotStuffClient{cc}
}

func (c *basicHotStuffClient) SendMsg(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.BasicHotStuff/SendMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basicHotStuffClient) SendRequest(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.BasicHotStuff/SendRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basicHotStuffClient) SendReply(ctx context.Context, in *Msg, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.BasicHotStuff/SendReply", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BasicHotStuffServer is the server API for BasicHotStuff service.
type BasicHotStuffServer interface {
	SendMsg(context.Context, *Msg) (*Empty, error)
	SendRequest(context.Context, *Msg) (*Empty, error)
	SendReply(context.Context, *Msg) (*Empty, error)
}

// UnimplementedBasicHotStuffServer can be embedded to have forward compatible implementations.
type UnimplementedBasicHotStuffServer struct {
}

func (*UnimplementedBasicHotStuffServer) SendMsg(context.Context, *Msg) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsg not implemented")
}
func (*UnimplementedBasicHotStuffServer) SendRequest(context.Context, *Msg) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendRequest not implemented")
}
func (*UnimplementedBasicHotStuffServer) SendReply(context.Context, *Msg) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendReply not implemented")
}

func RegisterBasicHotStuffServer(s *grpc.Server, srv BasicHotStuffServer) {
	s.RegisterService(&_BasicHotStuff_serviceDesc, srv)
}

func _BasicHotStuff_SendMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicHotStuffServer).SendMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.BasicHotStuff/SendMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicHotStuffServer).SendMsg(ctx, req.(*Msg))
	}
	return interceptor(ctx, in, info, handler)
}

func _BasicHotStuff_SendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicHotStuffServer).SendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.BasicHotStuff/SendRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicHotStuffServer).SendRequest(ctx, req.(*Msg))
	}
	return interceptor(ctx, in, info, handler)
}

func _BasicHotStuff_SendReply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasicHotStuffServer).SendReply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.BasicHotStuff/SendReply",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasicHotStuffServer).SendReply(ctx, req.(*Msg))
	}
	return interceptor(ctx, in, info, handler)
}

var _BasicHotStuff_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.BasicHotStuff",
	HandlerType: (*BasicHotStuffServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMsg",
			Handler:    _BasicHotStuff_SendMsg_Handler,
		},
		{
			MethodName: "SendRequest",
			Handler:    _BasicHotStuff_SendRequest_Handler,
		},
		{
			MethodName: "SendReply",
			Handler:    _BasicHotStuff_SendReply_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "basic.proto",
}
