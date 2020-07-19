package main

import (
	"testing"
)

func TestGapBufferInsert(t *testing.T) {

	type cmd struct {
		at int
		s  string
	}

	tests := []struct {
		seq    []cmd
		expect string
	}{
		{
			[]cmd{cmd{0, "foo"}, cmd{0, "bar"}},
			"barfoo",
		},
		{
			[]cmd{cmd{0, "Hello   "}, cmd{8, " World!"}},
			"Hello    World!",
		},
		{
			[]cmd{cmd{0, "foo"}, cmd{1, "bar"}},
			"fbaroo",
		},
		{
			[]cmd{cmd{0, "Hello"}, cmd{2, " World!"}},
			"He World!llo",
		},
	}

	gb := NewGapBuffer()

	for i, test := range tests {
		gb.Clear()
		for _, c := range test.seq {
			gb.Insert(c.at, []byte(c.s))
		}
		s := gb.ToString()
		if s != test.expect {
			t.Fatalf("tests[%d] failed. expected='%s', got='%s'",
				i, test.expect, s)
		}
	}
}

func TestGapBufferClear(t *testing.T) {
	gb := NewGapBuffer()
	gb.Insert(0, []byte("Hello World!"))
	gb.Clear()

	if gb.ToString() != "" {
		t.Fatalf("test failed.\n")
	}
}

func TestGapBufferDelAt(t *testing.T) {

	tests := []struct {
		s      string
		seq    []int
		expect string
	}{
		{
			"abcdefgh",
			[]int{3},
			"abcefgh",
		},
		{
			"abcdefgh",
			[]int{0},
			"bcdefgh",
		},
		{
			"abcdefgh",
			[]int{7},
			"abcdefg",
		},
		{
			"abcdefgh",
			[]int{0, 0, 0},
			"defgh",
		},
		{
			"abcdefgh",
			[]int{3, 3},
			"abcfgh",
		},
		{
			"abcdefgh",
			[]int{1, 6},
			"acdefg",
		},
		{
			"abcdefgh",
			[]int{6, 2},
			"abdefh",
		},
	}

	for i, test := range tests {
		gb := NewGapBuffer()
		gb.Insert(0, []byte(test.s))
		for _, at := range test.seq {
			gb.DelAt(at)
		}

		s := gb.ToString()

		if s != test.expect {
			t.Fatalf("tests[%d] failed. expected='%s', got='%s'",
				i, test.expect, s)
		}
	}
}

func TestGapBufferDel(t *testing.T) {

	type pair struct {
		begin int
		l     int
	}

	tests := []struct {
		s      string
		pairs  []pair
		expect string
	}{
		{
			"abcdefgh",
			[]pair{pair{0, 3}},
			"defgh",
		},
		{
			"abcdefgh",
			[]pair{pair{2, 3}},
			"abfgh",
		},
		{
			"abcdefgh",
			[]pair{pair{0, 8}},
			"",
		},
		{
			"abcdefgh",
			[]pair{pair{3, 5}},
			"abc",
		},
	}

	for i, test := range tests {
		gb := NewGapBuffer()
		gb.Insert(0, []byte(test.s))

		for _, p := range test.pairs {
			gb.Del(p.begin, p.l)
		}

		s := gb.ToString()
		if s != test.expect {
			t.Fatalf("tests[%d] failed. expected='%s', got='%s'",
				i, test.expect, s)
		}
	}

}

func BenchmarkGapBufferInsert(b *testing.B) {
	gb := NewGapBuffer()

	for i := 0; i < 10000; i++ {
		gb.Insert(0, []byte("test"))
	}

}
