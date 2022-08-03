package treerank

import (
	"fmt"
	"strconv"
)

const (
	DefaultFreeListSize = 32
)

type Item interface {
	Key() string
	Less(than Item) bool
}

type color int8

// enum
const (
	RED   color = 0
	BLACK color = 1
)

type node struct {
	color color
	item  Item
	left  *node
	right *node
	p     *node
	count int
}

type FreeList struct {
	freelist []*node
}

func NewFreeList(size int) *FreeList {
	return &FreeList{freelist: make([]*node, 0, size)}
}

func (f *FreeList) newNode() (n *node) {
	index := len(f.freelist) - 1
	if index < 0 {
		return new(node)
	}
	n = f.freelist[index]
	f.freelist[index] = nil
	f.freelist = f.freelist[:index]
	return
}

func (f *FreeList) freeNode(n *node) (out bool) {
	n.item = nil
	n.left = nil
	n.right = nil
	n.p = nil
	if len(f.freelist) < cap(f.freelist) {
		f.freelist = append(f.freelist, n)
		out = true
	}
	return
}

type RBTreeRank struct {
	root     *node
	nil      *node
	freelist *FreeList
	length   int
	dict     map[string]*node
}

func New() *RBTreeRank {
	t := &RBTreeRank{
		nil: &node{
			color: BLACK,
		},
		freelist: NewFreeList(DefaultFreeListSize),
		dict:     make(map[string]*node),
	}
	t.root = t.nil
	return t
}

/*
   x               y
  / \             / \
 a   y    ->     x	 c
    / \         / \
   b   c       a   b
*/
func (t *RBTreeRank) leftRotate(x *node) {
	y := x.right
	x.count = x.left.count + y.left.count + 1
	y.count = x.count + y.right.count + 1

	// y的左节点改成x的右节点
	x.right = y.left
	if y.left != t.nil {
		y.left.p = x
	}

	// x 改成y的左节点
	y.left = x
	if x.p == t.nil {
		t.root = y
	} else if x.p.left == x {
		x.p.left = y
	} else {
		x.p.right = y
	}
	y.p = x.p
	x.p = y
}

/*
    y  	       x
   / \        / \
  x	  c  ->  a   y
 / \            / \
a   b          b   c
*/
func (t *RBTreeRank) rightRotate(y *node) {
	x := y.left

	y.count = x.right.count + y.right.count + 1
	x.count = x.left.count + y.count + 1

	y.left = x.right
	if x.right != t.nil {
		x.right.p = y
	}

	x.right = y
	if y.p == t.nil {
		t.root = x
	} else if y.p.left == y {
		y.p.left = x
	} else {
		y.p.right = x
	}
	x.p = y.p
	y.p = x
}

func (t *RBTreeRank) Add(item Item) {
	if item == nil {
		panic("nil item is not allowed in RBTree")
	}

	key := item.Key()
	n, ok := t.dict[key]
	if ok {
		t.delete(n)
	}

	t.insert(item)
}

func (t *RBTreeRank) insert(item Item) *node {
	y := t.nil
	insertLeft := true
	for x := t.root; x != t.nil; {
		y = x
		if x.item.Less(item) {
			x = x.right
			insertLeft = false
		} else {
			x = x.left
			insertLeft = true
		}
	}

	z := t.freelist.newNode()
	z.item = item
	z.p = y
	if y == t.nil {
		t.root = z
	} else if insertLeft {
		y.left = z
	} else {
		y.right = z
	}
	z.left = t.nil
	z.right = t.nil
	z.color = RED
	z.count = 1
	t.length++
	t.dict[z.item.Key()] = z

	for p := z.p; p != t.nil; p = p.p {
		p.count++
	}
	t.insertFixup(z)
	return z
}

func (t *RBTreeRank) insertFixup(z *node) {
	for z.p.color == RED {
		if z.p == z.p.p.left { // z的父节点是左节点
			y := z.p.p.right
			if y.color == RED { // case 1(a): z的叔节点是红
				z.p.color = BLACK
				y.color = BLACK
				z.p.p.color = RED
				z = z.p.p
			} else {
				if z == z.p.right { // case 2: z叔节点是黑色且z是是右孩子
					z = z.p
					t.leftRotate(z)
				}
				// case 3: z叔节点是黑色且z是左孩子
				z.p.color = BLACK
				z.p.p.color = RED
				t.rightRotate(z.p.p)
			}
		} else if z.p == z.p.p.right { // z的父节点是右节点
			y := z.p.p.left
			if y.color == RED { // case 1(b): z叔节点是红
				z.p.color = BLACK
				y.color = BLACK
				z.p.p.color = RED
				z = z.p.p
			} else {
				if z == z.p.left {
					z = z.p
					t.rightRotate(z)
				}
				z.p.color = BLACK
				z.p.p.color = RED
				t.leftRotate(z.p.p)
			}
		}
	}
	t.root.color = BLACK
}

