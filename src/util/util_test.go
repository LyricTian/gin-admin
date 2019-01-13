package util

import "testing"

func TestGetLevelCode(t *testing.T) {
	levelCodes := []string{
		"01",
		"0102",
		"0103",
	}
	levelCode := GetLevelCode(levelCodes)
	if levelCode != "0101" {
		t.Error("无效的分级码：", levelCode)
		return
	}
}

func TestParseLevelCodes(t *testing.T) {
	codes := ParseLevelCodes("1010", "101010")

	if len(codes) != 3 ||
		codes[0] != "10" ||
		codes[1] != "1010" ||
		codes[2] != "101010" {
		t.Error("分级码解析错误：", codes)
	}
}

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

	target := make([]Foo, len(source))

	err := FillStructs(source, target)
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
