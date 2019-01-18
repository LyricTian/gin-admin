package util

import "testing"

func TestFillStruct(t *testing.T) {
	type Foo struct {
		Bar  string
		FBar string
	}

	type Foo2 struct {
		Bar string
	}

	source := Foo{
		Bar:  "bar",
		FBar: "fbar",
	}
	target := &Foo2{}

	err := FillStruct(source, target)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if target.Bar != source.Bar {
		t.Error("not the expected value:", target)
	}
}

func TestFillStructs(t *testing.T) {
	type Foo struct {
		Bar  string
		FBar string
	}

	source := []Foo{
		{Bar: "bar1", FBar: "fbar1"},
	}

	var target []Foo
	err := FillStructs(source, &target)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(target) != len(source) ||
		target[0].Bar != source[0].Bar ||
		target[0].FBar != source[0].FBar {
		t.Error("not the expected value:", target)
	}
}
