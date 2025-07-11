// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.21.12
// source: proto/saga/saga.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// SagaTransaction represents a distributed transaction
type SagaTransaction struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type          string                 `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Status        string                 `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	Steps         []*SagaStep            `protobuf:"bytes,4,rep,name=steps,proto3" json:"steps,omitempty"`
	Metadata      map[string]string      `protobuf:"bytes,5,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SagaTransaction) Reset() {
	*x = SagaTransaction{}
	mi := &file_proto_saga_saga_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SagaTransaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SagaTransaction) ProtoMessage() {}

func (x *SagaTransaction) ProtoReflect() protoreflect.Message {
	mi := &file_proto_saga_saga_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SagaTransaction.ProtoReflect.Descriptor instead.
func (*SagaTransaction) Descriptor() ([]byte, []int) {
	return file_proto_saga_saga_proto_rawDescGZIP(), []int{0}
}

func (x *SagaTransaction) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SagaTransaction) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *SagaTransaction) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *SagaTransaction) GetSteps() []*SagaStep {
	if x != nil {
		return x.Steps
	}
	return nil
}

func (x *SagaTransaction) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *SagaTransaction) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *SagaTransaction) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// SagaStep represents a step in the saga transaction
type SagaStep struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	Id                 string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name               string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Status             string                 `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	Service            string                 `protobuf:"bytes,4,opt,name=service,proto3" json:"service,omitempty"`
	Action             string                 `protobuf:"bytes,5,opt,name=action,proto3" json:"action,omitempty"`
	CompensationAction string                 `protobuf:"bytes,6,opt,name=compensation_action,json=compensationAction,proto3" json:"compensation_action,omitempty"`
	Payload            map[string]string      `protobuf:"bytes,7,rep,name=payload,proto3" json:"payload,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	ErrorMessage       string                 `protobuf:"bytes,8,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	ExecutedAt         *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=executed_at,json=executedAt,proto3" json:"executed_at,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *SagaStep) Reset() {
	*x = SagaStep{}
	mi := &file_proto_saga_saga_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SagaStep) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SagaStep) ProtoMessage() {}

func (x *SagaStep) ProtoReflect() protoreflect.Message {
	mi := &file_proto_saga_saga_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SagaStep.ProtoReflect.Descriptor instead.
func (*SagaStep) Descriptor() ([]byte, []int) {
	return file_proto_saga_saga_proto_rawDescGZIP(), []int{1}
}

func (x *SagaStep) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SagaStep) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SagaStep) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *SagaStep) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *SagaStep) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *SagaStep) GetCompensationAction() string {
	if x != nil {
		return x.CompensationAction
	}
	return ""
}

func (x *SagaStep) GetPayload() map[string]string {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *SagaStep) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

func (x *SagaStep) GetExecutedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ExecutedAt
	}
	return nil
}

// StartOrderSagaRequest represents a request to start an order saga
type StartOrderSagaRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderId       string                 `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	UserId        string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Amount        float64                `protobuf:"fixed64,3,opt,name=amount,proto3" json:"amount,omitempty"`
	PaymentMethod string                 `protobuf:"bytes,4,opt,name=payment_method,json=paymentMethod,proto3" json:"payment_method,omitempty"`
	Metadata      map[string]string      `protobuf:"bytes,5,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StartOrderSagaRequest) Reset() {
	*x = StartOrderSagaRequest{}
	mi := &file_proto_saga_saga_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StartOrderSagaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartOrderSagaRequest) ProtoMessage() {}

func (x *StartOrderSagaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_saga_saga_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartOrderSagaRequest.ProtoReflect.Descriptor instead.
func (*StartOrderSagaRequest) Descriptor() ([]byte, []int) {
	return file_proto_saga_saga_proto_rawDescGZIP(), []int{2}
}

func (x *StartOrderSagaRequest) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *StartOrderSagaRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *StartOrderSagaRequest) GetAmount() float64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *StartOrderSagaRequest) GetPaymentMethod() string {
	if x != nil {
		return x.PaymentMethod
	}
	return ""
}

