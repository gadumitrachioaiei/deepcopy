// this file is copied from github.com/mitchellh/copystructure
// so we can test our own code

//The MIT License (MIT)
//
//Copyright (c) 2014 Mitchell Hashimoto
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in
//all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//THE SOFTWARE.

package deepcopy

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestCopy_complex(t *testing.T) {
	v := map[string]interface{}{
		"foo": []string{"a", "b"},
		"bar": "baz",
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_interfacePointer(t *testing.T) {
	type Nested struct {
		Field string
	}

	type Test struct {
		Value *interface{}
	}

	ifacePtr := func(v interface{}) *interface{} {
		return &v
	}

	v := Test{
		Value: ifacePtr(Nested{Field: "111"}),
	}
	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_primitive(t *testing.T) {
	cases := []interface{}{
		42,
		"foo",
		1.2,
	}

	for _, tc := range cases {
		result, err := Copy(tc)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if result != tc {
			t.Fatalf("bad: %#v", result)
		}
	}
}

func TestCopy_primitivePtr(t *testing.T) {
	i := 42
	s := "foo"
	f := 1.2
	cases := []interface{}{
		&i,
		&s,
		&f,
	}

	for i, tc := range cases {
		result, err := Copy(tc)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		if !reflect.DeepEqual(result, tc) {
			t.Fatalf("%d exptected: %#v\nbad: %#v", i, tc, result)
		}
	}
}

func TestCopy_map(t *testing.T) {
	v := map[string]interface{}{
		"bar": "baz",
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_array(t *testing.T) {
	v := [2]string{"bar", "baz"}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_pointerToArray(t *testing.T) {
	v := &[2]string{"bar", "baz"}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_slice(t *testing.T) {
	v := []string{"bar", "baz"}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_pointerToSlice(t *testing.T) {
	v := &[]string{"bar", "baz"}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_pointerToMap(t *testing.T) {
	v := &map[string]string{"bar": "baz"}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_struct(t *testing.T) {
	type test struct {
		Value string
	}

	v := test{Value: "foo"}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_structPtr(t *testing.T) {
	type test struct {
		Value string
	}

	v := &test{Value: "foo"}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_structNil(t *testing.T) {
	type test struct {
		Value string
	}

	var v *test
	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if v, ok := result.(*test); !ok {
		t.Fatalf("bad: %#v", result)
	} else if v != nil {
		t.Fatalf("bad: %#v", v)
	}
}

func TestCopy_structNested(t *testing.T) {
	type TestInner struct{}

	type Test struct {
		Test *TestInner
	}

	v := Test{}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_structWithNestedArray(t *testing.T) {
	type TestInner struct {
		Value string
	}

	type Test struct {
		Value [2]TestInner
	}

	v := Test{
		Value: [2]TestInner{
			{Value: "bar"},
			{Value: "baz"},
		},
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_structWithPointerToSliceField(t *testing.T) {
	type Test struct {
		Value *[]string
	}

	v := Test{
		Value: &[]string{"bar", "baz"},
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_structWithPointerToArrayField(t *testing.T) {
	type Test struct {
		Value *[2]string
	}

	v := Test{
		Value: &[2]string{"bar", "baz"},
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_structWithPointerToMapField(t *testing.T) {
	type Test struct {
		Value *map[string]string
	}

	v := Test{
		Value: &map[string]string{"bar": "baz"},
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_structUnexported(t *testing.T) {
	type test struct {
		Value string

		private string
	}

	v := test{Value: "foo"}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_structUnexportedMap(t *testing.T) {
	type Sub struct {
		Foo map[string]interface{}
	}

	type test struct {
		Value string

		private Sub
	}

	v := test{
		Value: "foo",
		private: Sub{
			Foo: map[string]interface{}{
				"yo": 42,
			},
		},
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// private should not be copied
	v.private = Sub{}
	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad:\n\n%#v\n\n%#v", result, v)
	}
}

func TestCopy_structUnexportedArray(t *testing.T) {
	type Sub struct {
		Foo [2]string
	}

	type test struct {
		Value string

		private Sub
	}

	v := test{
		Value: "foo",
		private: Sub{
			Foo: [2]string{"bar", "baz"},
		},
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// private should not be copied
	v.private = Sub{}
	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad:\n\n%#v\n\n%#v", result, v)
	}
}

// This is testing an unexported field containing a slice of pointers, which
// was a crashing case found in Terraform.
func TestCopy_structUnexportedPtrMap(t *testing.T) {
	type Foo interface{}

	type Sub struct {
		List []Foo
	}

	type test struct {
		Value string

		private *Sub
	}

	v := test{
		Value: "foo",
		private: &Sub{
			List: []Foo{&Sub{}},
		},
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// private should not be copied
	v.private = nil
	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad:\n\n%#v\n\n%#v", result, v)
	}
}

func TestCopy_nestedStructUnexported(t *testing.T) {
	type subTest struct {
		mine string
	}

	type test struct {
		Value   string
		private subTest
	}

	v := test{Value: "foo"}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_time(t *testing.T) {
	type test struct {
		Value time.Time
	}

	v := test{Value: time.Now().UTC()}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

func TestCopy_aliased(t *testing.T) {
	type (
		Int   int
		Str   string
		Map   map[Int]interface{}
		Slice []Str
	)

	v := Map{
		1: Map{10: 20},
		2: Map(nil),
		3: Slice{"a", "b"},
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("bad: %#v", result)
	}
}

type EmbeddedLocker struct {
	sync.Mutex
	Map map[int]int
}

type LockedField struct {
	String string
	Locker *EmbeddedLocker
	// this should not get locked or have its state copied
	Mutex    sync.Mutex
	nilMutex *sync.Mutex
}

// test something that doesn't contain a lock internally
type lockedMap map[int]int

var mapLock sync.Mutex

func (m lockedMap) Lock()   { mapLock.Lock() }
func (m lockedMap) Unlock() { mapLock.Unlock() }

// Use an RLock if available
type RLocker struct {
	sync.RWMutex
	Map map[int]int
}

type PointerLocker struct {
	Mu sync.Mutex
}

func (p *PointerLocker) Lock()   { p.Mu.Lock() }
func (p *PointerLocker) Unlock() { p.Mu.Unlock() }

func TestCopy_sliceWithNil(t *testing.T) {
	v := [](*int){nil}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("expected:\n%#v\ngot:\n%#v", v, result)
	}
}

func TestCopy_mapWithNil(t *testing.T) {
	v := map[int](*int){0: nil}

	result, err := Copy(v)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(result, v) {
		t.Fatalf("expected:\n%#v\ngot:\n%#v", v, result)
	}
}

func TestCopy_mapWithPointers(t *testing.T) {
	type T struct {
		S string
	}
	v := map[string]interface{}{
		"a": &T{S: "hello"},
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(v, result) {
		t.Fatalf("%#v", result)
	}
}

func TestCopy_structWithMapWithPointers(t *testing.T) {
	type T struct {
		S string
		M map[string]interface{}
	}
	v := &T{
		S: "a",
		M: map[string]interface{}{
			"b": &T{
				S: "b",
			},
		},
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(v, result) {
		t.Fatal(result)
	}
}

type testT struct {
	N   int
	Spp **string
	X   testX
	Xp  *testX
	Xpp **testX
}

type testX struct {
	Tp  *testT
	Tpp **testT
	Ip  *interface{}
	Ep  *error
	S   fmt.Stringer
}

type stringer struct{}

func (s *stringer) String() string {
	return "test string"
}

func TestCopy_structWithPointersAndInterfaces(t *testing.T) {
	// test that we can copy various nested and chained pointers and interfaces
	s := "val"
	sp := &s
	spp := &sp
	i := interface{}(11)

	tp := &testT{
		N: 2,
	}

	xp := &testX{
		Tp:  tp,
		Tpp: &tp,
		Ip:  &i,
		S:   &stringer{},
	}

	v := &testT{
		N:   1,
		Spp: spp,
		X:   testX{},
		Xp:  xp,
		Xpp: &xp,
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(v, result) {
		t.Fatal(result)
	}
}

func Test_pointerInterfacePointer(t *testing.T) {
	s := "hi"
	si := interface{}(&s)
	sip := &si

	result, err := Copy(sip)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(sip, result) {
		t.Fatalf("%#v != %#v\n", sip, result)
	}
}

func Test_pointerInterfacePointer2(t *testing.T) {
	type T struct {
		I *interface{}
		J **fmt.Stringer
	}

	x := 1
	y := &stringer{}

	i := interface{}(&x)
	j := fmt.Stringer(y)
	jp := &j

	v := &T{
		I: &i,
		J: &jp,
	}
	result, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(v, result) {
		t.Fatalf("%#v != %#v\n", v, result)
	}
}

// This test catches a bug that happened when unexported fields were
// first their subsequent fields wouldn't be copied.
func TestCopy_unexportedFieldFirst(t *testing.T) {
	type P struct {
		mu       sync.Mutex
		Old, New string
	}

	type T struct {
		M map[string]*P
	}

	v := &T{
		M: map[string]*P{
			"a": {Old: "", New: "2"},
		},
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(v, result) {
		t.Fatalf("\n%#v\n\n%#v", v, result)
	}
}

func TestCopy_nilPointerInSlice(t *testing.T) {
	type T struct {
		Ps []*int
	}

	v := &T{
		Ps: []*int{nil},
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(v, result) {
		t.Fatalf("\n%#v\n\n%#v", v, result)
	}
}

//-------------------------------------------------------------------
// The tests below all tests various pointer cases around copying
// a structure that uses a defined Copier. This was originally raised
// around issue #26.

func TestCopy_timePointer(t *testing.T) {
	type T struct {
		Value *time.Time
	}

	now := time.Now()
	v := &T{
		Value: &now,
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(v, result) {
		t.Fatalf("\n%#v\n\n%#v", v, result)
	}
}

func TestCopy_timeNonPointer(t *testing.T) {
	type T struct {
		Value time.Time
	}

	v := &T{
		Value: time.Now(),
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(v, result) {
		t.Fatalf("\n%#v\n\n%#v", v, result)
	}
}

func TestCopy_timeDoublePointer(t *testing.T) {
	type T struct {
		Value **time.Time
	}

	now := time.Now()
	nowP := &now
	nowPP := &nowP
	v := &T{
		Value: nowPP,
	}

	result, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(v, result) {
		t.Fatalf("\n%#v\n\n%#v", v, result)
	}
}
