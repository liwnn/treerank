# TreeRank
This Go package provides an implementation of ranking-list based on red-black tree. It is in-memory-only and does not store any data. 

## Usage
All you have to do is to implement a comparison function Less(Item) bool for your Item which will be store in the tree, here are some examples.
``` go
package main

import (
	"fmt"

	"github.com/liwnn/treerank"
)

type User struct {
	Name  string
	Score int
}

func (u User) Key() string {
	return u.Name
}

func (u User) Less(than treerank.Item) bool {
	if u.Score == than.(User).Score {
		return u.Name < than.(User).Name
	}
	return u.Score < than.(User).Score
}

func main() {
	tr := treerank.New()

	// Add
	tr.Add("Hurst", User{Name: "Hurst", Score: 88})
	tr.Add("Peek", User{Name: "Peek", Score: 100})
	tr.Add("Beaty", User{Name: "Beaty", Score: 66})

	// Rank
	rank := tr.Rank("Hurst", true)
	fmt.Printf("Hurst's rank is %v\n", rank) // expected 2

	// Range
	tr.Range(0, 3, true, func(key string, v treerank.Item, rank int) bool {
		fmt.Printf("%v's rank is %v\n", v.(User).Name, rank)
		return true
	})

	// Remove
	tr.Remove("Peek")

	// Rank
	rank = tr.Rank("Hurst", true)
	fmt.Printf("Hurst's rank is %v\n", rank) // expected 1
}
```
Output:
```
Hurst's rank is 2
Peek's rank is 1
Hurst's rank is 2
Beaty's rank is 3
Hurst's rank is 1
```