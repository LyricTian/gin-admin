package aes

import "testing"

func TestEncrypt(t *testing.T) {
	origData := "root@bobcatminer.com"
	encryptData, err := EncryptToBase64([]byte(origData), SecretKey)
	if err != nil {
		t.Error("AES encrypt error: ", err.Error())
		return
	}
	t.Logf("encrypt data: %s", encryptData)

	decryptData, err := DecryptFromBase64(encryptData, SecretKey)
	if err != nil {
		t.Error("AES decrypt error: ", err.Error())
		return
	}

	if string(decryptData) != origData {
		t.Errorf("Not expect result,want:%s,given:%s", origData, decryptData)
	}
}
