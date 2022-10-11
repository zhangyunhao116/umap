package umap

import (
	"os"
	"strconv"
	"testing"
)

var benchRuntimeMap bool

func init() {
	if os.Getenv("BENCH_TYPE") == "runtime" {
		benchRuntimeMap = true
	}
}

var (
	cases = []int{6, 12, 18, 24, 30,
		64,
		128,
		256,
		512,
		1024,
		2048,
		4096,
		8192,
		1 << 16}
)

func BenchmarkMapAccessHit(b *testing.B) {
	if benchRuntimeMap {
		b.Run("Uint64", runWith(benchmarkMapAccessHitUint64Runtime, cases...))
	} else {
		b.Run("Uint64", runWith(benchmarkMapAccessHitUint64, cases...))
	}

}

func BenchmarkMapRange(b *testing.B) {
	if benchRuntimeMap {
		b.Run("Uint64", runWith(benchmarkMapRangeUint64Runtime, cases...))
	} else {
		b.Run("Uint64", runWith(benchmarkMapRangeUint64, cases...))
	}
}

func BenchmarkMapAssignGrow(b *testing.B) {
	if benchRuntimeMap {
		b.Run("Uint64", runWith(benchmarkMapAssignGrowUint64Runtime, cases...))
	} else {
		b.Run("Uint64", runWith(benchmarkMapAssignGrowUint64, cases...))
	}
}

func BenchmarkMapAssignReuse(b *testing.B) {
	if benchRuntimeMap {
		b.Run("Uint64", runWith(benchmarkMapAssignReuseUint64Runtime, cases...))
	} else {
		b.Run("Uint64", runWith(benchmarkMapAssignReuseUint64, cases...))
	}
}

func benchmarkMapAccessHitUint64(b *testing.B, n int) {
	m := New64(0)
	for i := 0; i < n; i++ {
		m.Store(uint64(i), uint64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Load(uint64(i & (n - 1)))
	}
}

func benchmarkMapAssignGrowUint64(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		m := New64(0)
		for j := uint64(0); int(j) < n; j++ {
			m.Store(j, j)
		}
	}
}

var rangecount uint64

func benchmarkMapRangeUint64(b *testing.B, n int) {
	m := New64(0)
	for i := 0; i < n; i++ {
		m.Store(uint64(i), uint64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Range(func(key, value uint64) bool {
			rangecount += key
			rangecount += value
			return true
		})
	}
}

func benchmarkMapAssignReuseUint64(b *testing.B, n int) {
	m := New64(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := uint64(0); int(j) < n; j++ {
			m.Store(j, j)
		}
		m.Range(func(key, _ uint64) bool {
			m.Delete(key)
			return true
		})
	}
}

func benchmarkMapAccessHitUint64Runtime(b *testing.B, n int) {
	m := make(map[uint64]uint64, 0)
	for i := 0; i < n; i++ {
		m[uint64(i)] = uint64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[uint64(i&(n-1))]
	}
}

func benchmarkMapAssignGrowUint64Runtime(b *testing.B, n int) {
	for i := 0; i < b.N; i++ {
		m := make(map[uint64]uint64, 0)
		for j := uint64(0); int(j) < n; j++ {
			m[uint64(j)] = uint64(j)
		}
	}
}

func benchmarkMapRangeUint64Runtime(b *testing.B, n int) {
	m := make(map[uint64]uint64, 0)
	for i := 0; i < n; i++ {
		m[uint64(i)] = uint64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for key, value := range m {
			rangecount += key
			rangecount += value
		}
	}
}

func benchmarkMapAssignReuseUint64Runtime(b *testing.B, n int) {
	m := make(map[uint64]uint64, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := uint64(0); int(j) < n; j++ {
			m[j] = j
		}
		for key, _ := range m {
			delete(m, key)
		}
	}
}

func runWith(f func(*testing.B, int), v ...int) func(*testing.B) {
	return func(b *testing.B) {
		for _, n := range v {
			b.Run(strconv.Itoa(n), func(b *testing.B) { f(b, n) })
		}
	}
}
