// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: composable/txboundary/v1beta1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryDelegateBoundaryRequest is the request type for the
// Query/DelegateBoundary RPC method.
type QueryDelegateBoundaryRequest struct {
}

func (m *QueryDelegateBoundaryRequest) Reset()         { *m = QueryDelegateBoundaryRequest{} }
func (m *QueryDelegateBoundaryRequest) String() string { return proto.CompactTextString(m) }
func (*QueryDelegateBoundaryRequest) ProtoMessage()    {}
func (*QueryDelegateBoundaryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_822e9740fb1ae211, []int{0}
}
func (m *QueryDelegateBoundaryRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryDelegateBoundaryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryDelegateBoundaryRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryDelegateBoundaryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryDelegateBoundaryRequest.Merge(m, src)
}
func (m *QueryDelegateBoundaryRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryDelegateBoundaryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryDelegateBoundaryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryDelegateBoundaryRequest proto.InternalMessageInfo

// QueryDelegateBoundaryResponse is the response type for the
// Query/DelegateBoundary RPC method.
type QueryDelegateBoundaryResponse struct {
	// boundary defines the boundary for the delegate tx
	Boundary Boundary `protobuf:"bytes,1,opt,name=boundary,proto3" json:"boundary"`
}

func (m *QueryDelegateBoundaryResponse) Reset()         { *m = QueryDelegateBoundaryResponse{} }
func (m *QueryDelegateBoundaryResponse) String() string { return proto.CompactTextString(m) }
func (*QueryDelegateBoundaryResponse) ProtoMessage()    {}
func (*QueryDelegateBoundaryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_822e9740fb1ae211, []int{1}
}
func (m *QueryDelegateBoundaryResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryDelegateBoundaryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryDelegateBoundaryResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryDelegateBoundaryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryDelegateBoundaryResponse.Merge(m, src)
}
func (m *QueryDelegateBoundaryResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryDelegateBoundaryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryDelegateBoundaryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryDelegateBoundaryResponse proto.InternalMessageInfo

func (m *QueryDelegateBoundaryResponse) GetBoundary() Boundary {
	if m != nil {
		return m.Boundary
	}
	return Boundary{}
}

// QueryRedelegateBoundaryRequest is the request type for the
// Query/ReDelegateBoundary RPC method.
type QueryRedelegateBoundaryRequest struct {
}

func (m *QueryRedelegateBoundaryRequest) Reset()         { *m = QueryRedelegateBoundaryRequest{} }
func (m *QueryRedelegateBoundaryRequest) String() string { return proto.CompactTextString(m) }
func (*QueryRedelegateBoundaryRequest) ProtoMessage()    {}
func (*QueryRedelegateBoundaryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_822e9740fb1ae211, []int{2}
}
func (m *QueryRedelegateBoundaryRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryRedelegateBoundaryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryRedelegateBoundaryRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryRedelegateBoundaryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryRedelegateBoundaryRequest.Merge(m, src)
}
func (m *QueryRedelegateBoundaryRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryRedelegateBoundaryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryRedelegateBoundaryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryRedelegateBoundaryRequest proto.InternalMessageInfo

// QueryRedelegateBoundaryResponse is the response type for the
// Query/ReDelegateBoundary RPC method.
type QueryRedelegateBoundaryResponse struct {
	// boundary defines the boundary for the redelegate tx
	Boundary Boundary `protobuf:"bytes,1,opt,name=boundary,proto3" json:"boundary"`
}

func (m *QueryRedelegateBoundaryResponse) Reset()         { *m = QueryRedelegateBoundaryResponse{} }
func (m *QueryRedelegateBoundaryResponse) String() string { return proto.CompactTextString(m) }
func (*QueryRedelegateBoundaryResponse) ProtoMessage()    {}
func (*QueryRedelegateBoundaryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_822e9740fb1ae211, []int{3}
}
func (m *QueryRedelegateBoundaryResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryRedelegateBoundaryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryRedelegateBoundaryResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryRedelegateBoundaryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryRedelegateBoundaryResponse.Merge(m, src)
}
func (m *QueryRedelegateBoundaryResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryRedelegateBoundaryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryRedelegateBoundaryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryRedelegateBoundaryResponse proto.InternalMessageInfo

func (m *QueryRedelegateBoundaryResponse) GetBoundary() Boundary {
	if m != nil {
		return m.Boundary
	}
	return Boundary{}
}

func init() {
	proto.RegisterType((*QueryDelegateBoundaryRequest)(nil), "composable.txboundary.v1beta1.QueryDelegateBoundaryRequest")
	proto.RegisterType((*QueryDelegateBoundaryResponse)(nil), "composable.txboundary.v1beta1.QueryDelegateBoundaryResponse")
	proto.RegisterType((*QueryRedelegateBoundaryRequest)(nil), "composable.txboundary.v1beta1.QueryRedelegateBoundaryRequest")
	proto.RegisterType((*QueryRedelegateBoundaryResponse)(nil), "composable.txboundary.v1beta1.QueryRedelegateBoundaryResponse")
}

func init() {
	proto.RegisterFile("composable/txboundary/v1beta1/query.proto", fileDescriptor_822e9740fb1ae211)
}

var fileDescriptor_822e9740fb1ae211 = []byte{
	// 344 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x52, 0x3d, 0x4b, 0xfb, 0x40,
	0x18, 0xcf, 0xfd, 0x5f, 0x44, 0xce, 0x45, 0x4e, 0x07, 0x09, 0xed, 0xb5, 0x64, 0x51, 0xd1, 0xe4,
	0x68, 0x8b, 0x93, 0x2f, 0x43, 0x71, 0x71, 0xb4, 0xa3, 0xdb, 0xa5, 0x7d, 0x08, 0x95, 0x34, 0x4f,
	0x9a, 0xbb, 0x4a, 0xbb, 0xfa, 0x09, 0x04, 0x27, 0xbf, 0x8e, 0x53, 0x71, 0x2a, 0xb8, 0x38, 0x89,
	0xb4, 0x7e, 0x10, 0x69, 0x9a, 0xb4, 0x50, 0x4d, 0x94, 0xe2, 0x16, 0x78, 0x7e, 0xef, 0x17, 0xba,
	0xdf, 0xc4, 0x4e, 0x88, 0x4a, 0xba, 0x3e, 0x08, 0xdd, 0x77, 0xb1, 0x17, 0xb4, 0x64, 0x34, 0x10,
	0x37, 0x15, 0x17, 0xb4, 0xac, 0x88, 0x6e, 0x0f, 0xa2, 0x81, 0x13, 0x46, 0xa8, 0x91, 0x15, 0x17,
	0x50, 0x67, 0x01, 0x75, 0x12, 0xa8, 0xb9, 0xed, 0xa1, 0x87, 0x31, 0x52, 0x4c, 0xbf, 0x66, 0x24,
	0xb3, 0xe0, 0x21, 0x7a, 0x3e, 0x08, 0x19, 0xb6, 0x85, 0x0c, 0x02, 0xd4, 0x52, 0xb7, 0x31, 0x50,
	0xc9, 0xf5, 0x30, 0xdf, 0x7d, 0xee, 0x11, 0xa3, 0x2d, 0x4e, 0x0b, 0x97, 0xd3, 0x3c, 0xe7, 0xe0,
	0x83, 0x27, 0x35, 0xd4, 0x93, 0x73, 0x03, 0xba, 0x3d, 0x50, 0xda, 0xba, 0xa6, 0xc5, 0x8c, 0xbb,
	0x0a, 0x31, 0x50, 0xc0, 0x2e, 0xe8, 0x7a, 0x2a, 0xb9, 0x43, 0xca, 0x64, 0x6f, 0xa3, 0xba, 0xeb,
	0xe4, 0x96, 0x72, 0x52, 0x89, 0xfa, 0xbf, 0xe1, 0x6b, 0xc9, 0x68, 0xcc, 0xe9, 0x56, 0x99, 0xf2,
	0xd8, 0xab, 0x01, 0xad, 0x8c, 0x34, 0x3e, 0x2d, 0x65, 0x22, 0x7e, 0x3d, 0x4f, 0xf5, 0xe1, 0x2f,
	0xfd, 0x1f, 0xdb, 0xb1, 0x47, 0x42, 0x37, 0x97, 0x17, 0x60, 0xc7, 0xdf, 0xe8, 0xe6, 0xed, 0x6a,
	0x9e, 0xac, 0x46, 0x9e, 0x95, 0xb4, 0x6a, 0xb7, 0xcf, 0xef, 0xf7, 0x7f, 0x6c, 0x76, 0x20, 0x9a,
	0xa8, 0x3a, 0xa8, 0xbe, 0x7a, 0xe8, 0x74, 0xa1, 0xf4, 0xc0, 0x9e, 0x08, 0x65, 0x9f, 0x87, 0x63,
	0xa7, 0x3f, 0x49, 0x92, 0xf9, 0x24, 0xe6, 0xd9, 0xaa, 0xf4, 0xa4, 0xca, 0x51, 0x5c, 0x45, 0x30,
	0x3b, 0xa7, 0x4a, 0x04, 0xcb, 0x65, 0xea, 0xf6, 0x70, 0xcc, 0xc9, 0x68, 0xcc, 0xc9, 0xdb, 0x98,
	0x93, 0xbb, 0x09, 0x37, 0x46, 0x13, 0x6e, 0xbc, 0x4c, 0xb8, 0x71, 0xb5, 0xd5, 0x17, 0xba, 0x6f,
	0xcf, 0x35, 0xf4, 0x20, 0x04, 0xe5, 0xae, 0xc5, 0x7f, 0x7b, 0xed, 0x23, 0x00, 0x00, 0xff, 0xff,
	0x6d, 0x4e, 0xd5, 0x97, 0x9b, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// DelegateBoundary returns the  boundary for the delegate tx.
	DelegateBoundary(ctx context.Context, in *QueryDelegateBoundaryRequest, opts ...grpc.CallOption) (*QueryDelegateBoundaryResponse, error)
	// RedelegateBoundary returns the  boundary for the redelegate tx.
	RedelegateBoundary(ctx context.Context, in *QueryRedelegateBoundaryRequest, opts ...grpc.CallOption) (*QueryRedelegateBoundaryResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) DelegateBoundary(ctx context.Context, in *QueryDelegateBoundaryRequest, opts ...grpc.CallOption) (*QueryDelegateBoundaryResponse, error) {
	out := new(QueryDelegateBoundaryResponse)
	err := c.cc.Invoke(ctx, "/composable.txboundary.v1beta1.Query/DelegateBoundary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) RedelegateBoundary(ctx context.Context, in *QueryRedelegateBoundaryRequest, opts ...grpc.CallOption) (*QueryRedelegateBoundaryResponse, error) {
	out := new(QueryRedelegateBoundaryResponse)
	err := c.cc.Invoke(ctx, "/composable.txboundary.v1beta1.Query/RedelegateBoundary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// DelegateBoundary returns the  boundary for the delegate tx.
	DelegateBoundary(context.Context, *QueryDelegateBoundaryRequest) (*QueryDelegateBoundaryResponse, error)
	// RedelegateBoundary returns the  boundary for the redelegate tx.
	RedelegateBoundary(context.Context, *QueryRedelegateBoundaryRequest) (*QueryRedelegateBoundaryResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) DelegateBoundary(ctx context.Context, req *QueryDelegateBoundaryRequest) (*QueryDelegateBoundaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelegateBoundary not implemented")
}
func (*UnimplementedQueryServer) RedelegateBoundary(ctx context.Context, req *QueryRedelegateBoundaryRequest) (*QueryRedelegateBoundaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RedelegateBoundary not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_DelegateBoundary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryDelegateBoundaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DelegateBoundary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/composable.txboundary.v1beta1.Query/DelegateBoundary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DelegateBoundary(ctx, req.(*QueryDelegateBoundaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_RedelegateBoundary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRedelegateBoundaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).RedelegateBoundary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/composable.txboundary.v1beta1.Query/RedelegateBoundary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).RedelegateBoundary(ctx, req.(*QueryRedelegateBoundaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "composable.txboundary.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DelegateBoundary",
			Handler:    _Query_DelegateBoundary_Handler,
		},
		{
			MethodName: "RedelegateBoundary",
			Handler:    _Query_RedelegateBoundary_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "composable/txboundary/v1beta1/query.proto",
}

func (m *QueryDelegateBoundaryRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryDelegateBoundaryRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryDelegateBoundaryRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryDelegateBoundaryResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryDelegateBoundaryResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryDelegateBoundaryResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Boundary.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryRedelegateBoundaryRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryRedelegateBoundaryRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryRedelegateBoundaryRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryRedelegateBoundaryResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryRedelegateBoundaryResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryRedelegateBoundaryResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Boundary.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryDelegateBoundaryRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryDelegateBoundaryResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Boundary.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryRedelegateBoundaryRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryRedelegateBoundaryResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Boundary.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryDelegateBoundaryRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryDelegateBoundaryRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryDelegateBoundaryRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryDelegateBoundaryResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryDelegateBoundaryResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryDelegateBoundaryResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Boundary", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Boundary.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryRedelegateBoundaryRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryRedelegateBoundaryRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryRedelegateBoundaryRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryRedelegateBoundaryResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryRedelegateBoundaryResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryRedelegateBoundaryResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Boundary", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Boundary.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
