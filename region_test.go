package gts

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

var regionAccessorTests = []struct {
	in   Region
	len  int
	head int
	tail int
}{
	{Segment{3, 6}, 3, 3, 6},
	{Segment{6, 3}, 3, 6, 3},
	{Regions{}, 0, 0, 0},
	{Regions{Segment{3, 6}, Segment{13, 16}}, 6, 3, 16},
}

func TestRegionAccessor(t *testing.T) {
	for _, tt := range regionAccessorTests {
		if tt.in.Len() != tt.len {
			t.Errorf("%#v.Len() = %d, want %d", tt.in, tt.in.Len(), tt.len)
		}
		if tt.in.Head() != tt.head {
			t.Errorf("%#v.Head() = %d, want %d", tt.in, tt.in.Head(), tt.head)
		}
		if tt.in.Tail() != tt.tail {
			t.Errorf("%#v.Tail() = %d, want %d", tt.in, tt.in.Tail(), tt.len)
		}
	}
}

var regionResizeTests = []struct {
	in       Region
	modifier Modifier
	out      Region
}{
	{Segment{3, 6}, Head(+0), Segment{3, 3}},
	{Segment{3, 6}, Head(+1), Segment{4, 4}},
	{Segment{3, 6}, Head(-1), Segment{2, 2}},

	{Segment{3, 6}, Tail(+0), Segment{6, 6}},
	{Segment{3, 6}, Tail(+1), Segment{7, 7}},
	{Segment{3, 6}, Tail(-1), Segment{5, 5}},

	{Segment{3, 6}, HeadTail{+0, +0}, Segment{3, 6}},
	{Segment{3, 6}, HeadTail{+0, +1}, Segment{3, 7}},
	{Segment{3, 6}, HeadTail{+2, +0}, Segment{5, 6}},
	{Segment{3, 6}, HeadTail{+0, -1}, Segment{3, 5}},
	{Segment{3, 6}, HeadTail{-2, +0}, Segment{1, 6}},
	{Segment{3, 6}, HeadTail{+2, -1}, Segment{5, 5}},
	{Segment{3, 6}, HeadTail{+2, -2}, Segment{5, 5}},
	{Segment{3, 6}, HeadTail{-2, +1}, Segment{1, 7}},
	{Segment{3, 6}, HeadTail{+2, +1}, Segment{5, 7}},
	{Segment{3, 6}, HeadTail{-2, -1}, Segment{1, 5}},

	{Segment{3, 6}, HeadHead{+0, +0}, Segment{3, 3}},
	{Segment{3, 6}, HeadHead{+0, +1}, Segment{3, 4}},
	{Segment{3, 6}, HeadHead{+2, +0}, Segment{5, 5}},
	{Segment{3, 6}, HeadHead{+0, -1}, Segment{3, 3}},
	{Segment{3, 6}, HeadHead{-2, +0}, Segment{1, 3}},
	{Segment{3, 6}, HeadHead{+2, -1}, Segment{5, 5}},
	{Segment{3, 6}, HeadHead{-2, +1}, Segment{1, 4}},
	{Segment{3, 6}, HeadHead{+2, +1}, Segment{5, 5}},
	{Segment{3, 6}, HeadHead{-2, -1}, Segment{1, 2}},

	{Segment{3, 6}, TailTail{+0, +0}, Segment{6, 6}},
	{Segment{3, 6}, TailTail{+0, +1}, Segment{6, 7}},
	{Segment{3, 6}, TailTail{+2, +0}, Segment{8, 8}},
	{Segment{3, 6}, TailTail{+0, -1}, Segment{6, 6}},
	{Segment{3, 6}, TailTail{-2, +0}, Segment{4, 6}},
	{Segment{3, 6}, TailTail{+2, -1}, Segment{8, 8}},
	{Segment{3, 6}, TailTail{-2, +1}, Segment{4, 7}},
	{Segment{3, 6}, TailTail{+2, +1}, Segment{8, 8}},
	{Segment{3, 6}, TailTail{-2, -1}, Segment{4, 5}},

	{Segment{6, 3}, Head(+0), Segment{6, 6}},
	{Segment{6, 3}, Head(+1), Segment{5, 5}},
	{Segment{6, 3}, Head(-1), Segment{7, 7}},

	{Segment{6, 3}, Tail(+0), Segment{3, 3}},
	{Segment{6, 3}, Tail(+1), Segment{2, 2}},
	{Segment{6, 3}, Tail(-1), Segment{4, 4}},

	{Segment{6, 3}, HeadTail{+0, +0}, Segment{6, 3}},
	{Segment{6, 3}, HeadTail{+0, +1}, Segment{6, 2}},
	{Segment{6, 3}, HeadTail{+2, +0}, Segment{4, 3}},
	{Segment{6, 3}, HeadTail{+0, -1}, Segment{6, 4}},
	{Segment{6, 3}, HeadTail{-2, +0}, Segment{8, 3}},
	{Segment{6, 3}, HeadTail{+2, -1}, Segment{4, 4}},
	{Segment{6, 3}, HeadTail{-2, +1}, Segment{8, 2}},
	{Segment{6, 3}, HeadTail{+2, +1}, Segment{4, 2}},
	{Segment{6, 3}, HeadTail{-2, -1}, Segment{8, 4}},

	{Segment{6, 3}, HeadHead{+0, +0}, Segment{6, 6}},
	{Segment{6, 3}, HeadHead{+0, +1}, Segment{6, 5}},
	{Segment{6, 3}, HeadHead{+2, +0}, Segment{4, 4}},
	{Segment{6, 3}, HeadHead{+0, -1}, Segment{6, 6}},
	{Segment{6, 3}, HeadHead{-2, +0}, Segment{8, 6}},
	{Segment{6, 3}, HeadHead{+2, -1}, Segment{4, 4}},
	{Segment{6, 3}, HeadHead{-2, +1}, Segment{8, 5}},
	{Segment{6, 3}, HeadHead{+2, +1}, Segment{4, 4}},
	{Segment{6, 3}, HeadHead{-2, -1}, Segment{8, 7}},

	{Segment{6, 3}, TailTail{+0, +0}, Segment{3, 3}},
	{Segment{6, 3}, TailTail{+0, +1}, Segment{3, 2}},
	{Segment{6, 3}, TailTail{+2, +0}, Segment{1, 1}},
	{Segment{6, 3}, TailTail{+0, -1}, Segment{3, 3}},
	{Segment{6, 3}, TailTail{-2, +0}, Segment{5, 3}},
	{Segment{6, 3}, TailTail{+2, -1}, Segment{1, 1}},
	{Segment{6, 3}, TailTail{-2, +1}, Segment{5, 2}},
	{Segment{6, 3}, TailTail{+2, +1}, Segment{1, 1}},
	{Segment{6, 3}, TailTail{-2, -1}, Segment{5, 4}},

	{Regions{Segment{3, 6}, Segment{13, 16}}, Head(+0), Segment{3, 3}},
	{Regions{Segment{3, 6}, Segment{13, 16}}, Head(+7), Segment{17, 17}},
	{Regions{Segment{3, 6}, Segment{13, 16}}, Tail(+0), Segment{16, 16}},
	{Regions{Segment{3, 6}, Segment{13, 16}}, Tail(-7), Segment{2, 2}},
	{Regions{Segment{13, 16}, Segment{3, 6}}, Head(+0), Segment{13, 13}},
	{Regions{Segment{13, 16}, Segment{3, 6}}, Head(+7), Segment{7, 7}},
	{Regions{Segment{13, 16}, Segment{3, 6}}, Tail(+0), Segment{6, 6}},
	{Regions{Segment{13, 16}, Segment{3, 6}}, Tail(-7), Segment{12, 12}},

	{Regions{Segment{3, 6}, Segment{13, 16}}, HeadTail{0, 0}, Regions{Segment{3, 6}, Segment{13, 16}}},
	{Regions{Segment{3, 6}, Segment{13, 16}}, HeadTail{4, -4}, Segment{14, 14}},
	{Regions{Segment{3, 6}, Segment{13, 16}}, HeadHead{-2, 4}, Regions{Segment{1, 6}, Segment{13, 14}}},
	{Regions{Segment{3, 6}, Segment{13, 16}}, TailTail{-4, 2}, Regions{Segment{5, 6}, Segment{13, 18}}},
}

