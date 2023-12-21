// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: centauri/stakingmiddleware/v1beta1/genesis.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
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

// GenesisState defines the stakingmiddleware module's genesis state.
type GenesisState struct {
	// last_total_power tracks the total amounts of bonded tokens recorded during
	// the previous end block.
	LastTotalPower github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,1,opt,name=last_total_power,json=lastTotalPower,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"last_total_power"`
	Params         Params                                 `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
	// delegations defines the delegations active at genesis.
	Delegations                []Delegation                `protobuf:"bytes,3,rep,name=delegations,proto3" json:"delegations"`
	Begindelegations           []BeginRedelegate           `protobuf:"bytes,4,rep,name=begindelegations,proto3" json:"begindelegations"`
	Undelegates                []Undelegate                `protobuf:"bytes,5,rep,name=undelegates,proto3" json:"undelegates"`
	Cancelunbondingdelegations []CancelUnbondingDelegation `protobuf:"bytes,6,rep,name=cancelunbondingdelegations,proto3" json:"cancelunbondingdelegations"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_1aa0bd912277c095, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetDelegations() []Delegation {
	if m != nil {
		return m.Delegations
	}
	return nil
}

func (m *GenesisState) GetBegindelegations() []BeginRedelegate {
	if m != nil {
		return m.Begindelegations
	}
	return nil
}

func (m *GenesisState) GetUndelegates() []Undelegate {
	if m != nil {
		return m.Undelegates
	}
	return nil
}

func (m *GenesisState) GetCancelunbondingdelegations() []CancelUnbondingDelegation {
	if m != nil {
		return m.Cancelunbondingdelegations
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "centauri.stakingmiddleware.v1beta1.GenesisState")
}

func init() {
	proto.RegisterFile("centauri/stakingmiddleware/v1beta1/genesis.proto", fileDescriptor_1aa0bd912277c095)
}

var fileDescriptor_1aa0bd912277c095 = []byte{
	// 404 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xc1, 0x8a, 0xd3, 0x40,
	0x18, 0x80, 0x33, 0x6e, 0x2d, 0x98, 0x2e, 0xb2, 0x06, 0x0f, 0xb1, 0x87, 0x6c, 0xd9, 0x83, 0x94,
	0x05, 0x27, 0xb6, 0xbd, 0x09, 0x5e, 0xa2, 0xa0, 0xde, 0x4a, 0xb4, 0x20, 0x7a, 0x28, 0x93, 0xe4,
	0x67, 0x1c, 0x9b, 0xcc, 0x84, 0xcc, 0xc4, 0xea, 0x1b, 0x78, 0xf4, 0x31, 0x3c, 0xfa, 0x18, 0x3d,
	0xf6, 0x28, 0x1e, 0x8a, 0x34, 0x07, 0x5f, 0x43, 0x66, 0x92, 0x4a, 0xa0, 0xc5, 0xcd, 0xa5, 0x0d,
	0x43, 0xbe, 0xef, 0xcb, 0x3f, 0xfc, 0xf6, 0xe3, 0x18, 0xb8, 0x22, 0x65, 0xc1, 0x7c, 0xa9, 0xc8,
	0x8a, 0x71, 0x9a, 0xb1, 0x24, 0x49, 0x61, 0x4d, 0x0a, 0xf0, 0x3f, 0x4d, 0x22, 0x50, 0x64, 0xe2,
	0x53, 0xe0, 0x20, 0x99, 0xc4, 0x79, 0x21, 0x94, 0x70, 0xae, 0x0e, 0x04, 0x3e, 0x22, 0x70, 0x43,
	0x0c, 0xef, 0x53, 0x41, 0x85, 0x79, 0xdd, 0xd7, 0x4f, 0x35, 0x39, 0x7c, 0xd2, 0xa1, 0x75, 0xec,
	0xac, 0xd9, 0x7b, 0x24, 0x63, 0x5c, 0xf8, 0xe6, 0xb7, 0x3e, 0xba, 0xaa, 0x7a, 0xf6, 0xf9, 0x8b,
	0xfa, 0xd3, 0x5e, 0x2b, 0xa2, 0xc0, 0x79, 0x6b, 0x5f, 0xa4, 0x44, 0xaa, 0xa5, 0x12, 0x8a, 0xa4,
	0xcb, 0x5c, 0xac, 0xa1, 0x70, 0xd1, 0x08, 0x8d, 0xcf, 0x03, 0xbc, 0xd9, 0x5d, 0x5a, 0xbf, 0x76,
	0x97, 0x0f, 0x29, 0x53, 0x1f, 0xca, 0x08, 0xc7, 0x22, 0xf3, 0x63, 0x21, 0x33, 0x21, 0x9b, 0xbf,
	0x47, 0x32, 0x59, 0xf9, 0xea, 0x4b, 0x0e, 0x12, 0xbf, 0xe2, 0x2a, 0xbc, 0xab, 0x3d, 0x6f, 0xb4,
	0x66, 0xae, 0x2d, 0xce, 0x4b, 0xbb, 0x9f, 0x93, 0x82, 0x64, 0xd2, 0xbd, 0x35, 0x42, 0xe3, 0xc1,
	0xf4, 0x1a, 0xdf, 0x7c, 0x09, 0x78, 0x6e, 0x88, 0xa0, 0xa7, 0xdb, 0x61, 0xc3, 0x3b, 0xef, 0xed,
	0x41, 0x02, 0x29, 0x50, 0xa2, 0x98, 0xe0, 0xd2, 0x3d, 0x1b, 0x9d, 0x8d, 0x07, 0x53, 0xdc, 0x45,
	0xf7, 0xfc, 0x1f, 0x16, 0xdc, 0xd1, 0xca, 0xef, 0x7f, 0x7e, 0x5c, 0xa3, 0xb0, 0x6d, 0x73, 0x3e,
	0xda, 0x17, 0x11, 0x50, 0xc6, 0xdb, 0x85, 0x9e, 0x29, 0xcc, 0xba, 0x14, 0x02, 0xcd, 0x86, 0xd0,
	0xd0, 0xd0, 0xce, 0x1c, 0x79, 0xf5, 0x20, 0xe5, 0xe1, 0x00, 0xa4, 0x7b, 0xbb, 0xfb, 0x20, 0x0b,
	0x7e, 0xaa, 0xd0, 0xb6, 0x39, 0x5f, 0x91, 0x3d, 0x8c, 0x09, 0x8f, 0x21, 0x2d, 0x79, 0x24, 0x78,
	0xc2, 0x38, 0x6d, 0xcf, 0xd4, 0x37, 0xb1, 0xa7, 0x5d, 0x62, 0xcf, 0x8c, 0x65, 0x71, 0xb0, 0x9c,
	0xbe, 0xc4, 0xff, 0xb4, 0x82, 0xd9, 0x66, 0xef, 0xa1, 0xed, 0xde, 0x43, 0xbf, 0xf7, 0x1e, 0xfa,
	0x56, 0x79, 0xd6, 0xb6, 0xf2, 0xac, 0x9f, 0x95, 0x67, 0xbd, 0x7b, 0xf0, 0xf9, 0xc4, 0x1e, 0x9b,
	0x1d, 0x8a, 0xfa, 0x66, 0x43, 0x67, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x66, 0xe2, 0x9b, 0xf3,
	0x5e, 0x03, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Cancelunbondingdelegations) > 0 {
		for iNdEx := len(m.Cancelunbondingdelegations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Cancelunbondingdelegations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.Undelegates) > 0 {
		for iNdEx := len(m.Undelegates) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Undelegates[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Begindelegations) > 0 {
		for iNdEx := len(m.Begindelegations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Begindelegations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Delegations) > 0 {
		for iNdEx := len(m.Delegations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Delegations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.LastTotalPower.Size()
		i -= size
		if _, err := m.LastTotalPower.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.LastTotalPower.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.Delegations) > 0 {
		for _, e := range m.Delegations {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Begindelegations) > 0 {
		for _, e := range m.Begindelegations {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Undelegates) > 0 {
		for _, e := range m.Undelegates {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Cancelunbondingdelegations) > 0 {
		for _, e := range m.Cancelunbondingdelegations {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastTotalPower", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LastTotalPower.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Delegations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Delegations = append(m.Delegations, Delegation{})
			if err := m.Delegations[len(m.Delegations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Begindelegations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Begindelegations = append(m.Begindelegations, BeginRedelegate{})
			if err := m.Begindelegations[len(m.Begindelegations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Undelegates", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Undelegates = append(m.Undelegates, Undelegate{})
			if err := m.Undelegates[len(m.Undelegates)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cancelunbondingdelegations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Cancelunbondingdelegations = append(m.Cancelunbondingdelegations, CancelUnbondingDelegation{})
			if err := m.Cancelunbondingdelegations[len(m.Cancelunbondingdelegations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
