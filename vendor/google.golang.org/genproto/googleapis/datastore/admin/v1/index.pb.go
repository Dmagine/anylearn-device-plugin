// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/datastore/admin/v1/index.proto

package admin

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// For an ordered index, specifies whether each of the entity's ancestors
// will be included.
type Index_AncestorMode int32

const (
	// The ancestor mode is unspecified.
	Index_ANCESTOR_MODE_UNSPECIFIED Index_AncestorMode = 0
	// Do not include the entity's ancestors in the index.
	Index_NONE Index_AncestorMode = 1
	// Include all the entity's ancestors in the index.
	Index_ALL_ANCESTORS Index_AncestorMode = 2
)

var Index_AncestorMode_name = map[int32]string{
	0: "ANCESTOR_MODE_UNSPECIFIED",
	1: "NONE",
	2: "ALL_ANCESTORS",
}

var Index_AncestorMode_value = map[string]int32{
	"ANCESTOR_MODE_UNSPECIFIED": 0,
	"NONE":                      1,
	"ALL_ANCESTORS":             2,
}

func (x Index_AncestorMode) String() string {
	return proto.EnumName(Index_AncestorMode_name, int32(x))
}

func (Index_AncestorMode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_809cc5775e1cdeb3, []int{0, 0}
}

// The direction determines how a property is indexed.
type Index_Direction int32

const (
	// The direction is unspecified.
	Index_DIRECTION_UNSPECIFIED Index_Direction = 0
	// The property's values are indexed so as to support sequencing in
	// ascending order and also query by <, >, <=, >=, and =.
	Index_ASCENDING Index_Direction = 1
	// The property's values are indexed so as to support sequencing in
	// descending order and also query by <, >, <=, >=, and =.
	Index_DESCENDING Index_Direction = 2
)

var Index_Direction_name = map[int32]string{
	0: "DIRECTION_UNSPECIFIED",
	1: "ASCENDING",
	2: "DESCENDING",
}

var Index_Direction_value = map[string]int32{
	"DIRECTION_UNSPECIFIED": 0,
	"ASCENDING":             1,
	"DESCENDING":            2,
}

func (x Index_Direction) String() string {
	return proto.EnumName(Index_Direction_name, int32(x))
}

func (Index_Direction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_809cc5775e1cdeb3, []int{0, 1}
}

// The possible set of states of an index.
type Index_State int32

const (
	// The state is unspecified.
	Index_STATE_UNSPECIFIED Index_State = 0
	// The index is being created, and cannot be used by queries.
	// There is an active long-running operation for the index.
	// The index is updated when writing an entity.
	// Some index data may exist.
	Index_CREATING Index_State = 1
	// The index is ready to be used.
	// The index is updated when writing an entity.
	// The index is fully populated from all stored entities it applies to.
	Index_READY Index_State = 2
	// The index is being deleted, and cannot be used by queries.
	// There is an active long-running operation for the index.
	// The index is not updated when writing an entity.
	// Some index data may exist.
	Index_DELETING Index_State = 3
	// The index was being created or deleted, but something went wrong.
	// The index cannot by used by queries.
	// There is no active long-running operation for the index,
	// and the most recently finished long-running operation failed.
	// The index is not updated when writing an entity.
	// Some index data may exist.
	Index_ERROR Index_State = 4
)

var Index_State_name = map[int32]string{
	0: "STATE_UNSPECIFIED",
	1: "CREATING",
	2: "READY",
	3: "DELETING",
	4: "ERROR",
}

var Index_State_value = map[string]int32{
	"STATE_UNSPECIFIED": 0,
	"CREATING":          1,
	"READY":             2,
	"DELETING":          3,
	"ERROR":             4,
}

func (x Index_State) String() string {
	return proto.EnumName(Index_State_name, int32(x))
}

func (Index_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_809cc5775e1cdeb3, []int{0, 2}
}

