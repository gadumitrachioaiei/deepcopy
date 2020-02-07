package deepcopy

import (
	"encoding/gob"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/mitchellh/copystructure"
)

func TestCopyNilMap(t *testing.T) {
	var x map[int]int
	vi, err := Copy(x)
	if err != nil {
		t.Fatal(err)
	}
	v := vi.(map[int]int)
	if v != nil {
		t.Fatalf("got value: %v, expected nil", v)
	}
}

func TestCopyNilSlice(t *testing.T) {
	var x []int
	vi, err := Copy(x)
	if err != nil {
		t.Fatal(err)
	}
	v := vi.([]int)
	if v != nil {
		t.Fatalf("got value: %v, expected nil", v)
	}
}

func TestCopyNilPointer(t *testing.T) {
	var x *int
	vi, err := Copy(x)
	if err != nil {
		t.Fatal(err)
	}
	v := vi.(*int)
	if v != nil {
		t.Fatalf("got value: %v, expected nil", v)
	}
}

func TestCopyNilInterface(t *testing.T) {
	var y *int
	var x interface{} = y
	vi, err := Copy(x)
	if err != nil {
		t.Fatal(err)
	}
	v := vi.(*int)
	if v != nil {
		t.Fatalf("got value: %v, expected nil", v)
	}
}

func TestCopyInterface(t *testing.T) {
	data := []byte("interface test")
	u := struct {
		R io.Reader
	}{
		R: &R{Data: data},
	}
	vi, err := Copy(u)
	if err != nil {
		t.Fatal(err)
	}
	v := vi.(struct{ R io.Reader })
	if diff := cmp.Diff(u, v); diff != "" {
		t.Fatal(diff)
	}
	// test that we can read from v
	dv, err := ioutil.ReadAll(v.R)
	if err != nil {
		t.Fatal(err)
	}
	if string(dv) != string(data) {
		t.Fatalf("got: %s, expected: %s", dv, data)
	}
	// test that we can read from u
	du, err := ioutil.ReadAll(u.R)
	if err != nil {
		t.Fatal(err)
	}
	if string(du) != string(data) {
		t.Fatalf("got: %s, expected: %s", du, data)
	}
}

func TestCopyEquality(t *testing.T) {
	type A struct {
		B []byte
		C int
		D time.Time
	}
	type T struct {
		A int
		B string
		C []int
		D map[int]int
		E map[int][]int
		F *[]int
		G *A
		R
	}
	u := T{}
	f := fuzz.New()
	f.Fuzz(&u)
	v, err := Copy(u)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(u, v); diff != "" {
		t.Fatal(diff)
	}
}

var v interface{}
var vt T

func BenchmarkCopy(b *testing.B) {
	u := T{}
	f := fuzz.New()
	f.Fuzz(&u)
	var err error
	b.ResetTimer()
	b.Run("deepcopy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v, err = Copy(u)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("copystructure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v, err = copystructure.Copy(u)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("gob ser/deser", func(b *testing.B) {
		r, w := io.Pipe()
		ge := gob.NewEncoder(w)
		gd := gob.NewDecoder(r)
		var vT T
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			go func() {
				if err := ge.Encode(u); err != nil {
					b.Fatal(err)
				}
			}()
			if err := gd.Decode(&vT); err != nil {
				b.Fatal(err)
			}
		}
	})
}

type R struct {
	Data []byte
	Off  int
}

func (r *R) Read(p []byte) (int, error) {
	if r.Off > len(r.Data)-1 {
		return 0, io.EOF
	}
	n := copy(p, r.Data[r.Off:])
	r.Off += n
	return n, nil
}

type A struct {
	B []byte
	C int
	D time.Time
	b []byte
}
type T struct {
	A int
	B string
	C []int
	D map[int]int
	E map[int][]int
	F *[]int
	G A
}
