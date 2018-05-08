package util

type Int  int

func (i Int) Less(j Object) bool {
	return i < j.(Int)
}
