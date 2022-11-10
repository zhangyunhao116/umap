package umap

import (
	"math/rand"
	"sort"
	"testing"
)

func TestQQQ(t *testing.T) {
	for i := 0; i < 100000; i++ {
		rand.Seed(int64(i))
		// println("SEED", i)
		testQQQ(t)
	}
	// rand.Seed(131)
	// testQQQ(t)
}

func testQQQ(t *testing.T) {
	const count = 100
	m1 := New64(0)
	m2 := make(map[uint64]uint64)
	for i := 0; i < count; i++ {
		choose := rand.Intn(100)
		rk, rv := uint64(rand.Int63n(count)), uint64(rand.Int63())
		if choose >= 50 {
			m1.Store(rk, rv)
			m2[rk] = rv
		} else if choose <= 10 {
			m1.Delete(rk)
			delete(m2, rk)
		} else if choose == 27 && rand.Intn(100) == 0 {
			// var (
			// 	tmp1 = make(map[uint64]uint64)
			// 	tmp2 = make(map[uint64]uint64)
			// )
			// m1.Range(func(key, value uint64) bool {
			// 	tmp1[key] = value
			// 	return true
			// })
			// for key, value := range m2 {
			// 	tmp2[key] = value
			// }
			// var x, y int
			// for range tmp1 {
			// 	x++
			// }
			// for range tmp2 {
			// 	y++
			// }
			// if x != y {
			// 	t.Fatal(x, y)
			// }
			// for k2, v2 := range tmp1 {
			// 	v1, ok := tmp2[k2]
			// 	if !ok || v1 != v2 {
			// 		t.Fatal("invalid key value", k2, v2, v1, ok)
			// 	}
			// }
			// for k2, v2 := range tmp2 {
			// 	v1, ok := tmp1[k2]
			// 	if !ok || v1 != v2 {
			// 		t.Fatal("invalid key value", k2, v2, v1, ok)
			// 	}
			// }
		} else {
			// Load
			v1, ok1 := m1.Load(rk)
			v2, ok2 := m2[rk]
			if ok1 != ok2 || v1 != v2 {
				m1.debug()
				t.Fatalf("key:%v got:(%v,%v) expect:(%v,%v)", rk, v1, ok1, v2, ok2)
			}
		}
		if m1.Len() != len(m2) {
			t.Fatal(m1.Len(), len(m2))
		}
	}
}

func TestCorrectness(t *testing.T) {
	m := New64(0)
	if m.Len() != 0 {
		t.Fatal()
	}

	m.Store(1, 111)
	assertHasValue(1, 111, t, m)
	if m.Len() != 1 {
		t.Fatal()
	}

	m.Store(1, 222)
	assertHasValue(1, 222, t, m)
	if m.Len() != 1 {
		t.Fatal()
	}

	var (
		rangecount     int
		rangelastkey   uint64
		rangelastvalue uint64
	)
	m.Range(func(key, value uint64) bool {
		rangecount++
		rangelastkey = key
		rangelastvalue = value
		return true
	})
	if rangecount != 1 || rangelastkey != 1 || rangelastvalue != 222 {
		t.Fatal("invalid")
	}

	m.Delete(1)
	if _, ok := m.Load(1); ok {
		t.Fatal("got deleted key")
	}
	if m.Len() != 0 {
		t.Fatal()
	}

	const mingrowsize = 1000
	for i := uint64(0); i < mingrowsize; i++ {
		m.Store(i, i+10)
		assertHasValue(i, i+10, t, m)
	}
	if m.Len() != mingrowsize {
		t.Fatal()
	}
	for i := uint64(0); i < mingrowsize; i++ {
		assertHasValue(i, i+10, t, m)
	}
	if m.Len() != mingrowsize {
		t.Fatal()
	}
	for i := uint64(0); i < mingrowsize; i++ {
		m.Delete(i)
		if m.Len() != int(mingrowsize-i-1) {
			t.Fatal(m.Len(), int(mingrowsize-i-1))
		}
	}
}

func TestStoreWithoutGrow(t *testing.T) {
	for i := 0; i < 300; i++ {
		m := New64(i)
		initB := m.bucketmask
		for j := 0; j < i; j++ {
			m.Store(uint64(i), uint64(i))
		}
		if initB != m.bucketmask {
			t.Fatal()
		}
	}
}

func assertHasValue(k, v uint64, t *testing.T, m *Uint64Map) {
	got, ok := m.Load(k)
	if !ok || got != v {
		t.Errorf("key %v expected value %v, but got (%v,%v)", k, v, got, ok)
	}
}