func TestRegionResize(t *testing.T) {
	for i, tt := range regionResizeTests {
		out := tt.in.Resize(tt.modifier)
		if !reflect.DeepEqual(out, tt.out) {
			t.Errorf(
				"\ncase [%d]:\n  resize by %s\n   in: %#v\n  out: %#v\n  exp: %#v",
				i+1, tt.modifier, tt.in, out, tt.out,
			)
		}
	}
}

var regionWithinTests = []struct {
	in   Region
	l, u int
	out  bool
}{
	{Segment{3, 6}, 2, 7, true},
	{Segment{3, 6}, 3, 6, true},
	{Segment{3, 6}, 4, 6, false},
	{Segment{3, 6}, 3, 5, false},
	{Segment{3, 6}, 4, 5, false},
	{Segment{3, 6}, 6, 3, false},

	{Segment{6, 3}, 2, 7, true},
	{Segment{6, 3}, 3, 6, true},
	{Segment{6, 3}, 4, 6, false},
	{Segment{6, 3}, 3, 5, false},
	{Segment{6, 3}, 4, 5, false},
	{Segment{6, 3}, 6, 3, false},

	{Regions{Segment{3, 6}, Segment{13, 16}}, 3, 16, true},
	{Regions{Segment{3, 6}, Segment{13, 16}}, 3, 6, false},
	{Regions{Segment{3, 6}, Segment{13, 16}}, 13, 16, false},
}

