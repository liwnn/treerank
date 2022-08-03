# TreeRank
This Go package provides an implementation of ranking-list based on red-black tree. It is in-memory-only and does not store any data. 

## Usage
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
	tr.Add(User{Name: "Hurst", Score: 88})
	tr.Add(User{Name: "Peek", Score: 100})
	tr.Add(User{Name: "Beaty", Score: 66})

	// Rank
	rank := tr.Rank(User{Name: "Hurst"}, true)
	fmt.Printf("Hurst's rank is %v\n", rank) // expected 2

	// Range
	rg := tr.Range(0, 3, true)
	for i, v := range rg {
		fmt.Printf("%v's rank is %v\n", v.(User).Name, i+1)
	}

	// Delete
	tr.Delete(User{Name: "Peek"})

	// Rank
	rank = tr.Rank(User{Name: "Hurst"}, true)
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