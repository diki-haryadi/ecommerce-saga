// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.21.12
// source: proto/cart/cart.proto

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

// CartItem represents an item in the cart
type CartItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ProductId     string                 `protobuf:"bytes,2,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Price         float64                `protobuf:"fixed64,4,opt,name=price,proto3" json:"price,omitempty"`
	Quantity      int32                  `protobuf:"varint,5,opt,name=quantity,proto3" json:"quantity,omitempty"`
	AddedAt       *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=added_at,json=addedAt,proto3" json:"added_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CartItem) Reset() {
	*x = CartItem{}
	mi := &file_proto_cart_cart_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CartItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CartItem) ProtoMessage() {}

func (x *CartItem) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cart_cart_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CartItem.ProtoReflect.Descriptor instead.
func (*CartItem) Descriptor() ([]byte, []int) {
	return file_proto_cart_cart_proto_rawDescGZIP(), []int{0}
}

func (x *CartItem) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CartItem) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *CartItem) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CartItem) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *CartItem) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *CartItem) GetAddedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.AddedAt
	}
	return nil
}

// AddItemRequest represents a request to add an item to the cart
type AddItemRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProductId     string                 `protobuf:"bytes,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Quantity      int32                  `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddItemRequest) Reset() {
	*x = AddItemRequest{}
	mi := &file_proto_cart_cart_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddItemRequest) ProtoMessage() {}

func (x *AddItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cart_cart_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddItemRequest.ProtoReflect.Descriptor instead.
func (*AddItemRequest) Descriptor() ([]byte, []int) {
	return file_proto_cart_cart_proto_rawDescGZIP(), []int{1}
}

func (x *AddItemRequest) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *AddItemRequest) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

// AddItemResponse represents the response after adding an item
type AddItemResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Item          *CartItem              `protobuf:"bytes,3,opt,name=item,proto3" json:"item,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddItemResponse) Reset() {
	*x = AddItemResponse{}
	mi := &file_proto_cart_cart_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddItemResponse) ProtoMessage() {}

func (x *AddItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cart_cart_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddItemResponse.ProtoReflect.Descriptor instead.
func (*AddItemResponse) Descriptor() ([]byte, []int) {
	return file_proto_cart_cart_proto_rawDescGZIP(), []int{2}
}

func (x *AddItemResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *AddItemResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *AddItemResponse) GetItem() *CartItem {
	if x != nil {
		return x.Item
	}
	return nil
}

// RemoveItemRequest represents a request to remove an item from the cart
type RemoveItemRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CartItemId    string                 `protobuf:"bytes,1,opt,name=cart_item_id,json=cartItemId,proto3" json:"cart_item_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RemoveItemRequest) Reset() {
	*x = RemoveItemRequest{}
	mi := &file_proto_cart_cart_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveItemRequest) ProtoMessage() {}

func (x *RemoveItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cart_cart_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveItemRequest.ProtoReflect.Descriptor instead.
func (*RemoveItemRequest) Descriptor() ([]byte, []int) {
	return file_proto_cart_cart_proto_rawDescGZIP(), []int{3}
}

func (x *RemoveItemRequest) GetCartItemId() string {
	if x != nil {
		return x.CartItemId
	}
	return ""
}

// RemoveItemResponse represents the response after removing an item
type RemoveItemResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RemoveItemResponse) Reset() {
	*x = RemoveItemResponse{}
	mi := &file_proto_cart_cart_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveItemResponse) ProtoMessage() {}

