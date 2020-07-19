package main

import (
	"fmt"
)

type gapBuffer struct {
	buf   []byte
	begin int
	end   int
}

func NewGapBuffer() *gapBuffer {
	return &gapBuffer{make([]byte, 8, 8), 0, 8}
}

func (gb *gapBuffer) InsertAt(idx int, c byte) (err error) {
	// there is no gap
	if gb.begin == gb.end {
		c := len(gb.buf)
		gb.buf = append(gb.buf, make([]byte, c)...)
		gb.begin = c
		gb.end = len(gb.buf)
		for i := gb.begin; i < c-1; i++ {
			gb.buf[i+c] = gb.buf[i]
		}
	}

	// there is no gap at the insertion point
	if idx < gb.begin {
		r := gb.end - gb.begin
		for i := idx; i < gb.begin; i++ {
			gb.buf[i+r] = gb.buf[i]
		}
		gb.begin = idx
		gb.end = idx + r
	} else if gb.end < idx {
		r := gb.end - gb.begin
		for i := gb.end; i < idx; i++ {
			gb.buf[i-r] = gb.buf[i]
		}
		gb.begin = idx - r
		gb.end = idx
	}

	// there is gap at the insertion point
	if idx == gb.begin {
		gb.buf[idx] = c
		gb.begin++
		return nil
	} else if idx == gb.end-1 {
		gb.buf[idx] = c
		gb.end--
		return nil
	}

	return fmt.Errorf("out of range [%d] with length %d", idx, gb.Length())
}

func (gb *gapBuffer) Insert(idx int, p []byte) {
	for i, c := range p {
		gb.InsertAt(idx+i, c)
	}
}

func (gb *gapBuffer) DelAt(idx int) (err error) {
	if idx > gb.Length()-1 {
		return fmt.Errorf("out of range [%d] with length %d", idx, gb.Length())
	}

	at := idx
	if idx >= gb.begin {
		at = idx + (gb.end - gb.begin)
	}

	// there is no gap
	if gb.begin == gb.end {
		gb.begin = at
		gb.end = at + 1
		return nil
	}

	// there is a no gap where to delete
	if at < gb.begin {
		for i := at; i < gb.begin; i++ {
			gb.buf[i] = gb.buf[i+1]
		}
		gb.begin--
		return nil
	} else if gb.end < at {
		for i := at; i > gb.end; i-- {
			gb.buf[i] = gb.buf[i-1]
		}
		gb.end++
		return nil
	}

	// there is a gap next to where to delete
	if gb.begin == at {
		gb.begin--
		return nil
	} else if gb.end == at {
		gb.end++
		return nil
	}

	return fmt.Errorf("")
}

func (gb *gapBuffer) Del(begin int, l int) {
	for i := 0; i < l; i++ {
		gb.DelAt(begin)
	}
}

func (gb *gapBuffer) At(idx int) byte {
	return gb.ToString()[idx]
}

func (gb *gapBuffer) Clear() {
	gb.buf = make([]byte, 8, 8)
	gb.begin = 0
	gb.end = 8
}

func (gb *gapBuffer) Length() int {
	return cap(gb.buf) - (gb.end - gb.begin)
}

func (gb *gapBuffer) ToString() string {
	if gb.begin == 0 && gb.end == cap(gb.buf) {
		return ""
	}

	head := []byte{}
	if gb.begin > 0 {
		head = make([]byte, gb.begin)
		copy(head, gb.buf[0:gb.begin])
	}
	tail := make([]byte, cap(gb.buf)-gb.end)
	copy(tail, gb.buf[gb.end:cap(gb.buf)])

	return string(append(head, tail...))
}

func (gb *gapBuffer) Repr() string {
	s := ""

	for i, c := range gb.buf {
		if gb.begin <= i && i < gb.end {
			s += "_"
		} else {
			s += string(c)
		}
	}

	return s
}

func main() {
	gb := NewGapBuffer()
	gb.Insert(0, []byte("Hello "))
	gb.Insert(6, []byte("World!"))
	gb.Del(4, 1)
	fmt.Printf("gb: %s, begin=%d, end=%d\n", gb.Repr(), gb.begin, gb.end)
	fmt.Printf("gb: %s\n", gb.ToString())
}
