// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pairing/fixated_epoch_blocks_overlap.proto

package types

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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

type FixatedEpochBlocksOverlap struct {
	Index                string `protobuf:"bytes,1,opt,name=index,proto3" json:"index,omitempty"`
	ServicersToPairCount uint64 `protobuf:"varint,2,opt,name=servicersToPairCount,proto3" json:"servicersToPairCount,omitempty"`
	FixationBlock        uint64 `protobuf:"varint,3,opt,name=fixationBlock,proto3" json:"fixationBlock,omitempty"`
}

func (m *FixatedEpochBlocksOverlap) Reset()         { *m = FixatedEpochBlocksOverlap{} }
func (m *FixatedEpochBlocksOverlap) String() string { return proto.CompactTextString(m) }
func (*FixatedEpochBlocksOverlap) ProtoMessage()    {}
func (*FixatedEpochBlocksOverlap) Descriptor() ([]byte, []int) {
	return fileDescriptor_5cd9e38f2544d0ee, []int{0}
}
func (m *FixatedEpochBlocksOverlap) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FixatedEpochBlocksOverlap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FixatedEpochBlocksOverlap.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FixatedEpochBlocksOverlap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FixatedEpochBlocksOverlap.Merge(m, src)
}
func (m *FixatedEpochBlocksOverlap) XXX_Size() int {
	return m.Size()
}
func (m *FixatedEpochBlocksOverlap) XXX_DiscardUnknown() {
	xxx_messageInfo_FixatedEpochBlocksOverlap.DiscardUnknown(m)
}

var xxx_messageInfo_FixatedEpochBlocksOverlap proto.InternalMessageInfo

func (m *FixatedEpochBlocksOverlap) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *FixatedEpochBlocksOverlap) GetServicersToPairCount() uint64 {
	if m != nil {
		return m.ServicersToPairCount
	}
	return 0
}

func (m *FixatedEpochBlocksOverlap) GetFixationBlock() uint64 {
	if m != nil {
		return m.FixationBlock
	}
	return 0
}

func init() {
	proto.RegisterType((*FixatedEpochBlocksOverlap)(nil), "lavanet.lava.pairing.FixatedEpochBlocksOverlap")
}

func init() {
	proto.RegisterFile("pairing/fixated_epoch_blocks_overlap.proto", fileDescriptor_5cd9e38f2544d0ee)
}

var fileDescriptor_5cd9e38f2544d0ee = []byte{
	// 228 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2a, 0x48, 0xcc, 0x2c,
	0xca, 0xcc, 0x4b, 0xd7, 0x4f, 0xcb, 0xac, 0x48, 0x2c, 0x49, 0x4d, 0x89, 0x4f, 0x2d, 0xc8, 0x4f,
	0xce, 0x88, 0x4f, 0xca, 0xc9, 0x4f, 0xce, 0x2e, 0x8e, 0xcf, 0x2f, 0x4b, 0x2d, 0xca, 0x49, 0x2c,
	0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0xc9, 0x49, 0x2c, 0x4b, 0xcc, 0x4b, 0x2d, 0xd1,
	0x03, 0xd1, 0x7a, 0x50, 0x8d, 0x4a, 0xdd, 0x8c, 0x5c, 0x92, 0x6e, 0x10, 0xcd, 0xae, 0x20, 0xbd,
	0x4e, 0x60, 0xad, 0xfe, 0x10, 0x9d, 0x42, 0x22, 0x5c, 0xac, 0x99, 0x79, 0x29, 0xa9, 0x15, 0x12,
	0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x10, 0x8e, 0x90, 0x11, 0x97, 0x48, 0x71, 0x6a, 0x51, 0x59,
	0x66, 0x72, 0x6a, 0x51, 0x71, 0x48, 0x7e, 0x40, 0x62, 0x66, 0x91, 0x73, 0x7e, 0x69, 0x5e, 0x89,
	0x04, 0x93, 0x02, 0xa3, 0x06, 0x4b, 0x10, 0x56, 0x39, 0x21, 0x15, 0x2e, 0x5e, 0xb0, 0x1b, 0x33,
	0xf3, 0xf3, 0xc0, 0x56, 0x48, 0x30, 0x83, 0x15, 0xa3, 0x0a, 0x3a, 0x39, 0x9e, 0x78, 0x24, 0xc7,
	0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e, 0xcb, 0x31, 0x5c, 0x78, 0x2c,
	0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x7a, 0x7a, 0x66, 0x49, 0x46, 0x69, 0x92, 0x5e, 0x72,
	0x7e, 0xae, 0x3e, 0xd4, 0x23, 0x60, 0x5a, 0xbf, 0x42, 0x1f, 0x16, 0x06, 0x25, 0x95, 0x05, 0xa9,
	0xc5, 0x49, 0x6c, 0x60, 0xdf, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xfa, 0x11, 0xd4, 0x20,
	0x1b, 0x01, 0x00, 0x00,
}