func (x *RemoveItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cart_cart_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveItemResponse.ProtoReflect.Descriptor instead.
func (*RemoveItemResponse) Descriptor() ([]byte, []int) {
	return file_proto_cart_cart_proto_rawDescGZIP(), []int{4}
}

func (x *RemoveItemResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *RemoveItemResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// UpdateQuantityRequest represents a request to update item quantity
type UpdateQuantityRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CartItemId    string                 `protobuf:"bytes,1,opt,name=cart_item_id,json=cartItemId,proto3" json:"cart_item_id,omitempty"`
	Quantity      int32                  `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateQuantityRequest) Reset() {
	*x = UpdateQuantityRequest{}
	mi := &file_proto_cart_cart_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateQuantityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateQuantityRequest) ProtoMessage() {}

func (x *UpdateQuantityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cart_cart_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateQuantityRequest.ProtoReflect.Descriptor instead.
func (*UpdateQuantityRequest) Descriptor() ([]byte, []int) {
	return file_proto_cart_cart_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateQuantityRequest) GetCartItemId() string {
	if x != nil {
		return x.CartItemId
	}
	return ""
}

func (x *UpdateQuantityRequest) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

// UpdateQuantityResponse represents the response after updating quantity
type UpdateQuantityResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Item          *CartItem              `protobuf:"bytes,3,opt,name=item,proto3" json:"item,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateQuantityResponse) Reset() {
	*x = UpdateQuantityResponse{}
	mi := &file_proto_cart_cart_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateQuantityResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateQuantityResponse) ProtoMessage() {}

func (x *UpdateQuantityResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cart_cart_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateQuantityResponse.ProtoReflect.Descriptor instead.
func (*UpdateQuantityResponse) Descriptor() ([]byte, []int) {
	return file_proto_cart_cart_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateQuantityResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UpdateQuantityResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *UpdateQuantityResponse) GetItem() *CartItem {
	if x != nil {
		return x.Item
	}
	return nil
}

// GetCartRequest represents a request to get the current cart
type GetCartRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCartRequest) Reset() {
	*x = GetCartRequest{}
	mi := &file_proto_cart_cart_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCartRequest) ProtoMessage() {}

func (x *GetCartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cart_cart_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCartRequest.ProtoReflect.Descriptor instead.
func (*GetCartRequest) Descriptor() ([]byte, []int) {
	return file_proto_cart_cart_proto_rawDescGZIP(), []int{7}
}

// GetCartResponse represents the response containing cart details
type GetCartResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Items         []*CartItem            `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	TotalAmount   float64                `protobuf:"fixed64,2,opt,name=total_amount,json=totalAmount,proto3" json:"total_amount,omitempty"`
	TotalItems    int32                  `protobuf:"varint,3,opt,name=total_items,json=totalItems,proto3" json:"total_items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCartResponse) Reset() {
	*x = GetCartResponse{}
	mi := &file_proto_cart_cart_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCartResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCartResponse) ProtoMessage() {}

func (x *GetCartResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cart_cart_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCartResponse.ProtoReflect.Descriptor instead.
func (*GetCartResponse) Descriptor() ([]byte, []int) {
	return file_proto_cart_cart_proto_rawDescGZIP(), []int{8}
}

func (x *GetCartResponse) GetItems() []*CartItem {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *GetCartResponse) GetTotalAmount() float64 {
	if x != nil {
		return x.TotalAmount
	}
	return 0
}

func (x *GetCartResponse) GetTotalItems() int32 {
	if x != nil {
		return x.TotalItems
	}
	return 0
}

// ClearCartRequest represents a request to clear the cart
type ClearCartRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ClearCartRequest) Reset() {
	*x = ClearCartRequest{}
	mi := &file_proto_cart_cart_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClearCartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClearCartRequest) ProtoMessage() {}

func (x *ClearCartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cart_cart_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClearCartRequest.ProtoReflect.Descriptor instead.
func (*ClearCartRequest) Descriptor() ([]byte, []int) {
	return file_proto_cart_cart_proto_rawDescGZIP(), []int{9}
}

// ClearCartResponse represents the response after clearing the cart
type ClearCartResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ClearCartResponse) Reset() {
	*x = ClearCartResponse{}
	mi := &file_proto_cart_cart_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClearCartResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClearCartResponse) ProtoMessage() {}

func (x *ClearCartResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cart_cart_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClearCartResponse.ProtoReflect.Descriptor instead.
func (*ClearCartResponse) Descriptor() ([]byte, []int) {
	return file_proto_cart_cart_proto_rawDescGZIP(), []int{10}
}

func (x *ClearCartResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *ClearCartResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_proto_cart_cart_proto protoreflect.FileDescriptor

const file_proto_cart_cart_proto_rawDesc = "" +
	"\n" +
	"\x15proto/cart/cart.proto\x12\x04cart\x1a\x1fgoogle/protobuf/timestamp.proto\"\xb6\x01\n" +
	"\bCartItem\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x1d\n" +
	"\n" +
	"product_id\x18\x02 \x01(\tR\tproductId\x12\x12\n" +
	"\x04name\x18\x03 \x01(\tR\x04name\x12\x14\n" +
	"\x05price\x18\x04 \x01(\x01R\x05price\x12\x1a\n" +
	"\bquantity\x18\x05 \x01(\x05R\bquantity\x125\n" +
	"\badded_at\x18\x06 \x01(\v2\x1a.google.protobuf.TimestampR\aaddedAt\"K\n" +
	"\x0eAddItemRequest\x12\x1d\n" +
	"\n" +
	"product_id\x18\x01 \x01(\tR\tproductId\x12\x1a\n" +
	"\bquantity\x18\x02 \x01(\x05R\bquantity\"i\n" +
	"\x0fAddItemResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\x12\"\n" +
	"\x04item\x18\x03 \x01(\v2\x0e.cart.CartItemR\x04item\"5\n" +
	"\x11RemoveItemRequest\x12 \n" +
	"\fcart_item_id\x18\x01 \x01(\tR\n" +
	"cartItemId\"H\n" +
	"\x12RemoveItemResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\"U\n" +
	"\x15UpdateQuantityRequest\x12 \n" +
	"\fcart_item_id\x18\x01 \x01(\tR\n" +
	"cartItemId\x12\x1a\n" +
	"\bquantity\x18\x02 \x01(\x05R\bquantity\"p\n" +
	"\x16UpdateQuantityResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\x12\"\n" +
	"\x04item\x18\x03 \x01(\v2\x0e.cart.CartItemR\x04item\"\x10\n" +
	"\x0eGetCartRequest\"{\n" +
	"\x0fGetCartResponse\x12$\n" +
	"\x05items\x18\x01 \x03(\v2\x0e.cart.CartItemR\x05items\x12!\n" +
	"\ftotal_amount\x18\x02 \x01(\x01R\vtotalAmount\x12\x1f\n" +
	"\vtotal_items\x18\x03 \x01(\x05R\n" +
	"totalItems\"\x12\n" +
	"\x10ClearCartRequest\"G\n" +
	"\x11ClearCartResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage2\xd3\x02\n" +
	"\vCartService\x128\n" +
	"\aAddItem\x12\x14.cart.AddItemRequest\x1a\x15.cart.AddItemResponse\"\x00\x12A\n" +
	"\n" +
	"RemoveItem\x12\x17.cart.RemoveItemRequest\x1a\x18.cart.RemoveItemResponse\"\x00\x12M\n" +
	"\x0eUpdateQuantity\x12\x1b.cart.UpdateQuantityRequest\x1a\x1c.cart.UpdateQuantityResponse\"\x00\x128\n" +
	"\aGetCart\x12\x14.cart.GetCartRequest\x1a\x15.cart.GetCartResponse\"\x00\x12>\n" +
	"\tClearCart\x12\x16.cart.ClearCartRequest\x1a\x17.cart.ClearCartResponse\"\x00BSZQgithub.com/diki-haryadi/ecommerce-saga/internal/features/cart/delivery/grpc/protob\x06proto3"

var (
	file_proto_cart_cart_proto_rawDescOnce sync.Once
	file_proto_cart_cart_proto_rawDescData []byte
)

func file_proto_cart_cart_proto_rawDescGZIP() []byte {
	file_proto_cart_cart_proto_rawDescOnce.Do(func() {
		file_proto_cart_cart_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_cart_cart_proto_rawDesc), len(file_proto_cart_cart_proto_rawDesc)))
	})
	return file_proto_cart_cart_proto_rawDescData
}

