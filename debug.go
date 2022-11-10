package umap

import (
	"fmt"
	"unsafe"
)

func (h *Uint64Map) debug() {
	println("-------------MAP DEBUG-------------")
	fmt.Printf("%+v\n", *h)
	if h.oldbuckets != nil {
		println("**********OLD-BUCKET**********")
		debugbucket(h.oldbuckets, uint(h.noldbuckets()))
	}
	println("**********BUCKET**********")
	debugbucket(h.buckets, uint(h.bucketmask)+1)
	println("-------------MAP DEBUG END-------------")
}

func debugbucket(b unsafe.Pointer, length uint) {
	for i := uint(0); i < length; i++ {
		bmap := bmapPointer(b, i)
		fmt.Printf("---BUCKET %v---\n", i)
		fmt.Printf("%+v\n", bmap.tophash)
		for j := 0; j < bucketCnt; j++ {
			fmt.Printf("<%v,%v>", bmap.data[j].key, bmap.data[j].value)
		}
		println()
	}
}
