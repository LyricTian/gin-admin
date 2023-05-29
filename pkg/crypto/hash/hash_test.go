package hash

import (
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	origin := "abc-123"
	hashPwd, err := GeneratePassword(origin)
	if err != nil {
		t.Error("GeneratePassword Failed: ", err.Error())
	}
	t.Log("test password: ", hashPwd, ",length: ", len(hashPwd))

	if err := CompareHashAndPassword(hashPwd, origin); err != nil {
		t.Error("Unmatched password: ", err.Error())
	}
}

func TestMD5(t *testing.T) {
	origin := "abc-123"
	hashVal := "6351623c8cef86fefabfa7da046fc619"
	if v := MD5String(origin); v != hashVal {
		t.Error("Failed to generate MD5 hash: ", v)
	}
}
