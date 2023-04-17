package toml

import "testing"

func TestTomlDecode(t *testing.T) {
	var config struct {
		Middlewares []struct {
			Name    string `toml:"name"`
			Options Value  `toml:"options"`
		} `toml:"middlewares"`
	}

	md, err := Decode(`
	middlewares = [
  		{name = "ratelimit", options = {max = 10, period = 10}},
	]
	`, &config)
	if err != nil {
		t.Error(err)
		return
	}

	var rateLimitConfig struct {
		Max    int `toml:"max"`
		Period int `toml:"period"`
	}
	err = md.PrimitiveDecode(config.Middlewares[0].Options, &rateLimitConfig)
	if err != nil {
		t.Error(err)
		return
	}
	if rateLimitConfig.Max != 10 || rateLimitConfig.Period != 10 {
		t.Errorf("Expected {Max: 10, Period: 10}, got %v", rateLimitConfig)
	}
}
