// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
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

// The request message containing the user's name.
type Face struct {
	Descriptor_          []float32 `protobuf:"fixed32,1,rep,packed,name=descriptor,proto3" json:"descriptor,omitempty"`
	Jumin                string    `protobuf:"bytes,2,opt,name=jumin,proto3" json:"jumin,omitempty"`
	Name                 string    `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Index                int32     `protobuf:"varint,4,opt,name=index,proto3" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Face) Reset()         { *m = Face{} }
func (m *Face) String() string { return proto.CompactTextString(m) }
func (*Face) ProtoMessage()    {}
func (*Face) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{0}
}

func (m *Face) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Face.Unmarshal(m, b)
}
func (m *Face) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Face.Marshal(b, m, deterministic)
}
func (m *Face) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Face.Merge(m, src)
}
func (m *Face) XXX_Size() int {
	return xxx_messageInfo_Face.Size(m)
}
func (m *Face) XXX_DiscardUnknown() {
	xxx_messageInfo_Face.DiscardUnknown(m)
}

var xxx_messageInfo_Face proto.InternalMessageInfo

func (m *Face) GetDescriptor_() []float32 {
	if m != nil {
		return m.Descriptor_
	}
	return nil
}

func (m *Face) GetJumin() string {
	if m != nil {
		return m.Jumin
	}
	return ""
}

func (m *Face) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Face) GetIndex() int32 {
	if m != nil {
		return m.Index
	}
	return 0
}

