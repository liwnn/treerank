package treerank

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func init() {
	seed := time.Now().Unix()
	rand.Seed(seed)
}

// perm returns a random permutation of n Int items in the range [0, n).
func perm(n int) (out []Int) {
	out = make([]Int, 0, n)
	for _, v := range rand.Perm(n) {
		out = append(out, Int(v))
	}
	return
}

// rang returns an ordered list of Int items in the range [0, n).
func rang(n int) (out []Item) {
	for i := 0; i < n; i++ {
		out = append(out, Int(i))
	}
	return
}

func TestRBtreeRank(t *testing.T) {
	const treeSize = 10000
	tr := New()
	for i := 0; i < 10; i++ {
		for _, v := range perm(treeSize) {
			tr.Add(v)
		}
		for _, v := range perm(treeSize) {
			if r := tr.Rank(v.Key(), false); r != int(v)+1 {
				t.Error("rank failed")
			}
			if r := tr.Rank(v.Key(), true); r != int(treeSize-v) {
				t.Error("rank failed")
			}
		}

		if r := tr.Range(0, 1, false); !reflect.DeepEqual(r, rang(2)) {
			t.Error("range error")
		}

		if r := tr.Range(0, 1, true); r[0] != Int(treeSize-1) || r[1] != Int(treeSize-2) {
			t.Error("range error")
		}

		for i := 0; i < treeSize/2; i++ {
			tr.Delete(Int(i).Key())
		}
		for i := treeSize + 1; i < treeSize; i++ {
			if r := tr.Rank(Int(i).Key(), false); r != i-treeSize/2 {
				t.Error("rank failed")
			}
		}
	}
}

const benchmarkTreeSize = 10000

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	insertP := perm(benchmarkTreeSize)
	b.StartTimer()
	i := 0
	for i < b.N {
		tr := New()
		for _, item := range insertP {
			tr.Add(item)
			i++
			if i >= b.N {
				return
			}
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	b.StopTimer()
	insertP := perm(benchmarkTreeSize)
	searchP := perm(benchmarkTreeSize)
	b.StartTimer()
	i := 0
	for i < b.N {
		b.StopTimer()
		tr := New()
		for _, v := range insertP {
			tr.Add(v)
		}
		b.StartTimer()
		for _, item := range searchP {
			tr.Search(item)
			i++
			if i >= b.N {
				return
			}
		}
	}
}

func BenchmarkDeleteInsert(b *testing.B) {
	b.StopTimer()
	insertP := perm(benchmarkTreeSize)
	tr := New()
	for _, item := range insertP {
		tr.Add(item)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tr.Delete(insertP[i%benchmarkTreeSize].Key())
		tr.Add(insertP[i%benchmarkTreeSize])
	}
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	insertP := perm(benchmarkTreeSize)
	removeP := perm(benchmarkTreeSize)
	b.StartTimer()
	i := 0
	for i < b.N {
		b.StopTimer()
		tr := New()
		for _, v := range insertP {
			tr.Add(v)
		}
		b.StartTimer()
		for _, item := range removeP {
			tr.Delete(item.Key())
			i++
			if i >= b.N {
				return
			}
		}
		if tr.Len() > 0 {
			panic(tr.Len())
		}
	}
}
