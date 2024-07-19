package tests

import . "github.com/fargusplumdoodle/dump_dir/src"

// ConfigOption is a function type that modifies a Config
type ConfigOption func(*Config)

// BuildConfig creates a Config with default values and applies the given options
func BuildConfig(options ...ConfigOption) Config {
	config := Config{
		Extensions:     []string{"go"},
		Directories:    []string{"./src"},
		SkipDirs:       []string{},
		SpecificFiles:  []string{},
		IncludeIgnored: false,
	}

	for _, option := range options {
		option(&config)
	}

	return config
}

// WithExtensions sets the Extensions field of the Config
func WithExtensions(extensions ...string) ConfigOption {
	return func(c *Config) {
		c.Extensions = extensions
	}
}

// WithDirectories sets the Directories field of the Config
func WithDirectories(directories ...string) ConfigOption {
	return func(c *Config) {
		c.Directories = directories
	}
}

// WithSkipDirs sets the SkipDirs field of the Config
func WithSkipDirs(skipDirs ...string) ConfigOption {
	return func(c *Config) {
		c.SkipDirs = skipDirs
	}
}

// WithSpecificFiles sets the SpecificFiles field of the Config
func WithSpecificFiles(specificFiles ...string) ConfigOption {
	return func(c *Config) {
		c.SpecificFiles = specificFiles
	}
}

// WithIncludeIgnored sets the IncludeIgnored field of the Config
func WithIncludeIgnored(includeIgnored bool) ConfigOption {
	return func(c *Config) {
		c.IncludeIgnored = includeIgnored
	}
}
