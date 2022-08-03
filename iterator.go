package treerank

type Iterator struct {
	t *RBTreeRank
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
