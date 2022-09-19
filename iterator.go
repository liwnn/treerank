package treerank

type Iterator struct {
	t *RBTree
	x *node
}

func (it Iterator) Valid() bool {
	return it.x != it.t.nil
}

func (it *Iterator) Next() {
	it.x = it.t.successor(it.x)
}

func (it *Iterator) Prev() {
	it.x = it.t.predecessor(it.x)
}

func (it Iterator) Value() Item {
	return it.x.item
}

type RangeIterator struct {
	t               *RBTree
	node            *node
	start, end, cur int
	reverse         bool
}

func (r *RangeIterator) Len() int {
	return r.end - r.start + 1
}

func (r *RangeIterator) Valid() bool {
	return r.cur <= r.end
}

func (r *RangeIterator) Next() {
	if r.reverse {
		r.node = r.t.predecessor(r.node)
	} else {
		r.node = r.t.successor(r.node)
	}
	r.cur++
}

func (r *RangeIterator) Item() Item {
	return r.node.item
}

func (r *RangeIterator) Key() string {
	return r.node.key
}

func (r *RangeIterator) Rank() int {
	return r.cur + 1
}
