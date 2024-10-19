//  SPDX-FileCopyrightText: 2024 Paulo Almeida <almeidapaulopt@gmail.com>
//  SPDX-License-Identifier: MIT

package core

import (
	"fmt"
	"strings"

	"github.com/creasty/defaults"

	"github.com/knadh/koanf/v2"

	"github.com/knadh/koanf/providers/env"
)

const prefix = "TSDPROXY_"

type (
	// Config stores complete configuration.
	Config struct {
		PublicURL string `default:"http://localhost:8080"`
		DataDir   string `default:"/data/"`

		Log  LogConfig
		HTTP HTTPConfig
	}

	// LogConfig stores logging configuration.
	LogConfig struct {
		Level string `default:"info"`
		JSON  bool   `default:"false"`
	}

	// HTTPConfig stores HTTP configuration.
	HTTPConfig struct {
		Hostname string `default:"0.0.0.0"`
		Port     uint16 `default:"8080"`
	}
)

// GetConfig loads and returns configuration.
func GetConfig() (*Config, error) {
	c := new(Config)

	// load default values
	if err := defaults.Set(c); err != nil {
		fmt.Printf("Error loading defaults: %v", err)
	}

	k := koanf.New(".")

	// load environment variables
	err := k.Load(
		env.Provider(
			prefix,
			".",
			func(s string) string {
				return strings.Replace(
					strings.ToLower(strings.TrimPrefix(s, prefix)),
					"_",
					".",
					-1)
			},
		),
		nil,
	)
	if err != nil {
		fmt.Printf("Error loading env: %v", err)
	}

	// unmarshal config to struct
	err = k.UnmarshalWithConf("", &c, koanf.UnmarshalConf{
		Tag: "env",
	})
	if err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	return c, nil
}
