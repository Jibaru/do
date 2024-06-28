package utils

import (
	"reflect"
	"testing"
)

func TestPtr_int(t *testing.T) {
	expected := 1
	actual := Ptr(expected)

	if expected != *actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestPtr_string(t *testing.T) {
	expected := "hello"
	actual := Ptr(expected)

	if expected != *actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestPtr_bool(t *testing.T) {
	actual := Ptr(true)

	if true != *actual {
		t.Errorf("expected %v, got %v", true, actual)
	}
}

func TestPtr_struct(t *testing.T) {
	type testStruct struct {
		Field1 string
		Field2 int
	}

	expected := testStruct{
		Field1: "hello",
		Field2: 1,
	}
	actual := Ptr(expected)

	if expected != *actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestPtr_slice(t *testing.T) {
	expected := []string{"hello", "world"}
	actual := Ptr(expected)

	if &expected[0] != &(*actual)[0] {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestPtr_map(t *testing.T) {
	expected := map[string]string{
		"hello": "world",
	}
	actual := Ptr(expected)

	if !reflect.DeepEqual(expected, *actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
