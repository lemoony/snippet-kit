package config

import (
	"emperror.dev/errors"
	"gopkg.in/yaml.v3"

	"github.com/lemoony/snipkit/internal/assistant/openai"
)

const (
	VersionFrom = "1.1.1"
	VersionTo   = "1.1.2"
)

func Migrate(old []byte) []byte {
	var config versionWrapper

	if err := yaml.Unmarshal(old, &config); err != nil {
		panic(err)
	}

	if config.Version != VersionFrom {
		panic(errors.Errorf("Invalid version for migration to v1.1.2: %s", config.Version))
	}

	config.Version = VersionTo
	config.Config.Assistant = make(map[string]interface{})

	config.Config.Assistant["openai"] = openai.AutoDiscoveryConfig()

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	return configBytes
}

type versionWrapper struct {
	Version string     `yaml:"version"`
	Config  configV112 `yaml:"config"`
}

type configV112 struct {
	Style              map[string]interface{} `yaml:"style"`
	Editor             string                 `yaml:"editor"`
	DefaultRootCommand string                 `yaml:"defaultRootCommand"`
	FuzzySearch        bool                   `yaml:"fuzzySearch"`
	Script             map[string]interface{} `yaml:"scripts"`
	Assistant          map[string]interface{} `yaml:"assistant"`
	Manager            map[string]interface{} `yaml:"manager"`
}
