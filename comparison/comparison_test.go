package comparison

import (
	"reflect"
	"testing"
)

func TestSquare(t *testing.T) {
	arg := 4
	want := 16
	got := Square(arg)
	if got != want {
		t.Errorf("Square(%d) = %d; want %d", arg, got, want)
	}
}

// Although the objects point to different locations
// they are being compared based on field matching,
// not memory address, so the comparison test passes.
func TestPerson(t *testing.T) {
	simon := Person{
		Name: "Simon",
		Age:  23,
	}
	simonClone := Person{
		Name: "Simon",
		Age:  23,
	}
	if simon != simonClone {
		t.Errorf("simon != simonClone")
	}
}

// Here, two pointers are being directly compared,
// so the test will pass if the two objects have
// different memory addresses.
func TestPersonPtr(t *testing.T) {
	simon := &Person{}
	simonClone := &Person{}
	t.Logf("simonPtr=%p, simonClonePtr=%p", simon, simonClone)
	if simon == simonClone {
		t.Errorf("simon == simonClone")
	}
}

// In the case of comparing objects with nested functions,
// or the fields of pointers to objects, the reflect package
// is useful. See godoc reflect.
func TestStructWithFun(t *testing.T) {
	test := StructWithFun{}
	testClone := StructWithFun{}
	if !reflect.DeepEqual(test, testClone) {
		t.Errorf("test != testClone")
	}

	testPtr := &StructWithFun{}
	testPtrClone := &StructWithFun{}
	if !reflect.DeepEqual(testPtr, testPtrClone) {
		t.Errorf("testPtr != testPtrClone")
	}
}
