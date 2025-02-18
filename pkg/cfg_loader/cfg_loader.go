package cfgloader

import (
	"log"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)


type Option struct {
	Prefix       string
	Delimiter    string
	Separator    string
	YamlFilePath string
	CallbackEnv  func(string) string
}

// defaultCallbackEnv processes environment variable keys based on provided prefix and separator
func defaultCallbackEnv(source, prefix, separator string) string {
	base := strings.ToLower(strings.TrimPrefix(source, prefix))
	return strings.ReplaceAll(base, separator, ".")
}

// Load function loads configuration from YAML file and environment variables based on provided options
func Load(options Option, config interface{}) error {
	k := koanf.New(options.Delimiter)

	// Load configuration from YAML file if provided
	if options.YamlFilePath != "" {
			if err := k.Load(file.Provider(options.YamlFilePath), yaml.Parser()); err != nil {
					log.Fatalf("Error loading config file: %v", err)
					return err
			}
	}

	// Define callback function for environment variables
	callback := options.CallbackEnv
	if callback == nil {
			// Set default callback using the prefix and separator from options
			callback = func(source string) string {
					return defaultCallbackEnv(source, options.Prefix, options.Separator)
			}
	}

	// Load environment variables with the specified prefix and callback
	if err := k.Load(env.Provider(options.Prefix, options.Delimiter, callback), nil); err != nil {
			log.Fatalf("Error loading environment variables: %v", err)
			return err
	}

	// Unmarshal into provided config structure (passing address)
	if err := k.Unmarshal("", &config); err != nil {
			log.Fatalf("Error unmarshalling config: %v", err)
			return err
	}

	return nil
}