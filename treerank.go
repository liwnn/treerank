// Package treerank implements ranking-list based on red-black tree.
package treerank

const (
	DefaultFreeListSize = 32
)

// Item represents a single object in the set.
type Item interface {
	Less(than Item) bool
}

// ItemIterator allows callers of Range* to iterate of the zset.
// When this function returns false, iteration will stop.
type ItemIterator func(key string, i Item, rank int) bool

type color int8

// enum
const (
	RED   color = 0
	BLACK color = 1
)

type node struct {
	color color
	left  *node
	right *node
	p     *node
	count int
	key   string
	item  Item
}

type FreeList struct {
	freelist []*node
}

func NewFreeList(size int) *FreeList {
	return &FreeList{freelist: make([]*node, 0, size)}
}

func (f *FreeList) newNode() (n *node) {
	if len(f.freelist) == 0 {
		return new(node)
	}
	index := len(f.freelist) - 1
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

type RBTree struct {
	root     *node
	nil      *node
	freelist *FreeList
}

//	 x               y
//	/ \             / \
//
// a   y    ->     x   c
//
//	 / \         / \
//	b   c       a   b
func (t *RBTree) leftRotate(x *node) {
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

//	   y          x
//	  / \        / \
//	 x   c  ->  a   y
//	/ \            / \
//
// a   b          b   c
func (t *RBTree) rightRotate(y *node) {
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

func (t *RBTree) insert(key string, item Item) *node {
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
	z.key = key
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

	for p := z.p; p != t.nil; p = p.p {
		p.count++
	}
	t.insertFixup(z)
	return z
}

func (t *RBTree) insertFixup(z *node) {
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
func (t *RBTree) transplant(u *node, v *node) {
	if u.p == t.nil {
		t.root = v
	} else if u.p.left == u {
		u.p.left = v
	} else {
		u.p.right = v
	}
	v.p = u.p
}

func (t *RBTree) search(x *node, item Item) *node {
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

func (t *RBTree) delete(z *node) {
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
			x.p = y // x maybe t.nil, reassign p to y.
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
	t.freelist.freeNode(z)

	if yOriginalColor == BLACK {
		t.deleteFixup(x)
	}
}

func (t *RBTree) updateItem(x *node, item Item) bool {
	successor := t.successor(x)
	if successor == t.nil || !successor.item.Less(item) {
		predecessor := t.predecessor(x)
		if predecessor == t.nil || !item.Less(predecessor.item) {
			x.item = item
			return true
		}
	}
	return false
}

func (t *RBTree) minimum(x *node) *node {
	for x.left != t.nil {
		x = x.left
	}
	return x
}

func (t *RBTree) maximum(x *node) *node {
	for x.right != t.nil {
		x = x.right
	}
	return x
}

func (t *RBTree) length() int {
	return t.root.count
}

func (t *RBTree) deleteFixup(x *node) {
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

func (t *RBTree) successor(x *node) *node {
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

func (t *RBTree) predecessor(x *node) *node {
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

func (t *RBTree) getLessCount(n *node) (count int) {
	x := t.root
	for x != t.nil {
		if x == n {
			return count + x.left.count
		}
		if x.item.Less(n.item) {
			count += x.count - x.right.count
			x = x.right
		} else {
			x = x.left
		}
	}
	return -1
}

func (t *RBTree) getNodeBySortIndex(index int) *node {
	x := t.root
	for x != t.nil {
		if x.left.count < index {
			index = index - x.left.count - 1
			x = x.right
		} else if x.left.count > index {
			x = x.left
		} else {
			return x
		}
	}
	return x
}

type RBTreeRank struct {
	rbTree RBTree
	dict   map[string]*node
}

func New() *RBTreeRank {
	t := &RBTreeRank{
		rbTree: RBTree{
			nil:      &node{color: BLACK},
			freelist: NewFreeList(DefaultFreeListSize),
		},
		dict: make(map[string]*node),
	}
	t.rbTree.root = t.rbTree.nil
	return t
}

func (t *RBTreeRank) Add(key string, item Item) {
	if item == nil {
		panic("nil item is not allowed in RBTree")
	}

	n, ok := t.dict[key]
	if ok {
		if t.rbTree.updateItem(n, item) {
			return
		}
		t.rbTree.delete(n)
	}

	t.dict[key] = t.rbTree.insert(key, item)
}

// Remove the element 'ele' from the rank.
func (t *RBTreeRank) Remove(key string) (removeItem Item) {
	n := t.dict[key]
	if n == nil {
		return nil
	}
	removeItem = n.item
	t.rbTree.delete(n)
	delete(t.dict, key)
	return
}

// Get return Item in dict.
func (t *RBTreeRank) Get(key string) Item {
	if node, ok := t.dict[key]; ok {
		return node.item
	}
	return nil
}

// Rank return 1-based rank or 0 if not exist
func (t *RBTreeRank) Rank(key string, reverse bool) (count int) {
	n := t.dict[key]
	if n == nil {
		return -1
	}

	lessCount := t.rbTree.getLessCount(n)
	if lessCount < 0 {
		return 0
	}
	if reverse {
		return t.rbTree.length() - lessCount
	}
	return lessCount + 1
}

func (t *RBTreeRank) Range(start, end int, reverse bool, iterator ItemIterator) {
	llen := t.rbTree.length()
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
		return
	}

	if end >= llen {
		end = llen - 1
	}

	var count = end - start + 1
	if reverse {
		// end = t.rbTree.length() - 1 - end
		x := t.rbTree.getNodeBySortIndex(t.rbTree.length() - 1 - start)
		for i := 1; i <= count; i++ {
			if iterator(x.key, x.item, start+i) {
				x = t.rbTree.predecessor(x)
			} else {
				break
			}
		}

	} else {
		x := t.rbTree.getNodeBySortIndex(start)
		for i := 1; i <= count; i++ {
			if iterator(x.key, x.item, start+i) {
				x = t.rbTree.successor(x)
			} else {
				break
			}
		}
	}
}

// RangeIterator return iterator for visit elements in [start, end].
// It is slower than Range.
func (t *RBTreeRank) RangeIterator(start, end int, reverse bool) RangeIterator {
	llen := t.rbTree.length()
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
		return RangeIterator{end: -1}
	}

	if end >= llen {
		end = llen - 1
	}

	var n *node
	if reverse {
		n = t.rbTree.getNodeBySortIndex(t.rbTree.length() - 1 - start)
	} else {
		n = t.rbTree.getNodeBySortIndex(start)
	}
	return RangeIterator{
		t:       &t.rbTree,
		start:   start,
		cur:     start,
		end:     end,
		node:    n,
		reverse: reverse,
	}
}

func (t *RBTreeRank) Length() int {
	return t.rbTree.root.count
}

func (t *RBTreeRank) NewAscendIterator() *Iterator {
	return &Iterator{t: &t.rbTree, x: t.rbTree.minimum(t.rbTree.root)}
}

type Int int

func (a Int) Less(b Item) bool {
	return a < b.(Int)
}