func TestRegionWithin(t *testing.T) {
	for _, tt := range regionWithinTests {
		out := tt.in.Within(tt.l, tt.u)
		if out != tt.out {
			t.Errorf("%v.Within(%d, %d) = %t, want %t", tt.in, tt.l, tt.u, out, tt.out)
		}
	}
}

var regionLocateTests = []struct {
	in  Region
	out Sequence
}{
	{Segment{2, 6}, New(nil, nil, []byte("gcat"))},
	{Segment{6, 2}, New(nil, nil, []byte("atgc"))},
	{Regions{Segment{0, 2}, Segment{4, 6}}, New(nil, nil, []byte("atat"))},
}

func TestRegionLocate(t *testing.T) {
	seq := New(nil, nil, []byte("atgcatgc"))
	for _, tt := range regionLocateTests {
		out, exp := tt.in.Locate(seq), tt.out
		if !reflect.DeepEqual(out.Info(), exp.Info()) {
			t.Errorf("Slice(in, %d, %d).Info() = %v, want %v", 2, 6, out.Info(), exp.Info())
		}
		if diff := deep.Equal(out.Features(), exp.Features()); diff != nil {
			t.Errorf("Slice(in, %d, %d).Features() = %v, want %v", 2, 6, out.Features(), exp.Features())
		}
		if diff := deep.Equal(out.Bytes(), exp.Bytes()); diff != nil {
			t.Errorf("Slice(in, %d, %d).Bytes() = %v, want %v", 2, 6, out.Bytes(), exp.Bytes())
		}

		cmp := tt.in.Complement()
		if cmp.Len() != tt.in.Len() {
			t.Errorf("%s.Len() = %d, want %d", cmp, cmp.Len(), tt.in.Len())
		}
		if !reflect.DeepEqual(cmp.Complement(), tt.in) {
			t.Errorf(
				"%s.Complement() = %s, want %s",
				cmp, cmp.Complement(), tt.in,
			)
		}
		out = cmp.Locate(seq)
		exp = Reverse(Complement(tt.out))
		if !Equal(out, exp) {
			t.Errorf(
				"%s.Locate(%q) = %q, want %q",
				cmp, string(seq.Bytes()),
				string(out.Bytes()), string(exp.Bytes()),
			)
		}
	}
}

func TestMinimize(t *testing.T) {
	in := Regions{Segment{1, 3}, Segment{6, 9}, Segment{5, 3}, Segment{6, 8}, Segment{1, 3}}
	exp := []Segment{{1, 5}, {6, 9}}
	out := Minimize(in)
	if !reflect.DeepEqual(out, exp) {
		t.Errorf("Minimize(%#v) = %#v, want %#v", in, out, exp)
	}
}