// A minimal index definition.
type Index struct {
	// Output only. Project ID.
	ProjectId string `protobuf:"bytes,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	// Output only. The resource ID of the index.
	IndexId string `protobuf:"bytes,3,opt,name=index_id,json=indexId,proto3" json:"index_id,omitempty"`
	// Required. The entity kind to which this index applies.
	Kind string `protobuf:"bytes,4,opt,name=kind,proto3" json:"kind,omitempty"`
	// Required. The index's ancestor mode.  Must not be ANCESTOR_MODE_UNSPECIFIED.
	Ancestor Index_AncestorMode `protobuf:"varint,5,opt,name=ancestor,proto3,enum=google.datastore.admin.v1.Index_AncestorMode" json:"ancestor,omitempty"`
	// Required. An ordered sequence of property names and their index attributes.
	Properties []*Index_IndexedProperty `protobuf:"bytes,6,rep,name=properties,proto3" json:"properties,omitempty"`
	// Output only. The state of the index.
	State                Index_State `protobuf:"varint,7,opt,name=state,proto3,enum=google.datastore.admin.v1.Index_State" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Index) Reset()         { *m = Index{} }
func (m *Index) String() string { return proto.CompactTextString(m) }
func (*Index) ProtoMessage()    {}
func (*Index) Descriptor() ([]byte, []int) {
	return fileDescriptor_809cc5775e1cdeb3, []int{0}
}

func (m *Index) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Index.Unmarshal(m, b)
}
func (m *Index) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Index.Marshal(b, m, deterministic)
}
func (m *Index) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Index.Merge(m, src)
}
func (m *Index) XXX_Size() int {
	return xxx_messageInfo_Index.Size(m)
}
func (m *Index) XXX_DiscardUnknown() {
	xxx_messageInfo_Index.DiscardUnknown(m)
}

var xxx_messageInfo_Index proto.InternalMessageInfo

func (m *Index) GetProjectId() string {
	if m != nil {
		return m.ProjectId
	}
	return ""
}

func (m *Index) GetIndexId() string {
	if m != nil {
		return m.IndexId
	}
	return ""
}

func (m *Index) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *Index) GetAncestor() Index_AncestorMode {
	if m != nil {
		return m.Ancestor
	}
	return Index_ANCESTOR_MODE_UNSPECIFIED
}

func (m *Index) GetProperties() []*Index_IndexedProperty {
	if m != nil {
		return m.Properties
	}
	return nil
}

func (m *Index) GetState() Index_State {
	if m != nil {
		return m.State
	}
	return Index_STATE_UNSPECIFIED
}

