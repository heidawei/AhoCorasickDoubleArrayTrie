package util

import (
	"github.com/petar/GoLLRB/llrb"
)

type TreeSetItem struct {
	o Object
}

func (i *TreeSetItem) Less(item llrb.Item) bool {
	if i == nil {
		return true
	}
	switch it := item.(type) {
	case *TreeSetItem:
		return i.o.Less(it.o)
	default:
		return !item.Less(i)
	}
	return true
}

type TreeSet struct {
	size int
	tree *llrb.LLRB
}

func NewTreeSet() *TreeSet {
	return &TreeSet{tree: llrb.New()}
}

func(s *TreeSet) Add(o Object) bool {
	item := &TreeSetItem{o}
	e := s.tree.Get(item)
	if e != nil {
		return false
	}
	s.tree.ReplaceOrInsert(item)
	s.size++
	return true
}

func(s *TreeSet) Remove(o Object) bool {
	item := &TreeSetItem{o}
	if e := s.tree.Delete(item); e != nil {
		return true
	}
	return false
}

func(s *TreeSet) Size() int {
	return s.size
}

func(s *TreeSet) Contains(o Object) bool {
	return false
}
func(s *TreeSet) IsEmpty() bool {
	if s.size == 0 {
		return true
	}
	return false
}

func(s *TreeSet) Iterator() *TreeSetIterator {
	return nil
}

func(s *TreeSet) ContainAll(c Collection) bool {
	return false
}

func(s *TreeSet) AddAll(c Collection) bool {
	return false
}
func(s *TreeSet) Clear() {
	s.tree = llrb.New()
	s.size = 0
}
// 从集合中删除c集合中也有的元素
func(s *TreeSet) RemoveAll(c Collection) {
	return
}
// 从集合中删除集合c中不包含的元素
func(s *TreeSet) RetainAll(c Collection) {
	return
}

type TreeSetIterator struct {

}

func (it *TreeSetIterator) Next() Int {
	return 0
}