func (x *StartOrderSagaRequest) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// StartOrderSagaResponse represents the response after starting a saga
type StartOrderSagaResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Transaction   *SagaTransaction       `protobuf:"bytes,3,opt,name=transaction,proto3" json:"transaction,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StartOrderSagaResponse) Reset() {
	*x = StartOrderSagaResponse{}
	mi := &file_proto_saga_saga_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StartOrderSagaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartOrderSagaResponse) ProtoMessage() {}

func (x *StartOrderSagaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_saga_saga_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartOrderSagaResponse.ProtoReflect.Descriptor instead.
func (*StartOrderSagaResponse) Descriptor() ([]byte, []int) {
	return file_proto_saga_saga_proto_rawDescGZIP(), []int{3}
}

func (x *StartOrderSagaResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *StartOrderSagaResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *StartOrderSagaResponse) GetTransaction() *SagaTransaction {
	if x != nil {
		return x.Transaction
	}
	return nil
}

// GetSagaStatusRequest represents a request to get saga status
type GetSagaStatusRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SagaId        string                 `protobuf:"bytes,1,opt,name=saga_id,json=sagaId,proto3" json:"saga_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetSagaStatusRequest) Reset() {
	*x = GetSagaStatusRequest{}
	mi := &file_proto_saga_saga_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSagaStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSagaStatusRequest) ProtoMessage() {}

func (x *GetSagaStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_saga_saga_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSagaStatusRequest.ProtoReflect.Descriptor instead.
func (*GetSagaStatusRequest) Descriptor() ([]byte, []int) {
	return file_proto_saga_saga_proto_rawDescGZIP(), []int{4}
}

func (x *GetSagaStatusRequest) GetSagaId() string {
	if x != nil {
		return x.SagaId
	}
	return ""
}

// GetSagaStatusResponse represents the response containing saga status
type GetSagaStatusResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Transaction   *SagaTransaction       `protobuf:"bytes,1,opt,name=transaction,proto3" json:"transaction,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetSagaStatusResponse) Reset() {
	*x = GetSagaStatusResponse{}
	mi := &file_proto_saga_saga_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSagaStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSagaStatusResponse) ProtoMessage() {}

func (x *GetSagaStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_saga_saga_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSagaStatusResponse.ProtoReflect.Descriptor instead.
func (*GetSagaStatusResponse) Descriptor() ([]byte, []int) {
	return file_proto_saga_saga_proto_rawDescGZIP(), []int{5}
}

func (x *GetSagaStatusResponse) GetTransaction() *SagaTransaction {
	if x != nil {
		return x.Transaction
	}
	return nil
}

// CompensateTransactionRequest represents a request to compensate a transaction
type CompensateTransactionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SagaId        string                 `protobuf:"bytes,1,opt,name=saga_id,json=sagaId,proto3" json:"saga_id,omitempty"`
	StepId        string                 `protobuf:"bytes,2,opt,name=step_id,json=stepId,proto3" json:"step_id,omitempty"`
	Reason        string                 `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CompensateTransactionRequest) Reset() {
	*x = CompensateTransactionRequest{}
	mi := &file_proto_saga_saga_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CompensateTransactionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompensateTransactionRequest) ProtoMessage() {}

