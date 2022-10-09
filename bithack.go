package umap

import (
	"math/bits"
	"unsafe"
)

// This implementation assumes that
// - The deletedSlot must be 0b1000_0000.
// - The emptySlot must be 0b1111_1111.

const (
	bucketCnt = 8
)

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

func makeUint64BucketArray(size int) unsafe.Pointer {
	x := make([]bmapuint64, size)
	for i := range x {
		// The compiler will optimize this pattern.
		x[i].tophash[0] = emptySlot
		x[i].tophash[1] = emptySlot
		x[i].tophash[2] = emptySlot
		x[i].tophash[3] = emptySlot
		x[i].tophash[4] = emptySlot
		x[i].tophash[5] = emptySlot
		x[i].tophash[6] = emptySlot
		x[i].tophash[7] = emptySlot
	}
	return (*sliceHeader)(unsafe.Pointer(&x)).Data
}

func (b *bmapuint64) MatchEmptyOrDeleted() bitmask64 {
	// The high bit is set for both empty slot and deleted slot.
	tophashs := littleEndianBytesToUint64(b.tophash)
	return bitmask64(emptyOrDeletedMask & tophashs)
}

func (b *bmapuint64) MatchEmpty() bitmask64 {
	// Same as b.MatchTopHash(emptySlot), but faster.
	//
	// The high bit is set for both empty slot and deleted slot.
	// (tophashs & emptyOrDeletedMask) get all empty or deleted slots.
	// (tophashs << 1) clears the high bit for deletedSlot.
	// ANDing them we can get all the empty slots.
	tophashs := littleEndianBytesToUint64(b.tophash)
	return bitmask64((tophashs << 1) & tophashs & emptyOrDeletedMask)
}

func matchTopHash(tophash [bucketCnt]uint8, top uint8) bitmask64 {
	tophashs := littleEndianBytesToUint64(tophash)
	cmp := tophashs ^ (uint64(0x0101_0101_0101_0101) * uint64(top))
	return bitmask64((cmp - 0x0101_0101_0101_0101) & ^cmp & 0x8080_8080_8080_8080)
}

func matchFull(tophash [bucketCnt]uint8) bitmask64 {
	// If a slot is neither empty nor deleted, then it must be FUll.
	tophashs := littleEndianBytesToUint64(tophash)
	return bitmask64(emptyOrDeletedMask & ^tophashs)
}

func (b *bmapuint64) PrepareSameSizeGrow() {
	// Convert Deleted to Empty and Full to Deleted.
	tophashs := littleEndianBytesToUint64(b.tophash)
	full := ^tophashs & emptyOrDeletedMask
	full = ^full + (full >> 7)
	b.tophash = littleEndianUint64ToBytes(full)
}

func (b bitmask64) AnyMatch() bool {
	return b != 0
}

func (b *bitmask64) NextMatch() uint {
	return uint(bits.TrailingZeros64(uint64(*b)) / bucketCnt)
}

func (b *bitmask64) RemoveLowestBit() {
	*b = *b & (*b - 1)
}

func littleEndianBytesToUint64(v [8]uint8) uint64 {
	return uint64(v[0]) | uint64(v[1])<<8 | uint64(v[2])<<16 | uint64(v[3])<<24 | uint64(v[4])<<32 | uint64(v[5])<<40 | uint64(v[6])<<48 | uint64(v[7])<<56
}

func littleEndianUint64ToBytes(v uint64) [8]uint8 {
	return [8]uint8{uint8(v), uint8(v >> 8), uint8(v >> 16), uint8(v >> 24), uint8(v >> 32), uint8(v >> 40), uint8(v >> 48), uint8(v >> 56)}
}
