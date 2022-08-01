package config

import (
	"testing"

	"github.com/LyricTian/gin-admin/v9/pkg/util/toml"
)

func TestConfig(t *testing.T) {
	MustLoad("../../configs/config.toml")

	s, err := toml.MarshalToString(C)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", s, "\n")
}