func (x *CompensateTransactionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_saga_saga_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompensateTransactionRequest.ProtoReflect.Descriptor instead.
func (*CompensateTransactionRequest) Descriptor() ([]byte, []int) {
	return file_proto_saga_saga_proto_rawDescGZIP(), []int{6}
}

func (x *CompensateTransactionRequest) GetSagaId() string {
	if x != nil {
		return x.SagaId
	}
	return ""
}

func (x *CompensateTransactionRequest) GetStepId() string {
	if x != nil {
		return x.StepId
	}
	return ""
}

func (x *CompensateTransactionRequest) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

// CompensateTransactionResponse represents the response after compensation
type CompensateTransactionResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Transaction   *SagaTransaction       `protobuf:"bytes,3,opt,name=transaction,proto3" json:"transaction,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CompensateTransactionResponse) Reset() {
	*x = CompensateTransactionResponse{}
	mi := &file_proto_saga_saga_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CompensateTransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompensateTransactionResponse) ProtoMessage() {}

func (x *CompensateTransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_saga_saga_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompensateTransactionResponse.ProtoReflect.Descriptor instead.
func (*CompensateTransactionResponse) Descriptor() ([]byte, []int) {
	return file_proto_saga_saga_proto_rawDescGZIP(), []int{7}
}

func (x *CompensateTransactionResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *CompensateTransactionResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *CompensateTransactionResponse) GetTransaction() *SagaTransaction {
	if x != nil {
		return x.Transaction
	}
	return nil
}

// ListSagaTransactionsRequest represents a request to list transactions
type ListSagaTransactionsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Limit         int32                  `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Status        string                 `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	Type          string                 `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListSagaTransactionsRequest) Reset() {
	*x = ListSagaTransactionsRequest{}
	mi := &file_proto_saga_saga_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListSagaTransactionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSagaTransactionsRequest) ProtoMessage() {}

func (x *ListSagaTransactionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_saga_saga_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSagaTransactionsRequest.ProtoReflect.Descriptor instead.
func (*ListSagaTransactionsRequest) Descriptor() ([]byte, []int) {
	return file_proto_saga_saga_proto_rawDescGZIP(), []int{8}
}

func (x *ListSagaTransactionsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListSagaTransactionsRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *ListSagaTransactionsRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ListSagaTransactionsRequest) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

// ListSagaTransactionsResponse represents the response containing transactions
type ListSagaTransactionsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Transactions  []*SagaTransaction     `protobuf:"bytes,1,rep,name=transactions,proto3" json:"transactions,omitempty"`
	Total         int32                  `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
	Page          int32                  `protobuf:"varint,3,opt,name=page,proto3" json:"page,omitempty"`
	Limit         int32                  `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListSagaTransactionsResponse) Reset() {
	*x = ListSagaTransactionsResponse{}
	mi := &file_proto_saga_saga_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListSagaTransactionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSagaTransactionsResponse) ProtoMessage() {}

func (x *ListSagaTransactionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_saga_saga_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSagaTransactionsResponse.ProtoReflect.Descriptor instead.
func (*ListSagaTransactionsResponse) Descriptor() ([]byte, []int) {
	return file_proto_saga_saga_proto_rawDescGZIP(), []int{9}
}

func (x *ListSagaTransactionsResponse) GetTransactions() []*SagaTransaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

func (x *ListSagaTransactionsResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *ListSagaTransactionsResponse) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListSagaTransactionsResponse) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

var File_proto_saga_saga_proto protoreflect.FileDescriptor

