package config

import (
	"bytes"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestWithToml(t *testing.T) {
	log := Log{
		Level: 1,
	}

	log.Hooks = append(log.Hooks, LogHook{
		Type:      "logkit",
		Levels:    []string{"info", "warn", "error"},
		MaxBuffer: 1024,
		MaxThread: 2,
	})
	C.Log = log

	buf := new(bytes.Buffer)
	toml.NewEncoder(buf).Encode(C)
	t.Log(buf.String())
}
