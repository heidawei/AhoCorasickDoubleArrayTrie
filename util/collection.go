package util

type Object interface {
	Less(Object) bool
}

type Collection interface {
	Add(o Object) bool
	Remove(o Object) bool
	Size() int
	Contains(o Object) bool
	IsEmpty() bool
	Iterator() Iterator
	ContainAll(c Collection) bool
	AddAll(c Collection) bool
	Clear()
	// 从集合中删除c集合中也有的元素
	RemoveAll(c Collection)
	// 从集合中删除集合c中不包含的元素
	RetainAll(c Collection)
}

type Iterator interface {
	Next() Object
}