// A property of an index.
type Index_IndexedProperty struct {
	// Required. The property name to index.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Required. The indexed property's direction.  Must not be DIRECTION_UNSPECIFIED.
	Direction            Index_Direction `protobuf:"varint,2,opt,name=direction,proto3,enum=google.datastore.admin.v1.Index_Direction" json:"direction,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Index_IndexedProperty) Reset()         { *m = Index_IndexedProperty{} }
func (m *Index_IndexedProperty) String() string { return proto.CompactTextString(m) }
func (*Index_IndexedProperty) ProtoMessage()    {}
func (*Index_IndexedProperty) Descriptor() ([]byte, []int) {
	return fileDescriptor_809cc5775e1cdeb3, []int{0, 0}
}

func (m *Index_IndexedProperty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Index_IndexedProperty.Unmarshal(m, b)
}
func (m *Index_IndexedProperty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Index_IndexedProperty.Marshal(b, m, deterministic)
}
func (m *Index_IndexedProperty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Index_IndexedProperty.Merge(m, src)
}
func (m *Index_IndexedProperty) XXX_Size() int {
	return xxx_messageInfo_Index_IndexedProperty.Size(m)
}
func (m *Index_IndexedProperty) XXX_DiscardUnknown() {
	xxx_messageInfo_Index_IndexedProperty.DiscardUnknown(m)
}

var xxx_messageInfo_Index_IndexedProperty proto.InternalMessageInfo

func (m *Index_IndexedProperty) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Index_IndexedProperty) GetDirection() Index_Direction {
	if m != nil {
		return m.Direction
	}
	return Index_DIRECTION_UNSPECIFIED
}

func init() {
	proto.RegisterEnum("google.datastore.admin.v1.Index_AncestorMode", Index_AncestorMode_name, Index_AncestorMode_value)
	proto.RegisterEnum("google.datastore.admin.v1.Index_Direction", Index_Direction_name, Index_Direction_value)
	proto.RegisterEnum("google.datastore.admin.v1.Index_State", Index_State_name, Index_State_value)
	proto.RegisterType((*Index)(nil), "google.datastore.admin.v1.Index")
	proto.RegisterType((*Index_IndexedProperty)(nil), "google.datastore.admin.v1.Index.IndexedProperty")
}

func init() {
	proto.RegisterFile("google/datastore/admin/v1/index.proto", fileDescriptor_809cc5775e1cdeb3)
}

var fileDescriptor_809cc5775e1cdeb3 = []byte{
	// 529 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0xef, 0x8a, 0xd3, 0x4c,
	0x14, 0xc6, 0xdf, 0xa4, 0xcd, 0x6e, 0x73, 0xde, 0xdd, 0x35, 0x3b, 0xb0, 0x98, 0x16, 0xd7, 0x2d,
	0x05, 0xa5, 0x08, 0x26, 0x76, 0xfd, 0x28, 0x08, 0x69, 0x32, 0x2e, 0x91, 0x36, 0x2d, 0x49, 0x57,
	0xd0, 0x2f, 0x65, 0xb6, 0x33, 0xc6, 0xd1, 0x36, 0x53, 0xd2, 0x58, 0xdc, 0xab, 0xf0, 0x3e, 0xfc,
	0xe0, 0x35, 0x79, 0x29, 0x92, 0x99, 0xfe, 0xa3, 0xb8, 0xf4, 0x4b, 0x48, 0xce, 0x79, 0xce, 0xef,
	0x39, 0x79, 0x32, 0x81, 0x67, 0xa9, 0x10, 0xe9, 0x94, 0xb9, 0x94, 0x14, 0x64, 0x51, 0x88, 0x9c,
	0xb9, 0x84, 0xce, 0x78, 0xe6, 0x2e, 0x3b, 0x2e, 0xcf, 0x28, 0xfb, 0xe1, 0xcc, 0x73, 0x51, 0x08,
	0x54, 0x57, 0x32, 0x67, 0x23, 0x73, 0xa4, 0xcc, 0x59, 0x76, 0x1a, 0x57, 0x2b, 0x02, 0x99, 0x73,
	0xf7, 0x33, 0x67, 0x53, 0x3a, 0xbe, 0x63, 0x5f, 0xc8, 0x92, 0x8b, 0x5c, 0xcd, 0x36, 0x9e, 0xec,
	0x08, 0x48, 0x96, 0x89, 0x82, 0x14, 0x5c, 0x64, 0x0b, 0xd5, 0x6d, 0xfd, 0x36, 0xc0, 0x08, 0x4b,
	0x27, 0xd4, 0x02, 0x98, 0xe7, 0xe2, 0x2b, 0x9b, 0x14, 0x63, 0x4e, 0x6d, 0xad, 0xa9, 0xb5, 0xcd,
	0x6e, 0xe5, 0x8f, 0x57, 0x89, 0xcd, 0x55, 0x39, 0xa4, 0xe8, 0x29, 0xd4, 0xe4, 0x5a, 0xa5, 0xa2,
	0xb2, 0x55, 0x1c, 0xcb, 0x62, 0x48, 0xd1, 0x63, 0xa8, 0x7e, 0xe3, 0x19, 0xb5, 0xab, 0xeb, 0x9e,
	0x1e, 0xcb, 0x02, 0x8a, 0xa0, 0x46, 0xb2, 0x09, 0x2b, 0x77, 0xb7, 0x8d, 0xa6, 0xd6, 0x3e, 0xbb,
	0x7e, 0xe9, 0x3c, 0xf8, 0x4e, 0x8e, 0x5c, 0xc8, 0xf1, 0x56, 0x03, 0x7d, 0x41, 0x99, 0x62, 0x6d,
	0x18, 0xe8, 0x56, 0x2e, 0x3b, 0x67, 0x79, 0xc1, 0xd9, 0xc2, 0x3e, 0x6a, 0x56, 0xda, 0xff, 0x5f,
	0xbf, 0x3a, 0x48, 0x94, 0x57, 0x46, 0x87, 0x6a, 0xf2, 0x5e, 0x41, 0x77, 0x40, 0xc8, 0x03, 0x63,
	0x51, 0x90, 0x82, 0xd9, 0xc7, 0x72, 0xc7, 0xe7, 0x07, 0x89, 0x49, 0xa9, 0x56, 0x21, 0xa8, 0xc9,
	0xc6, 0x3d, 0x3c, 0xda, 0xb3, 0x29, 0x53, 0xc9, 0xc8, 0x8c, 0x6d, 0x33, 0xd5, 0x63, 0x59, 0x40,
	0x7d, 0x30, 0x29, 0xcf, 0xd9, 0xa4, 0xfc, 0x20, 0xb6, 0x2e, 0x2d, 0x5f, 0x1c, 0xb4, 0x0c, 0xd6,
	0x13, 0x8a, 0xb4, 0x25, 0xb4, 0xde, 0xc3, 0xc9, 0x6e, 0x66, 0xe8, 0x12, 0xea, 0x5e, 0xe4, 0xe3,
	0x64, 0x34, 0x88, 0xc7, 0xfd, 0x41, 0x80, 0xc7, 0xb7, 0x51, 0x32, 0xc4, 0x7e, 0xf8, 0x2e, 0xc4,
	0x81, 0xf5, 0x1f, 0xaa, 0x41, 0x35, 0x1a, 0x44, 0xd8, 0xd2, 0xd0, 0x39, 0x9c, 0x7a, 0xbd, 0xde,
	0x78, 0x2d, 0x4e, 0x2c, 0xbd, 0x85, 0xc1, 0xdc, 0x18, 0xa1, 0x3a, 0x5c, 0x04, 0x61, 0x8c, 0xfd,
	0x51, 0x38, 0x88, 0xf6, 0x20, 0xa7, 0x60, 0x7a, 0x89, 0x8f, 0xa3, 0x20, 0x8c, 0x6e, 0x2c, 0x0d,
	0x9d, 0x01, 0x04, 0x78, 0xf3, 0xac, 0xb7, 0x86, 0x60, 0xc8, 0x88, 0xd0, 0x05, 0x9c, 0x27, 0x23,
	0x6f, 0xb4, 0xbf, 0xc3, 0x09, 0xd4, 0xfc, 0x18, 0x7b, 0x23, 0x35, 0x6d, 0x82, 0x11, 0x63, 0x2f,
	0xf8, 0x68, 0xe9, 0x65, 0x23, 0xc0, 0x3d, 0x2c, 0x1b, 0x95, 0xb2, 0x81, 0xe3, 0x78, 0x10, 0x5b,
	0xd5, 0xee, 0x4f, 0x0d, 0x2e, 0x27, 0x62, 0xf6, 0x70, 0x4c, 0x5d, 0x90, 0x39, 0x0d, 0xcb, 0xe3,
	0x3d, 0xd4, 0x3e, 0xbd, 0x5d, 0x09, 0x53, 0x31, 0x25, 0x59, 0xea, 0x88, 0x3c, 0x75, 0x53, 0x96,
	0xc9, 0xc3, 0xef, 0xaa, 0x16, 0x99, 0xf3, 0xc5, 0x3f, 0x7e, 0xc0, 0x37, 0xf2, 0xe6, 0x97, 0x7e,
	0x75, 0xa3, 0x00, 0xfe, 0x54, 0x7c, 0xa7, 0x4e, 0xb0, 0xf1, 0xf3, 0xa4, 0xdf, 0x87, 0xce, 0xdd,
	0x91, 0x84, 0xbd, 0xfe, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x6f, 0xb4, 0x08, 0x5d, 0xcc, 0x03, 0x00,
	0x00,
}