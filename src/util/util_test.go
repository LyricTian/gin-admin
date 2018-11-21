package util

import "testing"

func TestGetLevelCode(t *testing.T) {
	levelCodes := []string{
		"01",
		"0101",
		"0102",
	}
	levelCode := GetLevelCode(levelCodes)
	if levelCode != "0103" {
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
