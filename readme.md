# umap

umap is a fast, reliable, simple hashmap that only supports the uint64 key/value pair. It is faster than the runtime hashmap in most cases.

The umap is mainly based on the SwissTable(https://github.com/abseil/abseil-cpp/blob/master/absl/container/internal/raw_hash_set.h).



### Goals

- Fast, reliable, and simple.



### Non-goals

- Support other key/value types. Other types require more complex logic, if you need a fast map that supports any type, see https://github.com/golang/go/issues/54766
- Panic in some illegal cases. umap disables all these features, such as panic when multiple goroutines insert items. It's the user's responsibility to check that.
- Concurrent safe. A mutex is required if you need this feature.



### Compared to map[uint64]uint64

The umap is compatible with map[uint64]uint64 for most cases, except:

- DO NOT support inserting items during iterating the map.

> Inserting items during iterating the runtime map is undefined behavior for most cases, compatible with this feature requires lots of work, so we ban this feature in the umap. **But deleting items during iterating the map is fine.**
>
> i.e. This code prints different numbers each time.
>
> ```golang
> func main() {
> 	m := make(map[int]int, 0)
> 	for i := 0; i < 100; i++ {
> 		m[i] = i
> 	}
> 	j := 1000
> 	for k, _ := range m {
> 		m[j+k] = j
> 		j++
> 	}
> 	println(len(m))
> }
> ```
>
> It means that the code below is INVALID! Although it will not panic.
>
> ```golang
> 	m.Range(func(key, value uint64) bool {
> 		m.Store(key+100, value)
> 		return true
> 	})
> ```
>
> 



## Benchmark

Go version: go version devel go1.20-a813be86df

CPU: Intel 11700k

OS: 22.04.1 LTS (Jammy Jellyfish)

MEMORY: 16G x 2 (3200MHz)

```
name                            old time/op  new time/op   delta
MapAccessHit/Uint64/6-16        3.67ns ± 0%   2.82ns ± 0%   -23.16%  (p=0.000 n=8+10)
MapAccessHit/Uint64/12-16       5.13ns ± 1%   2.82ns ± 0%   -45.00%  (p=0.000 n=8+10)
MapAccessHit/Uint64/18-16       5.39ns ±18%   2.82ns ± 0%   -47.63%  (p=0.000 n=10+10)
MapAccessHit/Uint64/24-16       5.34ns ± 6%   2.82ns ± 0%   -47.26%  (p=0.000 n=9+10)
MapAccessHit/Uint64/30-16       4.63ns ± 5%   2.82ns ± 0%   -39.09%  (p=0.000 n=10+10)
MapAccessHit/Uint64/64-16       4.79ns ± 5%   2.82ns ± 1%   -41.21%  (p=0.000 n=10+10)
MapAccessHit/Uint64/128-16      4.85ns ± 6%   2.82ns ± 0%   -41.76%  (p=0.000 n=10+10)
MapAccessHit/Uint64/256-16      4.75ns ± 2%   2.82ns ± 0%   -40.58%  (p=0.000 n=10+10)
MapAccessHit/Uint64/512-16      4.74ns ± 2%   2.82ns ± 0%   -40.49%  (p=0.000 n=9+10)
MapAccessHit/Uint64/1024-16     4.75ns ± 1%   2.82ns ± 0%   -40.71%  (p=0.000 n=10+10)
MapAccessHit/Uint64/2048-16     5.82ns ± 5%   2.93ns ± 0%   -49.69%  (p=0.000 n=10+10)
MapAccessHit/Uint64/4096-16     11.1ns ± 1%    3.0ns ± 0%   -72.71%  (p=0.000 n=10+10)
MapAccessHit/Uint64/8192-16     13.0ns ± 1%    3.2ns ± 1%   -75.23%  (p=0.000 n=9+10)
MapAccessHit/Uint64/65536-16    16.8ns ± 1%    4.6ns ± 0%   -72.37%  (p=0.000 n=9+8)
MapRange/Uint64/6-16            48.0ns ± 0%   12.1ns ± 1%   -74.85%  (p=0.000 n=9+10)
MapRange/Uint64/12-16           84.7ns ± 5%   23.0ns ± 0%   -72.82%  (p=0.000 n=10+10)
MapRange/Uint64/18-16            136ns ± 3%     40ns ± 1%   -70.75%  (p=0.000 n=10+9)
MapRange/Uint64/24-16            159ns ± 3%     44ns ± 1%   -72.09%  (p=0.000 n=10+10)
MapRange/Uint64/30-16            218ns ± 5%     71ns ± 0%   -67.52%  (p=0.000 n=10+10)
MapRange/Uint64/64-16            418ns ± 1%    148ns ± 0%   -64.51%  (p=0.000 n=9+10)
MapRange/Uint64/128-16           824ns ± 3%    289ns ± 0%   -64.93%  (p=0.000 n=10+9)
MapRange/Uint64/256-16          1.73µs ± 4%   0.57µs ± 0%   -66.95%  (p=0.000 n=10+10)
MapRange/Uint64/512-16          3.75µs ± 3%   1.14µs ± 0%   -69.64%  (p=0.000 n=10+10)
MapRange/Uint64/1024-16         7.88µs ± 1%   2.27µs ± 0%   -71.21%  (p=0.000 n=10+9)
MapRange/Uint64/2048-16         16.4µs ± 3%    4.5µs ± 0%   -72.29%  (p=0.000 n=10+9)
MapRange/Uint64/4096-16         32.9µs ± 1%    9.1µs ± 1%   -72.35%  (p=0.000 n=9+10)
MapRange/Uint64/8192-16         65.9µs ± 0%   19.1µs ± 1%   -71.08%  (p=0.000 n=9+10)
MapRange/Uint64/65536-16         525µs ± 0%    186µs ± 2%   -64.57%  (p=0.000 n=8+9)
MapAssignGrow/Uint64/6-16       53.1ns ± 0%  255.8ns ±24%  +381.74%  (p=0.000 n=10+9)
MapAssignGrow/Uint64/12-16       662ns ±24%    705ns ±13%      ~     (p=0.280 n=10+10)
MapAssignGrow/Uint64/18-16      1.66µs ±27%   1.47µs ±31%      ~     (p=0.113 n=10+9)
MapAssignGrow/Uint64/24-16      1.74µs ±11%   1.63µs ±14%      ~     (p=0.139 n=10+9)
MapAssignGrow/Uint64/30-16      4.72µs ±34%   2.65µs ±32%   -43.77%  (p=0.000 n=10+10)
MapAssignGrow/Uint64/64-16      9.79µs ±32%   5.60µs ±23%   -42.75%  (p=0.000 n=10+10)
MapAssignGrow/Uint64/128-16     21.8µs ±31%   10.7µs ±25%   -50.67%  (p=0.000 n=10+10)
MapAssignGrow/Uint64/256-16     36.4µs ±29%   20.1µs ±37%   -44.72%  (p=0.000 n=10+10)
MapAssignGrow/Uint64/512-16     81.8µs ±21%   34.0µs ±50%   -58.44%  (p=0.000 n=10+10)
MapAssignGrow/Uint64/1024-16     161µs ±30%     46µs ±73%   -71.75%  (p=0.000 n=10+10)
MapAssignGrow/Uint64/2048-16     330µs ±19%    121µs ±53%   -63.32%  (p=0.000 n=10+10)
MapAssignGrow/Uint64/4096-16     596µs ±23%    261µs ±31%   -56.16%  (p=0.000 n=10+10)
MapAssignGrow/Uint64/8192-16    1.39ms ±29%   0.68ms ±31%   -50.99%  (p=0.000 n=10+10)
MapAssignGrow/Uint64/65536-16   8.11ms ±14%   5.73ms ±12%   -29.41%  (p=0.000 n=10+10)
MapAssignReuse/Uint64/6-16       151ns ± 0%     92ns ± 1%   -38.77%  (p=0.000 n=8+10)
MapAssignReuse/Uint64/12-16      398ns ± 0%    159ns ± 1%   -60.20%  (p=0.000 n=9+10)
MapAssignReuse/Uint64/18-16      636ns ± 0%    219ns ± 1%   -65.53%  (p=0.000 n=8+8)
MapAssignReuse/Uint64/24-16      867ns ± 0%    300ns ± 4%   -65.36%  (p=0.000 n=10+9)
MapAssignReuse/Uint64/30-16     1.07µs ± 0%   0.36µs ± 1%   -66.17%  (p=0.000 n=9+10)
MapAssignReuse/Uint64/64-16     2.25µs ± 1%   0.77µs ± 0%   -66.00%  (p=0.000 n=10+9)
MapAssignReuse/Uint64/128-16    4.47µs ± 0%   1.53µs ± 3%   -65.82%  (p=0.000 n=9+8)
MapAssignReuse/Uint64/256-16    8.85µs ± 0%   3.05µs ± 0%   -65.51%  (p=0.000 n=10+8)
MapAssignReuse/Uint64/512-16    17.5µs ± 0%    6.1µs ± 2%   -64.99%  (p=0.000 n=9+10)
MapAssignReuse/Uint64/1024-16   35.1µs ± 0%   12.3µs ± 1%   -65.06%  (p=0.000 n=9+10)
MapAssignReuse/Uint64/2048-16   70.6µs ± 0%   24.6µs ± 1%   -65.09%  (p=0.000 n=9+10)
MapAssignReuse/Uint64/4096-16    143µs ± 1%     50µs ± 1%   -65.38%  (p=0.000 n=10+10)
MapAssignReuse/Uint64/8192-16    291µs ± 0%    100µs ± 1%   -65.60%  (p=0.000 n=10+10)
MapAssignReuse/Uint64/65536-16  2.58ms ± 0%   0.96ms ± 2%   -62.77%  (p=0.000 n=9+10)
```

Tips:

- umap is slower in `MapAssignGrow/Uint64/6`, the reason is that in this case, the runtime hashmap is allocated in the stack instead of the heap. umap is always allcated in the heap.

- You can run this benchmark via

  ```bash
  $ env BENCH_TYPE=runtime go test -bench=. -count=10 -timeout=10h > a.txt
  $ go test -bench=. -count=10 -timeout=10h > b.txt
  $ benchstat a.txt b.txt
  ```

