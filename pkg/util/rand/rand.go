package rand

import (
	"bytes"
	"crypto/rand"
	"errors"
)

// define a flag that generates a random string
const (
	Ldigit = 1 << iota
	LlowerCase
	LupperCase
	LlowerAndUpperCase = LlowerCase | LupperCase
	LdigitAndLowerCase = Ldigit | LlowerCase
	LdigitAndUpperCase = Ldigit | LupperCase
	LdigitAndLetter    = Ldigit | LlowerCase | LupperCase
)

var (
	digits           = []byte("0123456789")
	lowerCaseLetters = []byte("abcdefghijklmnopqrstuvwxyz")
	upperCaseLetters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// definition error
var (
	ErrInvalidFlag = errors.New("Invalid flag")
)

// Random generate a random string specifying the length of the random number
// and the random flag
func Random(length, flag int) (string, error) {
	if length < 1 {
		length = 6
	}

	source, err := getFlagSource(flag)
	if err != nil {
		return "", err
	}

	b, err := randomBytesMod(length, byte(len(source)))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	for _, c := range b {
		buf.WriteByte(source[c])
	}

	return buf.String(), nil
}

func getFlagSource(flag int) ([]byte, error) {
	var source []byte

	if flag&Ldigit > 0 {
		source = append(source, digits...)
	}

	if flag&LlowerCase > 0 {
		source = append(source, lowerCaseLetters...)
	}

	if flag&LupperCase > 0 {
		source = append(source, upperCaseLetters...)
	}

	sourceLen := len(source)
	if sourceLen == 0 {
		return nil, ErrInvalidFlag
	}
	return source, nil
}

func randomBytesMod(length int, mod byte) ([]byte, error) {
	b := make([]byte, length)
	max := 255 - 255%mod
	i := 0

LROOT:
	for {
		r, err := randomBytes(length + length/4)
		if err != nil {
			return nil, err
		}

		for _, c := range r {
			if c >= max {
				// Skip this number to avoid modulo bias
				continue
			}

			b[i] = c % mod
			i++
			if i == length {
				break LROOT
			}
		}

	}

	return b, nil
}

func randomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
