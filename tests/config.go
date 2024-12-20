package tests

import . "github.com/fargusplumdoodle/dump_dir/src"

// ConfigOption is a function type that modifies a Config
type ConfigOption func(*Config)

// BuildConfig creates a Config with default values and applies the given options
func BuildConfig(opts ...ConfigOption) *Config {
	config := &Config{
		Action:         "dump_dir",
		Extensions:     []string{},
		Directories:    []string{},
		SkipDirs:       []string{},
		SpecificFiles:  []string{},
		IncludeIgnored: false,
		MaxFileSize:    500 * 1024, // Default to 500KB
		GlobPatterns:   nil,
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}

func WithAction(action string) ConfigOption {
	return func(c *Config) {
		c.Action = action
	}
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

func WithMaxFileSize(maxFileSize int64) ConfigOption {
	return func(c *Config) {
		c.MaxFileSize = maxFileSize
	}
}

func WithGlobPatterns(patterns ...string) ConfigOption {
	return func(c *Config) {
		c.GlobPatterns = patterns
	}
}
