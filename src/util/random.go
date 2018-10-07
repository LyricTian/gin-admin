package util

import (
	"io"
	"math/rand"
	"time"
)

// NewRandom create the Random instances
// l the length of the generated random code
func NewRandom(l int) *Random {
	return &Random{vl: l}
}

// Random provide random code generation
type Random struct {
	vl int
}

// Number generate random code contains only Numbers
func (rd *Random) Number() string {
	source := rd.number()
	return rd.Source(source)
}

// LowerLetter the random code generation only contain lowercase letters
func (rd *Random) LowerLetter() string {
	source := rd.lowerLetter()
	return rd.Source(source)
}

// UpperLetter the random code generation only contains uppercase letters
func (rd *Random) UpperLetter() string {
	source := rd.upperLetter()
	return rd.Source(source)
}

// NumberAndUpperLetter 随机数包含数字和大写字母
func (rd *Random) NumberAndUpperLetter() string {
	source := rd.number()
	source = append(source, rd.upperLetter()...)
	return rd.Source(source)
}

// NumberAndLetter generated contains Numbers and letters (case-insensitive) random code
func (rd *Random) NumberAndLetter() string {
	source := rd.number()
	source = append(source, rd.lowerLetter()...)
	source = append(source, rd.upperLetter()...)
	return rd.Source(source)
}

// Source from the specified data source to generate random codes
func (rd *Random) Source(source []byte) string {
	if len(source) == 0 {
		return ""
	}

	rdn := rand.New(rand.NewSource(time.Now().UnixNano()))
	r, w := io.Pipe()
	go func() {
		for i := 0; i < rd.vl; i++ {
			defer func() {
				if err := w.Close(); err != nil {
					panic(err)
				}
			}()

			val := source[rdn.Intn(len(source))]
			_, err := w.Write([]byte{val})
			if err != nil {
				panic(err)
			}
		}
	}()

	var result []byte
	for {
		buf := make([]byte, rd.vl)
		n, err := r.Read(buf)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		result = append(result, buf[:n]...)
	}
	return string(result)
}

func (rd *Random) number() []byte {
	v := make([]byte, 10)
	for i, j := 48, 0; i <= 57; i, j = i+1, j+1 {
		v[j] = byte(i)
	}
	return v
}

func (rd *Random) lowerLetter() []byte {
	v := make([]byte, 26)
	for i, j := 97, 0; i < 123; i, j = i+1, j+1 {
		v[j] = byte(i)
	}
	return v
}

func (rd *Random) upperLetter() []byte {
	v := make([]byte, 26)
	for i, j := 65, 0; i < 91; i, j = i+1, j+1 {
		v[j] = byte(i)
	}
	return v
}