// v替换u
func (t *RBTreeRank) transplant(u *node, v *node) {
	if u.p == t.nil {
		t.root = v
	} else if u.p.left == u {
		u.p.left = v
	} else {
		u.p.right = v
	}
	v.p = u.p
}

func (t *RBTreeRank) Search(item Item) Item {
	if node, ok := t.dict[item.Key()]; ok {
		return node.item
	}
	return nil
}

func (t *RBTreeRank) search(x *node, item Item) *node {
	for x != t.nil {
		if item.Less(x.item) {
			x = x.left
		} else if x.item.Less(item) {
			x = x.right
		} else {
			break
		}
	}
	return x
}

func (t *RBTreeRank) Delete(key string) (removeItem Item) {
	n := t.dict[key]
	if n == nil {
		return nil
	}
	removeItem = n.item
	t.delete(n)
	return
}

func (t *RBTreeRank) delete(z *node) {
	var y = z
	yOriginalColor := y.color
	var x *node
	if z.left == t.nil {
		x = z.right
		t.transplant(z, z.right)
		for p := z.p; p != t.nil; p = p.p {
			p.count--
		}
	} else if z.right == t.nil {
		x = z.left
		t.transplant(z, z.left)
		for p := z.left.p; p != t.nil; p = p.p {
			p.count--
		}
	} else {
		y = t.minimum(z.right)
		yOriginalColor = y.color
		x = y.right
		for p := y.p; p != t.nil; p = p.p {
			p.count--
		}
		if y.p == z {
			x.p = y // t.nil
		} else {
			t.transplant(y, y.right)
			y.right = z.right
			y.right.p = y
		}
		t.transplant(z, y)
		y.left = z.left
		y.left.p = y
		y.color = z.color
		y.count = z.count
	}
	t.length--
	delete(t.dict, z.item.Key())
	t.freelist.freeNode(z)

	if yOriginalColor == BLACK {
		t.deleteFixup(x)
	}
}

func (t *RBTreeRank) minimum(x *node) *node {
	for x.left != t.nil {
		x = x.left
	}
	return x
}

func (t *RBTreeRank) maximum(x *node) *node {
	for x.right != t.nil {
		x = x.right
	}
	return x
}

func (t *RBTreeRank) deleteFixup(x *node) {
	for x != t.root && x.color == BLACK {
		if x == x.p.left {
			w := x.p.right
			if w.color == RED { // case 1: x的兄弟节点w是红色
				w.color = BLACK
				x.p.color = RED
				t.leftRotate(x.p)
				w = x.p.right
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				// case 2: x的兄弟节点w是黑色的, 而且w的两个孩子都是黑色
				w.color = RED
				x = x.p
			} else {
				if w.right.color == BLACK {
					// case 3: x的兄弟节点w是黑色的, w的左孩子是红色, w的右孩子是黑色
					w.left.color = BLACK
					w.color = RED
					t.rightRotate(w)
					w = x.p.right
				}
				// case 4: x的兄弟节点w是黑色的, w的左孩子黑色, w的右孩子是红色
				w.color = x.p.color
				x.p.color = BLACK
				w.right.color = BLACK
				t.leftRotate(x.p)
				x = t.root
			}
		} else {
			w := x.p.left
			if w.color == RED {
				w.color = BLACK
				x.p.color = RED
				t.rightRotate(x.p)
				w = x.p.left
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.p
			} else {
				if w.left.color == BLACK {
					w.right.color = BLACK
					w.color = RED
					t.leftRotate(w)
					w = x.p.left
				}
				w.color = x.p.color
				x.p.color = BLACK
				w.left.color = BLACK
				t.rightRotate(x.p)
				x = t.root
			}
		}
	}
	x.color = BLACK
}

func (t *RBTreeRank) successor(x *node) *node {
	if x.right != t.nil {
		return t.minimum(x.right)
	}
	y := x.p
	for y != t.nil && x == y.right {
		x = y
		y = y.p
	}
	return y
}

func (t *RBTreeRank) predecessor(x *node) *node {
	if x.left != t.nil {
		return t.maximum(x.left)
	}
	y := x.p
	for y != t.nil && x == y.left {
		x = y
		y = y.p
	}
	return y
}

func (t *RBTreeRank) Len() int {
	return t.length
}

func (t *RBTreeRank) NewAscendIterator() *Iterator {
	return &Iterator{t: t, x: t.minimum(t.root)}
}