func (m *FixatedEpochBlocksOverlap) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FixatedEpochBlocksOverlap) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FixatedEpochBlocksOverlap) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.FixationBlock != 0 {
		i = encodeVarintFixatedEpochBlocksOverlap(dAtA, i, uint64(m.FixationBlock))
		i--
		dAtA[i] = 0x18
	}
	if m.ServicersToPairCount != 0 {
		i = encodeVarintFixatedEpochBlocksOverlap(dAtA, i, uint64(m.ServicersToPairCount))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Index) > 0 {
		i -= len(m.Index)
		copy(dAtA[i:], m.Index)
		i = encodeVarintFixatedEpochBlocksOverlap(dAtA, i, uint64(len(m.Index)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintFixatedEpochBlocksOverlap(dAtA []byte, offset int, v uint64) int {
	offset -= sovFixatedEpochBlocksOverlap(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *FixatedEpochBlocksOverlap) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Index)
	if l > 0 {
		n += 1 + l + sovFixatedEpochBlocksOverlap(uint64(l))
	}
	if m.ServicersToPairCount != 0 {
		n += 1 + sovFixatedEpochBlocksOverlap(uint64(m.ServicersToPairCount))
	}
	if m.FixationBlock != 0 {
		n += 1 + sovFixatedEpochBlocksOverlap(uint64(m.FixationBlock))
	}
	return n
}

func sovFixatedEpochBlocksOverlap(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFixatedEpochBlocksOverlap(x uint64) (n int) {
	return sovFixatedEpochBlocksOverlap(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *FixatedEpochBlocksOverlap) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFixatedEpochBlocksOverlap
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
			return fmt.Errorf("proto: FixatedEpochBlocksOverlap: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FixatedEpochBlocksOverlap: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFixatedEpochBlocksOverlap
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
				return ErrInvalidLengthFixatedEpochBlocksOverlap
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFixatedEpochBlocksOverlap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Index = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ServicersToPairCount", wireType)
			}
			m.ServicersToPairCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFixatedEpochBlocksOverlap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ServicersToPairCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FixationBlock", wireType)
			}
			m.FixationBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFixatedEpochBlocksOverlap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FixationBlock |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFixatedEpochBlocksOverlap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFixatedEpochBlocksOverlap
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
func skipFixatedEpochBlocksOverlap(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFixatedEpochBlocksOverlap
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
					return 0, ErrIntOverflowFixatedEpochBlocksOverlap
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
					return 0, ErrIntOverflowFixatedEpochBlocksOverlap
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
				return 0, ErrInvalidLengthFixatedEpochBlocksOverlap
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupFixatedEpochBlocksOverlap
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthFixatedEpochBlocksOverlap
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthFixatedEpochBlocksOverlap        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFixatedEpochBlocksOverlap          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupFixatedEpochBlocksOverlap = fmt.Errorf("proto: unexpected end of group")
)
