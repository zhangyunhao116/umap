package umap

import (
	"testing"
)

func TestConst(t *testing.T) {
	if !matchEmpty(littleEndianUint64ToBytes(allEvacuatedWithEmpty)).AnyMatch() {
		t.Fatal()
	}
}

func TestBitHackEvacuate(t *testing.T) {
	s := [8]uint8{evacuatedSlot, evacuatedSlot, evacuatedSlot, evacuatedSlot, evacuatedSlot, evacuatedSlot, evacuatedSlot, evacuatedSlot}
	status := matchEmpty(s)
	if status.AnyMatch() {
		t.Fatal()
	}

	s = [8]uint8{evacuatedSlot, evacuatedSlot, evacuatedSlot, evacuatedSlot, evacuatedSlot, evacuatedSlot, evacuatedSlot, evacuatedSlot}
	status = matchEmptyOrDeleted(s)
	if status.AnyMatch() {
		t.Fatal()
	}

	s = [8]uint8{evacuatedSlot, evacuatedSlot, evacuatedSlot, evacuatedSlot, emptySlot, evacuatedSlot, evacuatedSlot, evacuatedSlot}
	status = matchEmpty(s)
	if !status.AnyMatch() {
		t.Fatal()
	}

	s = [8]uint8{evacuatedSlot, evacuatedSlot, deletedSlot, evacuatedSlot, emptySlot, evacuatedSlot, evacuatedSlot, evacuatedSlot}
	status = matchEmpty(s)
	if !status.AnyMatch() {
		t.Fatal()
	}

	s = [8]uint8{evacuatedSlot, evacuatedSlot, deletedSlot, evacuatedSlot, emptySlot, evacuatedSlot, evacuatedSlot, evacuatedSlot}
	status = matchEmptyOrDeleted(s)
	if !status.shouldHaveMatches(2, 4) {
		t.Fatal()
	}
}

func TestBitHackMatchTophash(t *testing.T) {
	s := [8]uint8{emptySlot, 12, deletedSlot, 1, 13, 14}
	status := matchTopHash(s, 12)
	if !status.shouldHaveMatch(1) {
		t.Fatal()
	}

	s = [8]uint8{11, 122, 0, deletedSlot, emptySlot, deletedSlot}
	status = matchTopHash(s, 12)
	if status.AnyMatch() {
		t.Fatal()
	}

	s = [8]uint8{1, 127, 123, 127, 0, deletedSlot, emptySlot, deletedSlot}
	status = matchTopHash(s, 127)
	if !status.shouldHaveMatches(1, 3) {
		t.Fatal()
	}
}

func TestBitHackMatchEmpty(t *testing.T) {
	s := [8]uint8{emptySlot, 12, deletedSlot}
	status := matchEmpty(s)
	if !status.shouldHaveMatch(0) {
		t.Fatal()
	}

	s = [8]uint8{11, 122, 0, deletedSlot, deletedSlot, deletedSlot, 0, 1}
	status = matchEmpty(s)
	if status.AnyMatch() {
		t.Fatal()
	}

	s = littleEndianUint64ToBytes(allEmpty)
	status = matchEmpty(s)
	if !status.AnyMatch() {
		t.Fatal()
	}

	s = [8]uint8{1, emptySlot, 123, 127, 0, deletedSlot, emptySlot, deletedSlot}
	status = matchEmpty(s)
	if !status.shouldHaveMatches(1, 6) {
		t.Fatal()
	}
}

func TestBitHackMatchEmprtyOrDeleted(t *testing.T) {
	s := [8]uint8{emptySlot, 12, 0, 1}
	status := matchEmptyOrDeleted(s)
	if !status.shouldHaveMatch(0) {
		t.Fatal()
	}

	s = [8]uint8{11, 122, 0, 1, 127, 55, 0, 1}
	status = matchEmptyOrDeleted(s)
	if status.AnyMatch() {
		t.Fatal()
	}

	s = [8]uint8{1, emptySlot, 123, 127, 0, deletedSlot, emptySlot, deletedSlot}
	status = matchEmptyOrDeleted(s)
	if !status.shouldHaveMatches(1, 5, 6, 7) {
		t.Fatal()
	}
}

func TestBitHackSameSizeGrow(t *testing.T) {
	s := [bucketCnt]uint8{deletedSlot, 2, emptySlot, 4, deletedSlot, deletedSlot, 7, emptySlot}
	ns := prepareSameSizeGrow(s)
	res := [bucketCnt]uint8{emptySlot, deletedSlot, emptySlot, deletedSlot, emptySlot, emptySlot, deletedSlot, emptySlot}
	for i := 0; i < bucketCnt; i++ {
		if ns[i] != res[i] {
			t.Fatalf("Expected ns[%d] == res[%d], got %d != %d", i, i, ns[i], res[i])
		}
	}
}

func (b *bitmask64) shouldHaveMatch(x uint) bool {
	cb := *b
	for {
		m := cb.NextMatch()
		if m == x {
			return true
		}
		if m >= bucketCnt {
			return false
		}
		cb.RemoveLowestBit()
	}
}

func (b *bitmask64) shouldHaveMatches(x ...uint) bool {
	for _, v := range x {
		if !b.shouldHaveMatch(v) {
			return false
		}
	}
	return true
}
