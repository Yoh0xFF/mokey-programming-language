package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &StringObject{Value: "Hello World"}
	hello2 := &StringObject{Value: "Hello World"}
	diff1 := &StringObject{Value: "My name is johnny"}
	diff2 := &StringObject{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &BoolObject{Value: true}
	true2 := &BoolObject{Value: true}
	false1 := &BoolObject{Value: false}
	false2 := &BoolObject{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("trues do not have same hash key")
	}

	if false1.HashKey() != false2.HashKey() {
		t.Errorf("falses do not have same hash key")
	}

	if true1.HashKey() == false1.HashKey() {
		t.Errorf("true has same hash key as false")
	}
}

func TestIntegerHashKey(t *testing.T) {
	one1 := &IntObject{Value: 1}
	one2 := &IntObject{Value: 1}
	two1 := &IntObject{Value: 2}
	two2 := &IntObject{Value: 2}

	if one1.HashKey() != one2.HashKey() {
		t.Errorf("integers with same content have twoerent hash keys")
	}

	if two1.HashKey() != two2.HashKey() {
		t.Errorf("integers with same content have twoerent hash keys")
	}

	if one1.HashKey() == two1.HashKey() {
		t.Errorf("integers with twoerent content have same hash keys")
	}
}
