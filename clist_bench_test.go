package clist

import (
	"math"
	"testing"

	"github.com/zhangyunhao116/clist/simplelist"
)

const randN = math.MaxUint32

func BenchmarkInsert(b *testing.B) {
	b.Run("simplelist", func(b *testing.B) {
		l := simplelist.NewInt()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Insert(int(fastrandn(randN)))
			}
		})
	})
	b.Run("concurrentlist", func(b *testing.B) {
		l := NewInt()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Insert(int(fastrandn(randN)))
			}
		})
	})
}

func BenchmarkInsertDupl(b *testing.B) {
	const randN = 500
	b.Run("simplelist", func(b *testing.B) {
		l := simplelist.NewInt()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Insert(int(fastrandn(randN)))
			}
		})
	})
	b.Run("concurrentlist", func(b *testing.B) {
		l := NewInt()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				l.Insert(int(fastrandn(randN)))
			}
		})
	})
}

func Benchmark30Insert70Contains(b *testing.B) {
	b.Run("simplelist", func(b *testing.B) {
		l := simplelist.NewInt()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 3 {
					l.Insert(int(fastrandn(randN)))
				} else {
					l.Contains(int(fastrandn(randN)))
				}
			}
		})
	})
	b.Run("concurrentlist", func(b *testing.B) {
		l := NewInt()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(10)
				if u < 3 {
					l.Insert(int(fastrandn(randN)))
				} else {
					l.Contains(int(fastrandn(randN)))
				}
			}
		})
	})
}

func Benchmark1Delete9Insert90Contains(b *testing.B) {
	const initsize = 1000
	b.Run("simplelist", func(b *testing.B) {
		l := simplelist.NewInt()
		for i := 0; i < initsize; i++ {
			l.Insert(i)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(100)
				if u < 9 {
					l.Insert(int(fastrandn(randN)))
				} else if u == 10 {
					l.Delete(int(fastrandn(randN)))
				} else {
					l.Contains(int(fastrandn(randN)))
				}
			}
		})
	})
	b.Run("concurrentlist", func(b *testing.B) {
		l := NewInt()
		for i := 0; i < initsize; i++ {
			l.Insert(i)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				u := fastrandn(100)
				if u < 9 {
					l.Insert(int(fastrandn(randN)))
				} else if u == 10 {
					l.Delete(int(fastrandn(randN)))
				} else {
					l.Contains(int(fastrandn(randN)))
				}
			}
		})
	})
}
