package json

import (
	"os"

	jsoniter "github.com/json-iterator/go"
)

// Define alias
var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)

func MarshalToString(v interface{}) string {
	s, err := jsoniter.MarshalToString(v)
	if err != nil {
		os.Stderr.WriteString("[warning] jsoniter marshal to string failed: " + err.Error())
		return ""
	}

	return s
}