// Rank return 1-based rank or 0 if not exist
func (t *RBTreeRank) Rank(key string, reverse bool) (count int) {
	n := t.dict[key]
	if n == nil {
		return -1
	}
	x := t.root
	for x != t.nil {
		if x == n {
			if reverse {
				return t.length - (count + x.left.count)
			}
			return count + x.left.count + 1
		}
		if x.item.Less(n.item) {
			count += x.count - x.right.count
			x = x.right
		} else {
			x = x.left
		}
	}
	return 0
}

func (t *RBTreeRank) Range(start, end int, reverse bool) []Item {
	llen := t.length
	if start < 0 {
		start = llen + start
	}
	if end < 0 {
		end = llen + end
	}
	if start < 0 {
		start = 0
	}

	if start > end || start >= llen {
		return nil
	}

	if end >= llen {
		end = llen - 1
	}

	var count = end - start + 1
	var items = make([]Item, 0, count)

	if reverse {
		start = t.length - 1 - start
		end = t.length - 1 - end
	}

	x := t.root
	var startNode *node
	for x != t.nil {
		if x.left.count < start {
			start = start - x.left.count - 1
			x = x.right
		} else if x.left.count == start {
			startNode = x
			break
		} else {
			x = x.left
		}
	}
	it := &Iterator{t: t, x: startNode}
	if reverse {
		for ; it.Valid(); it.Prev() {
			items = append(items, it.Value())
			if len(items) >= count {
				break
			}
		}
	} else {
		for ; it.Valid(); it.Next() {
			items = append(items, it.Value())
			if len(items) >= count {
				break
			}
		}
	}
	return items
}

type Int int

func (a Int) Less(b Item) bool {
	return a < b.(Int)
}
func (a Int) Key() string {
	return strconv.Itoa(int(a))
}

// PrintTree 打印树
func PrintTree(t *RBTreeRank) {
	const (
		nilStr = "nil"
		indent = 2
	)
	levelNode := make(map[int][]*node)
	levelNode[0] = []*node{t.root}
	for level := 0; ; level++ {
		var nodes = levelNode[level]
		var next []*node
		for _, n := range nodes {
			if n != nil {
				next = append(next, n.left, n.right)
			} else {
				next = append(next, nil, nil)
			}
		}
		var exit = true
		for _, v := range next {
			if v != nil {
				exit = false
				break
			}
		}
		if exit {
			break
		}
		levelNode[level+1] = next
	}
	depth := len(levelNode)
	for j := 0; j < depth; j++ {
		w := indent << (depth - 1 - j)
		if j > 0 {
			for i := 0; i < 1<<(j-1); i++ {
				if levelNode[j][i*2] == nil {
					fmt.Printf("%*c", w*4, ' ')
				} else {
					fmt.Printf("%*c", w, ' ') // w
					if w < 3 {
						leftW := 3
						if w == 1 {
							fmt.Printf("| ")
						} else {
							fmt.Printf("/ \\")
						}
						leftW -= 3 / w
						n := w - 3%w + leftW
						fmt.Printf("%*c", n, ' ')
					} else {
						fmt.Printf("%c", ' ')
						for k := 0; k < w-3; k++ {
							fmt.Printf("_")
						}
						fmt.Printf("/ \\")
						for k := 0; k < w-3; k++ {
							fmt.Printf("_")
						}
						fmt.Printf("%*c", w+2, ' ')
					}
				}
			}

			fmt.Printf("\n")
			for i := 0; i < 1<<(j-1); i++ {
				if levelNode[j][i*2] == nil {
					fmt.Printf("%*c", w*4, ' ')
				} else {
					if w < 3 {
						fmt.Printf("%*c%*c%*c", w, '/', w*2, '\\', w, ' ')
					} else {
						fmt.Printf("%*c%*c%*c", w+1, '/', w*2-2, '\\', w+1, ' ')
					}
				}
			}
			fmt.Printf("\n")
		}
		for i := 0; i < 1<<j; i++ {
			n := levelNode[j][i]
			if n == nil {
				fmt.Printf("%*c", w*2, ' ')
				continue
			}
			key := fmt.Sprintf("%v", n.item)
			if n.item == nil {
				key = nilStr
			}
			shiftLeft := (len(key) + 1) / 2
			if w < 3 {
				if i%2 == 0 || len(key) > 2 {
					shiftLeft = (len(key))/2 + 1
				} else {
					shiftLeft = (len(key) + 1) / 2
				}
			}
			if shiftLeft > w {
				shiftLeft = w
			}
			if w > shiftLeft {
				fmt.Printf("%*c", w-shiftLeft, ' ') // (key)
			}
			if n.color == RED {
				fmt.Printf("%c[1;41;37m%v%c[0m", 0x1B, key, 0x1B)
			} else {
				fmt.Printf("%c[1;40;37m%v%c[0m", 0x1B, key, 0x1B)
			}
			fmt.Printf("%*c", w-(len(key)-shiftLeft), ' ')
		}
		fmt.Printf("\n")
	}
}
