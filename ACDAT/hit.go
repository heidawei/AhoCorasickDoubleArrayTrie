package ACDAT

import "fmt"

type IHit func(begin, end int, v interface{})

type Hit struct {
	begin    int
	end      int
	value    interface{}
}

func NewHit(begin, end int, v interface{}) *Hit {
	return &Hit{begin:begin, end: end, value: v}
}

func (h *Hit)String() string {
	return fmt.Sprintf("[%d:%d]=%v", h.begin, h.end, h.value)
}