const file_proto_saga_saga_proto_rawDesc = "" +
	"\n" +
	"\x15proto/saga/saga.proto\x12\x04saga\x1a\x1fgoogle/protobuf/timestamp.proto\"\xe7\x02\n" +
	"\x0fSagaTransaction\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04type\x18\x02 \x01(\tR\x04type\x12\x16\n" +
	"\x06status\x18\x03 \x01(\tR\x06status\x12$\n" +
	"\x05steps\x18\x04 \x03(\v2\x0e.saga.SagaStepR\x05steps\x12?\n" +
	"\bmetadata\x18\x05 \x03(\v2#.saga.SagaTransaction.MetadataEntryR\bmetadata\x129\n" +
	"\n" +
	"created_at\x18\x06 \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\a \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt\x1a;\n" +
	"\rMetadataEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01\"\xfe\x02\n" +
	"\bSagaStep\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x16\n" +
	"\x06status\x18\x03 \x01(\tR\x06status\x12\x18\n" +
	"\aservice\x18\x04 \x01(\tR\aservice\x12\x16\n" +
	"\x06action\x18\x05 \x01(\tR\x06action\x12/\n" +
	"\x13compensation_action\x18\x06 \x01(\tR\x12compensationAction\x125\n" +
	"\apayload\x18\a \x03(\v2\x1b.saga.SagaStep.PayloadEntryR\apayload\x12#\n" +
	"\rerror_message\x18\b \x01(\tR\ferrorMessage\x12;\n" +
	"\vexecuted_at\x18\t \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"executedAt\x1a:\n" +
	"\fPayloadEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01\"\x8e\x02\n" +
	"\x15StartOrderSagaRequest\x12\x19\n" +
	"\border_id\x18\x01 \x01(\tR\aorderId\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\tR\x06userId\x12\x16\n" +
	"\x06amount\x18\x03 \x01(\x01R\x06amount\x12%\n" +
	"\x0epayment_method\x18\x04 \x01(\tR\rpaymentMethod\x12E\n" +
	"\bmetadata\x18\x05 \x03(\v2).saga.StartOrderSagaRequest.MetadataEntryR\bmetadata\x1a;\n" +
	"\rMetadataEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01\"\x85\x01\n" +
	"\x16StartOrderSagaResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\x127\n" +
	"\vtransaction\x18\x03 \x01(\v2\x15.saga.SagaTransactionR\vtransaction\"/\n" +
	"\x14GetSagaStatusRequest\x12\x17\n" +
	"\asaga_id\x18\x01 \x01(\tR\x06sagaId\"P\n" +
	"\x15GetSagaStatusResponse\x127\n" +
	"\vtransaction\x18\x01 \x01(\v2\x15.saga.SagaTransactionR\vtransaction\"h\n" +
	"\x1cCompensateTransactionRequest\x12\x17\n" +
	"\asaga_id\x18\x01 \x01(\tR\x06sagaId\x12\x17\n" +
	"\astep_id\x18\x02 \x01(\tR\x06stepId\x12\x16\n" +
	"\x06reason\x18\x03 \x01(\tR\x06reason\"\x8c\x01\n" +
	"\x1dCompensateTransactionResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\x127\n" +
	"\vtransaction\x18\x03 \x01(\v2\x15.saga.SagaTransactionR\vtransaction\"s\n" +
	"\x1bListSagaTransactionsRequest\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x14\n" +
	"\x05limit\x18\x02 \x01(\x05R\x05limit\x12\x16\n" +
	"\x06status\x18\x03 \x01(\tR\x06status\x12\x12\n" +
	"\x04type\x18\x04 \x01(\tR\x04type\"\x99\x01\n" +
	"\x1cListSagaTransactionsResponse\x129\n" +
	"\ftransactions\x18\x01 \x03(\v2\x15.saga.SagaTransactionR\ftransactions\x12\x14\n" +
	"\x05total\x18\x02 \x01(\x05R\x05total\x12\x12\n" +
	"\x04page\x18\x03 \x01(\x05R\x04page\x12\x14\n" +
	"\x05limit\x18\x04 \x01(\x05R\x05limit2\xed\x02\n" +
	"\vSagaService\x12M\n" +
	"\x0eStartOrderSaga\x12\x1b.saga.StartOrderSagaRequest\x1a\x1c.saga.StartOrderSagaResponse\"\x00\x12J\n" +
	"\rGetSagaStatus\x12\x1a.saga.GetSagaStatusRequest\x1a\x1b.saga.GetSagaStatusResponse\"\x00\x12b\n" +
	"\x15CompensateTransaction\x12\".saga.CompensateTransactionRequest\x1a#.saga.CompensateTransactionResponse\"\x00\x12_\n" +
	"\x14ListSagaTransactions\x12!.saga.ListSagaTransactionsRequest\x1a\".saga.ListSagaTransactionsResponse\"\x00BSZQgithub.com/diki-haryadi/ecommerce-saga/internal/features/saga/delivery/grpc/protob\x06proto3"

var (
	file_proto_saga_saga_proto_rawDescOnce sync.Once
	file_proto_saga_saga_proto_rawDescData []byte
)

func file_proto_saga_saga_proto_rawDescGZIP() []byte {
	file_proto_saga_saga_proto_rawDescOnce.Do(func() {
		file_proto_saga_saga_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_saga_saga_proto_rawDesc), len(file_proto_saga_saga_proto_rawDesc)))
	})
	return file_proto_saga_saga_proto_rawDescData
}

