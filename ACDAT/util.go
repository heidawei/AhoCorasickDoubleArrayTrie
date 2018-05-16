package ACDAT

import (
	"github.com/google/btree"
	"github.com/petar/GoLLRB/llrb"
)

type MapItem struct {
	key int
	value interface{}
}

func (m *MapItem)Less(item btree.Item) bool {
	switch i := item.(type) {
	case *MapItem:
		return m.key < i.key
	default:
		return !item.Less(m)
	}
}

type IntTreeMap struct {
	tree *btree.BTree
}

func (t *IntTreeMap)add(key int, value interface{}) {
	t.tree.ReplaceOrInsert(&MapItem{key: key, value: value})
}


type TreeSetItem struct {
	key int
	value interface{}
}

// 反序
func (i *TreeSetItem) Less(item llrb.Item) bool {
	if i == nil {
		return true
	}
	switch it := item.(type) {
	case *TreeSetItem:
		return i.key > it.key
	default:
		return item.Less(i)
	}
	return true
}

type TreeSet struct {
	size int
	tree *llrb.LLRB
	max *TreeSetItem
	min *TreeSetItem
}

func NewTreeSet() *TreeSet {
	return &TreeSet{tree: llrb.New(), }
}

func(s *TreeSet) Add(e int) bool {
	item := &TreeSetItem{key:e}
	it := s.tree.Get(item)
	if it != nil {
		return false
	}
	if s.size == 0 {
		s.max = item
		s.min = item
	} else {
		if s.max.Less(item) {
			s.max = item
		}
		if !s.min.Less(item) {
			s.min = item
		}
	}
	s.tree.ReplaceOrInsert(item)
	s.size++
	return true
}

func(s *TreeSet) Remove(e int) bool {
	item := &TreeSetItem{key:e}
	if it := s.tree.Delete(item); it != nil {
		s.size--
		return true
	}
	return false
}

func(s *TreeSet) Size() int {
	return s.size
}

func(s *TreeSet) Contains(key int) bool {
	return false
}

func(s *TreeSet) IsEmpty() bool {
	if s.size == 0 {
		return true
	}
	return false
}

func(s *TreeSet) Max() int {
	if s.max != nil {
		return s.max.key
	}
	return 0
}

func (s *TreeSet) Min() int {
	if s.min != nil {
		return s.min.key
	}
	return 0
}

func(s *TreeSet) Clear() {
	s.tree = llrb.New()
	s.size = 0
}

func (s *TreeSet) All() []int {
	var sets []int
	s.tree.AscendRange(llrb.Inf(1), llrb.Inf(-1), func (i llrb.Item) bool {
		sets = append(sets, i.(*TreeSetItem).key)
		return true
	})
	return sets
}


type ListHit struct {
	size_   int
	hits  []*Hit
}

func NewListHit() *ListHit {
	return &ListHit{size_: 0}
}

func (l *ListHit) size() int {
	return l.size_
}

// TODO check index > size
func (l *ListHit) get(index int) *Hit {
	if index < 0 {
		return nil
	}
	return l.hits[index]
}

func (l *ListHit) add(hit *Hit) {
	l.hits = append(l.hits, hit)
	l.size_++
}

func (l *ListHit) ListArray() []*Hit {
	return l.hits
}

type ArrayStateQueue struct {
	size int
	index int
	queue  []*State
}

func NewArrayStateQueue() *ArrayStateQueue {
	return &ArrayStateQueue{}
}

func (q *ArrayStateQueue) add(e *State) {
	q.queue = append(q.queue, e)
	q.size++
}

func (q *ArrayStateQueue) isEmpty() bool {
	return q.size == 0 || q.size == q.index + 1
}

func (q *ArrayStateQueue) remove() *State {
	if q.isEmpty() {
		return nil
	}
	e := q.queue[q.index]
	q.index++
	return e
}


type StringMapItem struct {
	key string
	value interface{}
}

func (m *StringMapItem)Less(item btree.Item) bool {
	switch i := item.(type) {
	case *StringMapItem:
		return m.key < i.key
	default:
		return !item.Less(m)
	}
}

type StringTreeMap struct {
	tree *btree.BTree
}

func NewStringTreeMap() *StringTreeMap {
	return &StringTreeMap{tree: btree.New(10)}
}

func (t *StringTreeMap)Add(key string, value interface{}) {
	t.tree.ReplaceOrInsert(&StringMapItem{key: key, value: value})
}

func (t *StringTreeMap) Values() []interface{} {
	var vs []interface{}
	t.tree.Ascend(func(i btree.Item) bool {
		vs = append(vs, i.(*StringMapItem).value)
		return true
	})
	return vs
}

func (t *StringTreeMap) Keys() []string {
	var ks []string
	t.tree.Ascend(func(i btree.Item) bool {
		ks = append(ks, i.(*StringMapItem).key)
		return true
	})
	return ks
}