func TestSameKeyValue(t *testing.T) {
	const count = 100000
	m1 := New64(0)
	m2 := make(map[uint64]uint64)
	for i := 0; i < count; i++ {
		choose := rand.Intn(100)
		rk, rv := uint64(rand.Int63n(count)), uint64(rand.Int63())
		if choose >= 50 {
			m1.Store(rk, rv)
			m2[rk] = rv
		} else if choose <= 10 {
			m1.Delete(rk)
			delete(m2, rk)
		} else if choose == 27 && rand.Intn(100) == 0 {
			// var (
			// 	tmp1 = make(map[uint64]uint64)
			// 	tmp2 = make(map[uint64]uint64)
			// )
			// m1.Range(func(key, value uint64) bool {
			// 	tmp1[key] = value
			// 	return true
			// })
			// for key, value := range m2 {
			// 	tmp2[key] = value
			// }
			// var x, y int
			// for range tmp1 {
			// 	x++
			// }
			// for range tmp2 {
			// 	y++
			// }
			// if x != y {
			// 	t.Fatal(x, y)
			// }
			// for k2, v2 := range tmp1 {
			// 	v1, ok := tmp2[k2]
			// 	if !ok || v1 != v2 {
			// 		t.Fatal("invalid key value", k2, v2, v1, ok)
			// 	}
			// }
			// for k2, v2 := range tmp2 {
			// 	v1, ok := tmp1[k2]
			// 	if !ok || v1 != v2 {
			// 		t.Fatal("invalid key value", k2, v2, v1, ok)
			// 	}
			// }
		} else {
			// Load
			v1, ok1 := m1.Load(rk)
			v2, ok2 := m2[rk]
			if ok1 != ok2 || v1 != v2 {
				t.Fatalf("key:%v got:(%v,%v) expect:(%v,%v)", rk, v1, ok1, v2, ok2)
			}
		}
		if m1.Len() != len(m2) {
			t.Fatal(m1.Len(), len(m2))
		}
	}
}

func TestNeedBucket(t *testing.T) {
	itemNeedBucket := func(length int) int {
		m := New64(length)
		return int(m.bucketmask) + 1
	}

	for i := 0; i <= maxItemInBucket; i++ {
		if itemNeedBucket(i) != 1 {
			t.Fatal(i, itemNeedBucket(i))
		}
	}

	if itemNeedBucket(maxItemInBucket+1) != 2 {
		t.Fatal()
	}

	if itemNeedBucket(4*maxItemInBucket-1) != 4 {
		t.Fatal()
	}

	if itemNeedBucket(4*maxItemInBucket) != 4 {
		t.Fatal()
	}

	if itemNeedBucket(4*maxItemInBucket+1) != 8 {
		t.Fatal()
	}
}

func TestSameSizeGrow(t *testing.T) {
	type kv struct {
		k uint64
		v uint64
	}
	mcap := (bucketCnt - 1) * 16
	m := New64(0)
	for i := 0; i < 1000; i++ {
		var addk []kv
		for i := m.Len(); i < mcap; i++ {
			r := rand.Int()
			k := uint64(r)
			v := uint64(i)
			m.Store(k, v)
			addk = append(addk, kv{k, v})
		}
		for _, elem := range addk {
			k, v := elem.k, elem.v
			got, ok := m.Load(k)
			if !ok || got != v {
				t.Fatalf("got wrong value (%v,%v), expected (%v,true)", got, ok, v)
			}
		}
		for _, elem := range addk {
			m.Delete(elem.k)
			if m.Len() == 16 {
				break
			}
		}
	}
}

func TestRange(t *testing.T) {
	m := New64(0)
	for i := uint64(0); i < maxItemInBucket; i++ {
		m.Store(i, i)
	}

	var times int
	m.Range(func(_, _ uint64) bool {
		times++
		return false
	})
	if times != 1 {
		t.Fatal()
	}

	times = 0
	m.Range(func(_, _ uint64) bool {
		times++
		return true
	})
	if times != maxItemInBucket {
		t.Fatal()
	}
}

func TestFullBucket(t *testing.T) {
	const bucketNum = 4
	m := New64(maxItemInBucket * bucketNum)

	var toBucket0 []uint64
	for i := uint64(0); len(toBucket0) < bucketNum*bucketCnt; i++ {
		if hashUint64(i)&(bucketNum-1) == 0 {
			toBucket0 = append(toBucket0, i)
		}
	}

	for _, v := range toBucket0 {
		m.Store(v, v)
	}

	var rangevals []uint64
	m.Range(func(key, value uint64) bool {
		if key != value {
			t.Fatal()
		}
		rangevals = append(rangevals, key)
		return true
	})

	if len(rangevals) != len(toBucket0) {
		t.Fatalf("invalid: %d != %d", len(rangevals), len(toBucket0))
	}

	sort.Sort(uint64Slice(rangevals))
	sort.Sort(uint64Slice(toBucket0))

	for i := 0; i < len(rangevals); i++ {
		if rangevals[i] != toBucket0[i] {
			t.Fatalf("invalid: %d != %d", rangevals[i], toBucket0[i])
		}
		if v, ok := m.Load(rangevals[i]); !ok || v != rangevals[i] {
			t.Fatalf("invalid: %t is not true || %d != %d", ok, rangevals[i], toBucket0[i])
		}
	}
}

type uint64Slice []uint64

func (x uint64Slice) Len() int           { return len(x) }
func (x uint64Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x uint64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
