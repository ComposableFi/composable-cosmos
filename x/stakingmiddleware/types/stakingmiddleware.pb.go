// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: centauri/stakingmiddleware/v1beta1/stakingmiddleware.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// MsgDelegate defines a SDK message for performing a delegation of coins
// from a delegator to a validator.
type Delegation struct {
	DelegatorAddress string     `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	ValidatorAddress string     `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	Amount           types.Coin `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
}

func (m *Delegation) Reset()         { *m = Delegation{} }
func (m *Delegation) String() string { return proto.CompactTextString(m) }
func (*Delegation) ProtoMessage()    {}
func (*Delegation) Descriptor() ([]byte, []int) {
	return fileDescriptor_b3eb4aece69e048d, []int{0}
}
func (m *Delegation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Delegation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Delegation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Delegation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Delegation.Merge(m, src)
}
func (m *Delegation) XXX_Size() int {
	return m.Size()
}
func (m *Delegation) XXX_DiscardUnknown() {
	xxx_messageInfo_Delegation.DiscardUnknown(m)
}

var xxx_messageInfo_Delegation proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Delegation)(nil), "centauri.stakingmiddleware.v1beta1.Delegation")
}

func init() {
	proto.RegisterFile("centauri/stakingmiddleware/v1beta1/stakingmiddleware.proto", fileDescriptor_b3eb4aece69e048d)
}

var fileDescriptor_b3eb4aece69e048d = []byte{
	// 344 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xb2, 0x4a, 0x4e, 0xcd, 0x2b,
	0x49, 0x2c, 0x2d, 0xca, 0xd4, 0x2f, 0x2e, 0x49, 0xcc, 0xce, 0xcc, 0x4b, 0xcf, 0xcd, 0x4c, 0x49,
	0xc9, 0x49, 0x2d, 0x4f, 0x2c, 0x4a, 0xd5, 0x2f, 0x33, 0x4c, 0x4a, 0x2d, 0x49, 0x34, 0xc4, 0x94,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x52, 0x82, 0xe9, 0xd5, 0xc3, 0x54, 0x01, 0xd5, 0x2b,
	0x25, 0x92, 0x9e, 0x9f, 0x9e, 0x0f, 0x56, 0xae, 0x0f, 0x62, 0x41, 0x74, 0x4a, 0x49, 0x26, 0xe7,
	0x17, 0xe7, 0xe6, 0x17, 0xc7, 0x43, 0x24, 0x20, 0x1c, 0xa8, 0x94, 0x1c, 0x84, 0xa7, 0x9f, 0x94,
	0x58, 0x8c, 0x70, 0x41, 0x72, 0x7e, 0x66, 0x1e, 0x54, 0x5e, 0x30, 0x31, 0x37, 0x33, 0x2f, 0x5f,
	0x1f, 0x4c, 0x42, 0x85, 0xc4, 0xa1, 0x5a, 0x72, 0x8b, 0xd3, 0xf5, 0xcb, 0x0c, 0x41, 0x14, 0x44,
	0x42, 0x69, 0x32, 0x13, 0x17, 0x97, 0x4b, 0x6a, 0x4e, 0x6a, 0x7a, 0x62, 0x49, 0x66, 0x7e, 0x9e,
	0x90, 0x2b, 0x97, 0x60, 0x0a, 0x84, 0x97, 0x5f, 0x14, 0x9f, 0x98, 0x92, 0x52, 0x94, 0x5a, 0x5c,
	0x2c, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0xe9, 0x24, 0x71, 0x69, 0x8b, 0xae, 0x08, 0xd4, 0x1d, 0x8e,
	0x10, 0x99, 0xe0, 0x92, 0xa2, 0xcc, 0xbc, 0xf4, 0x20, 0x01, 0xb8, 0x16, 0xa8, 0x38, 0xc8, 0x98,
	0xb2, 0xc4, 0x9c, 0xcc, 0x14, 0x14, 0x63, 0x98, 0x08, 0x19, 0x03, 0xd7, 0x02, 0x33, 0xc6, 0x86,
	0x8b, 0x2d, 0x31, 0x37, 0xbf, 0x34, 0xaf, 0x44, 0x82, 0x59, 0x81, 0x51, 0x83, 0xdb, 0x48, 0x52,
	0x0f, 0xaa, 0x11, 0xe4, 0x73, 0x58, 0xf8, 0xe9, 0x39, 0xe7, 0x67, 0xe6, 0x39, 0x71, 0x9e, 0xb8,
	0x27, 0xcf, 0xb0, 0xe2, 0xf9, 0x06, 0x2d, 0xc6, 0x20, 0xa8, 0x1e, 0x2b, 0xcb, 0x8e, 0x05, 0xf2,
	0x0c, 0x2f, 0x16, 0xc8, 0x33, 0x34, 0x3d, 0xdf, 0xa0, 0x85, 0xe9, 0xad, 0xae, 0xe7, 0x1b, 0xb4,
	0xc4, 0x20, 0xe6, 0xe9, 0x16, 0xa7, 0x64, 0xeb, 0xfb, 0x16, 0xa7, 0x43, 0x03, 0x22, 0xd5, 0xc9,
	0xf8, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf0, 0x58,
	0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b, 0x8f, 0xe5, 0x18, 0xa2, 0x24, 0x2b, 0xb0, 0xa4, 0x82,
	0x92, 0xca, 0x82, 0xd4, 0xe2, 0x24, 0x36, 0x70, 0x88, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff,
	0xdc, 0x56, 0x33, 0x1d, 0x30, 0x02, 0x00, 0x00,
}

func (m *Delegation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Delegation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Delegation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Amount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintStakingmiddleware(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintStakingmiddleware(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DelegatorAddress) > 0 {
		i -= len(m.DelegatorAddress)
		copy(dAtA[i:], m.DelegatorAddress)
		i = encodeVarintStakingmiddleware(dAtA, i, uint64(len(m.DelegatorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintStakingmiddleware(dAtA []byte, offset int, v uint64) int {
	offset -= sovStakingmiddleware(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Delegation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DelegatorAddress)
	if l > 0 {
		n += 1 + l + sovStakingmiddleware(uint64(l))
	}
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovStakingmiddleware(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sovStakingmiddleware(uint64(l))
	return n
}

func sovStakingmiddleware(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozStakingmiddleware(x uint64) (n int) {
	return sovStakingmiddleware(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Delegation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStakingmiddleware
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
			return fmt.Errorf("proto: Delegation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Delegation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelegatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStakingmiddleware
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStakingmiddleware
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStakingmiddleware
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DelegatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStakingmiddleware
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStakingmiddleware
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStakingmiddleware
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStakingmiddleware
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
				return ErrInvalidLengthStakingmiddleware
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthStakingmiddleware
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStakingmiddleware(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthStakingmiddleware
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
func skipStakingmiddleware(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStakingmiddleware
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
					return 0, ErrIntOverflowStakingmiddleware
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
					return 0, ErrIntOverflowStakingmiddleware
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
				return 0, ErrInvalidLengthStakingmiddleware
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupStakingmiddleware
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthStakingmiddleware
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthStakingmiddleware        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStakingmiddleware          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupStakingmiddleware = fmt.Errorf("proto: unexpected end of group")
)