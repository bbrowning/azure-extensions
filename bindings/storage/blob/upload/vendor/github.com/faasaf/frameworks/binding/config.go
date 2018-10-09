package binding

import (
	"fmt"
	"strings"
)

type Config struct {
	cfgStore map[string]string
}

func newConfig(pairs []string) (Config, error) {
	cfg := Config{
		cfgStore: map[string]string{},
	}
	for _, pair := range pairs {
		tokens := strings.SplitN(pair, "=", 2)
		if len(tokens) != 2 {
			return cfg, fmt.Errorf("error parsing setting: %s", pair)
		}
		cfg.cfgStore[tokens[0]] = tokens[1]
	}
	return cfg, nil
}

func (c Config) GetSetting(key, dflt string) string {
	val, ok := c.cfgStore[key]
	if !ok {
		return dflt
	}
	return val
}