var file_proto_cart_cart_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_proto_cart_cart_proto_goTypes = []any{
	(*CartItem)(nil),               // 0: cart.CartItem
	(*AddItemRequest)(nil),         // 1: cart.AddItemRequest
	(*AddItemResponse)(nil),        // 2: cart.AddItemResponse
	(*RemoveItemRequest)(nil),      // 3: cart.RemoveItemRequest
	(*RemoveItemResponse)(nil),     // 4: cart.RemoveItemResponse
	(*UpdateQuantityRequest)(nil),  // 5: cart.UpdateQuantityRequest
	(*UpdateQuantityResponse)(nil), // 6: cart.UpdateQuantityResponse
	(*GetCartRequest)(nil),         // 7: cart.GetCartRequest
	(*GetCartResponse)(nil),        // 8: cart.GetCartResponse
	(*ClearCartRequest)(nil),       // 9: cart.ClearCartRequest
	(*ClearCartResponse)(nil),      // 10: cart.ClearCartResponse
	(*timestamppb.Timestamp)(nil),  // 11: google.protobuf.Timestamp
}
var file_proto_cart_cart_proto_depIdxs = []int32{
	11, // 0: cart.CartItem.added_at:type_name -> google.protobuf.Timestamp
	0,  // 1: cart.AddItemResponse.item:type_name -> cart.CartItem
	0,  // 2: cart.UpdateQuantityResponse.item:type_name -> cart.CartItem
	0,  // 3: cart.GetCartResponse.items:type_name -> cart.CartItem
	1,  // 4: cart.CartService.AddItem:input_type -> cart.AddItemRequest
	3,  // 5: cart.CartService.RemoveItem:input_type -> cart.RemoveItemRequest
	5,  // 6: cart.CartService.UpdateQuantity:input_type -> cart.UpdateQuantityRequest
	7,  // 7: cart.CartService.GetCart:input_type -> cart.GetCartRequest
	9,  // 8: cart.CartService.ClearCart:input_type -> cart.ClearCartRequest
	2,  // 9: cart.CartService.AddItem:output_type -> cart.AddItemResponse
	4,  // 10: cart.CartService.RemoveItem:output_type -> cart.RemoveItemResponse
	6,  // 11: cart.CartService.UpdateQuantity:output_type -> cart.UpdateQuantityResponse
	8,  // 12: cart.CartService.GetCart:output_type -> cart.GetCartResponse
	10, // 13: cart.CartService.ClearCart:output_type -> cart.ClearCartResponse
	9,  // [9:14] is the sub-list for method output_type
	4,  // [4:9] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_proto_cart_cart_proto_init() }
func file_proto_cart_cart_proto_init() {
	if File_proto_cart_cart_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_cart_cart_proto_rawDesc), len(file_proto_cart_cart_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_cart_cart_proto_goTypes,
		DependencyIndexes: file_proto_cart_cart_proto_depIdxs,
		MessageInfos:      file_proto_cart_cart_proto_msgTypes,
	}.Build()
	File_proto_cart_cart_proto = out.File
	file_proto_cart_cart_proto_goTypes = nil
	file_proto_cart_cart_proto_depIdxs = nil
}