// The response message containing the greetings
type Res struct {
	Jumin                string   `protobuf:"bytes,1,opt,name=jumin,proto3" json:"jumin,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Res) Reset()         { *m = Res{} }
func (m *Res) String() string { return proto.CompactTextString(m) }
func (*Res) ProtoMessage()    {}
func (*Res) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{1}
}

func (m *Res) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Res.Unmarshal(m, b)
}
func (m *Res) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Res.Marshal(b, m, deterministic)
}
func (m *Res) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Res.Merge(m, src)
}
func (m *Res) XXX_Size() int {
	return xxx_messageInfo_Res.Size(m)
}
func (m *Res) XXX_DiscardUnknown() {
	xxx_messageInfo_Res.DiscardUnknown(m)
}

var xxx_messageInfo_Res proto.InternalMessageInfo

func (m *Res) GetJumin() string {
	if m != nil {
		return m.Jumin
	}
	return ""
}

func (m *Res) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Face)(nil), "proto.Face")
	proto.RegisterType((*Res)(nil), "proto.Res")
}

func init() { proto.RegisterFile("proto.proto", fileDescriptor_2fcc84b9998d60d8) }

var fileDescriptor_2fcc84b9998d60d8 = []byte{
	// 184 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8e, 0xc1, 0xca, 0xc2, 0x30,
	0x10, 0x84, 0x9b, 0x36, 0xf9, 0xe1, 0xdf, 0xde, 0x16, 0x0f, 0xc1, 0x83, 0x94, 0x1c, 0x24, 0xa7,
	0x2a, 0xfa, 0x0e, 0xde, 0x0d, 0xbe, 0x40, 0x4d, 0x57, 0xa9, 0xd0, 0xa6, 0x34, 0x15, 0x7c, 0x7c,
	0x69, 0x22, 0xb4, 0x82, 0x97, 0xdd, 0xfd, 0x86, 0xd9, 0x61, 0x20, 0xef, 0x07, 0x37, 0xba, 0x32,
	0x4c, 0x14, 0x61, 0xa9, 0x1b, 0xf0, 0x53, 0x65, 0x09, 0x37, 0x00, 0x35, 0x79, 0x3b, 0x34, 0xfd,
	0xe8, 0x06, 0xc9, 0x8a, 0x4c, 0xa7, 0x66, 0xa1, 0xe0, 0x0a, 0xc4, 0xe3, 0xd9, 0x36, 0x9d, 0x4c,
	0x0b, 0xa6, 0xff, 0x4d, 0x04, 0x44, 0xe0, 0x5d, 0xd5, 0x92, 0xcc, 0x82, 0x18, 0xee, 0xc9, 0xd9,
	0x74, 0x35, 0xbd, 0x24, 0x2f, 0x98, 0x16, 0x26, 0x82, 0xda, 0x41, 0x66, 0xc8, 0xcf, 0x31, 0xec,
	0x57, 0x4c, 0x3a, 0xc7, 0x1c, 0xce, 0xd3, 0x83, 0x45, 0x05, 0xc2, 0x90, 0x75, 0x77, 0xcc, 0x63,
	0xef, 0x72, 0x6a, 0xbb, 0x86, 0x0f, 0x18, 0xf2, 0x2a, 0xc1, 0x2d, 0xf0, 0x0b, 0xf9, 0xf1, 0xdb,
	0xb2, 0x04, 0x95, 0x68, 0xb6, 0x67, 0xd7, 0xbf, 0xa0, 0x1c, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff,
	0xc9, 0x53, 0x4e, 0xea, 0x08, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RecClient is the client API for Rec service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RecClient interface {
	// Sends a greeting
	Recog(ctx context.Context, in *Face, opts ...grpc.CallOption) (*Res, error)
	Test(ctx context.Context, opts ...grpc.CallOption) (Rec_TestClient, error)
}

type recClient struct {
	cc *grpc.ClientConn
}

func NewRecClient(cc *grpc.ClientConn) RecClient {
	return &recClient{cc}
}

func (c *recClient) Recog(ctx context.Context, in *Face, opts ...grpc.CallOption) (*Res, error) {
	out := new(Res)
	err := c.cc.Invoke(ctx, "/proto.Rec/Recog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recClient) Test(ctx context.Context, opts ...grpc.CallOption) (Rec_TestClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Rec_serviceDesc.Streams[0], "/proto.Rec/Test", opts...)
	if err != nil {
		return nil, err
	}
	x := &recTestClient{stream}
	return x, nil
}

type Rec_TestClient interface {
	Send(*Face) error
	Recv() (*Face, error)
	grpc.ClientStream
}

type recTestClient struct {
	grpc.ClientStream
}

func (x *recTestClient) Send(m *Face) error {
	return x.ClientStream.SendMsg(m)
}

func (x *recTestClient) Recv() (*Face, error) {
	m := new(Face)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RecServer is the server API for Rec service.
type RecServer interface {
	// Sends a greeting
	Recog(context.Context, *Face) (*Res, error)
	Test(Rec_TestServer) error
}

// UnimplementedRecServer can be embedded to have forward compatible implementations.
type UnimplementedRecServer struct {
}

func (*UnimplementedRecServer) Recog(ctx context.Context, req *Face) (*Res, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Recog not implemented")
}
func (*UnimplementedRecServer) Test(srv Rec_TestServer) error {
	return status.Errorf(codes.Unimplemented, "method Test not implemented")
}

func RegisterRecServer(s *grpc.Server, srv RecServer) {
	s.RegisterService(&_Rec_serviceDesc, srv)
}

func _Rec_Recog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Face)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecServer).Recog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Rec/Recog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecServer).Recog(ctx, req.(*Face))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rec_Test_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RecServer).Test(&recTestServer{stream})
}

type Rec_TestServer interface {
	Send(*Face) error
	Recv() (*Face, error)
	grpc.ServerStream
}

type recTestServer struct {
	grpc.ServerStream
}

func (x *recTestServer) Send(m *Face) error {
	return x.ServerStream.SendMsg(m)
}

func (x *recTestServer) Recv() (*Face, error) {
	m := new(Face)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Rec_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Rec",
	HandlerType: (*RecServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Recog",
			Handler:    _Rec_Recog_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Test",
			Handler:       _Rec_Test_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto.proto",
}
