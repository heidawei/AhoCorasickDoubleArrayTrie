package ac

import "fmt"

type Hit struct {
	begin    int
	end      int
	value    *Word
}

func NewHit(begin, end int, w *Word) *Hit {
	return &Hit{begin:begin, end: end, value: w}
}

func (h *Hit)String() string {
	return fmt.Sprintf("[%d:%d]=%s", h.begin, h.end, h.value)
}