var file_proto_saga_saga_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_proto_saga_saga_proto_goTypes = []any{
	(*SagaTransaction)(nil),               // 0: saga.SagaTransaction
	(*SagaStep)(nil),                      // 1: saga.SagaStep
	(*StartOrderSagaRequest)(nil),         // 2: saga.StartOrderSagaRequest
	(*StartOrderSagaResponse)(nil),        // 3: saga.StartOrderSagaResponse
	(*GetSagaStatusRequest)(nil),          // 4: saga.GetSagaStatusRequest
	(*GetSagaStatusResponse)(nil),         // 5: saga.GetSagaStatusResponse
	(*CompensateTransactionRequest)(nil),  // 6: saga.CompensateTransactionRequest
	(*CompensateTransactionResponse)(nil), // 7: saga.CompensateTransactionResponse
	(*ListSagaTransactionsRequest)(nil),   // 8: saga.ListSagaTransactionsRequest
	(*ListSagaTransactionsResponse)(nil),  // 9: saga.ListSagaTransactionsResponse
	nil,                                   // 10: saga.SagaTransaction.MetadataEntry
	nil,                                   // 11: saga.SagaStep.PayloadEntry
	nil,                                   // 12: saga.StartOrderSagaRequest.MetadataEntry
	(*timestamppb.Timestamp)(nil),         // 13: google.protobuf.Timestamp
}
var file_proto_saga_saga_proto_depIdxs = []int32{
	1,  // 0: saga.SagaTransaction.steps:type_name -> saga.SagaStep
	10, // 1: saga.SagaTransaction.metadata:type_name -> saga.SagaTransaction.MetadataEntry
	13, // 2: saga.SagaTransaction.created_at:type_name -> google.protobuf.Timestamp
	13, // 3: saga.SagaTransaction.updated_at:type_name -> google.protobuf.Timestamp
	11, // 4: saga.SagaStep.payload:type_name -> saga.SagaStep.PayloadEntry
	13, // 5: saga.SagaStep.executed_at:type_name -> google.protobuf.Timestamp
	12, // 6: saga.StartOrderSagaRequest.metadata:type_name -> saga.StartOrderSagaRequest.MetadataEntry
	0,  // 7: saga.StartOrderSagaResponse.transaction:type_name -> saga.SagaTransaction
	0,  // 8: saga.GetSagaStatusResponse.transaction:type_name -> saga.SagaTransaction
	0,  // 9: saga.CompensateTransactionResponse.transaction:type_name -> saga.SagaTransaction
	0,  // 10: saga.ListSagaTransactionsResponse.transactions:type_name -> saga.SagaTransaction
	2,  // 11: saga.SagaService.StartOrderSaga:input_type -> saga.StartOrderSagaRequest
	4,  // 12: saga.SagaService.GetSagaStatus:input_type -> saga.GetSagaStatusRequest
	6,  // 13: saga.SagaService.CompensateTransaction:input_type -> saga.CompensateTransactionRequest
	8,  // 14: saga.SagaService.ListSagaTransactions:input_type -> saga.ListSagaTransactionsRequest
	3,  // 15: saga.SagaService.StartOrderSaga:output_type -> saga.StartOrderSagaResponse
	5,  // 16: saga.SagaService.GetSagaStatus:output_type -> saga.GetSagaStatusResponse
	7,  // 17: saga.SagaService.CompensateTransaction:output_type -> saga.CompensateTransactionResponse
	9,  // 18: saga.SagaService.ListSagaTransactions:output_type -> saga.ListSagaTransactionsResponse
	15, // [15:19] is the sub-list for method output_type
	11, // [11:15] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_proto_saga_saga_proto_init() }
func file_proto_saga_saga_proto_init() {
	if File_proto_saga_saga_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_saga_saga_proto_rawDesc), len(file_proto_saga_saga_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_saga_saga_proto_goTypes,
		DependencyIndexes: file_proto_saga_saga_proto_depIdxs,
		MessageInfos:      file_proto_saga_saga_proto_msgTypes,
	}.Build()
	File_proto_saga_saga_proto = out.File
	file_proto_saga_saga_proto_goTypes = nil
	file_proto_saga_saga_proto_depIdxs = nil
}
