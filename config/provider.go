// Package config - custom `choice` flag from: https://github.com/spf13/pflag/issues/236#issuecomment-931600452
package config

import (
	"fmt"
	"strings"
)

// ProviderConfig is a "choice" style flag option
type ProviderConfig struct {
	Allowed []string
	Value   string
}

func NewProviderConfig(allowed []string, v string) *ProviderConfig {
	return &ProviderConfig{
		Allowed: allowed,
		Value:   v,
	}
}

func (pc *ProviderConfig) String() string {
	return pc.Value
}

func (pc *ProviderConfig) Set(v string) error {
	isIncluded := func(opts []string, val string) bool {
		for _, opt := range opts {
			if val == opt {
				return true
			}
		}
		return false
	}
	if !isIncluded(pc.Allowed, v) {
		return fmt.Errorf("%s is not included in [%s]", v, strings.Join(pc.Allowed, ", "))
	}
	pc.Value = v
	return nil
}

func (pc *ProviderConfig) Type() string {
	return "string"
}
