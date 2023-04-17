package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// md5 hash
func MD5(b []byte) string {
	h := md5.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// md5 hash
func MD5String(s string) string {
	return MD5([]byte(s))
}

// sha1 hash
func SHA1(b []byte) string {
	h := sha1.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// sha1 hash
func SHA1String(s string) string {
	return SHA1([]byte(s))
}

// Use bcrypt generate password hash
func GeneratePassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Use bcrypt compare hash password and password
func CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
