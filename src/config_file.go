package src

import (
	"fmt"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

const ConfigFileName = ".dump_dir.yml"

type ConfigFile struct {
	Include []string `yaml:"include,omitempty"`
	Ignore  []string `yaml:"ignore,omitempty"`
}

type ConfigLoader struct {
	fs afero.Fs
}

func NewConfigLoader(fs afero.Fs) *ConfigLoader {
	return &ConfigLoader{fs: fs}
}

func (cl *ConfigLoader) LoadAndMergeConfig(cmdConfig Config) (Config, error) {
	if cmdConfig.NoConfig {
		return cmdConfig, nil
	}
	configFile, err := cl.loadConfigFile()
	if err != nil {
		return cmdConfig, fmt.Errorf("error loading config file: %w", err)
	}
	if configFile == nil {
		return cmdConfig, nil
	}
	return MergeConfigs(cmdConfig, *configFile), nil
}

func MergeConfigs(cmdConfig Config, fileConfig ConfigFile) Config {
	mergedConfig := cmdConfig

	if fileConfig.Ignore != nil {
		for _, ignorePath := range fileConfig.Ignore {
			mergedConfig.AddSkipDir(ignorePath)
		}
	}

	if fileConfig.Include != nil {
		for _, includePath := range fileConfig.Include {
			if err := mergedConfig.AddIncludePath(includePath); err != nil {
				fmt.Printf("Warning: Could not process path %s: %v\n", includePath, err)
			}
		}
	}

	return mergedConfig
}

func (cl *ConfigLoader) loadConfigFile() (*ConfigFile, error) {
	exists, err := afero.Exists(cl.fs, ConfigFileName)
	if err != nil {
		return nil, fmt.Errorf("error checking config file existence: %w", err)
	}

	if !exists {
		return nil, nil
	}

	data, err := afero.ReadFile(cl.fs, ConfigFileName)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config ConfigFile
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}
