package toml

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

var (
	Unmarshal  = toml.Unmarshal
	DecodeFile = toml.DecodeFile
	Decode     = toml.Decode
)

type Value = toml.Primitive

func Marshal(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MarshalToString(v interface{}) (string, error) {
	b, err := Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
