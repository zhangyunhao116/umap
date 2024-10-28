package umap

import "math/bits"

func _wymix(a, b uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return hi ^ lo
}

func hashUint64(x uint64) uint64 {
	return _wymix(x^0xfffc2d0600147fc8, 0xf6ea71d5ec8a2980)
}
