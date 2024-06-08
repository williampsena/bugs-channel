// This package represents the dbless structure
package settings

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// Represents an error when attempting to parse a YAML file.
var ErrParseConfigFile = errors.New("an error occurred when trying to parse the configuration file yaml to struct")

// Represents an error when YAML file does not exists.
var ErrConfigFileNotFound = errors.New("an error occurred when trying to read the configuration file yaml to struct")

// Represents the configuration file structure from a YAML file.
type ConfigFile struct {
	// Configuration version
	Version string `yaml:"version"`
	// The organization
	Org string `yaml:"org"`
	// The service list
	Services []ConfigFileService `yaml:"services"`
}

// Represents service of the configuration file
type ConfigFileService struct {
	// Service id
	Id string `yaml:"id"`
	// Name of Service
	Name string `yaml:"name"`
	// The platform
	Platform string `yaml:"platform"`
	// Authentication keys
	AuthKeys []ConfigFileServiceAuthKey `yaml:"auth_keys"`
	// Service settings
	Settings ConfigFileServiceSettings `yaml:"settings"`
}

// Represents service authentication key of the configuration file
type ConfigFileServiceAuthKey struct {
	// Authorization key
	Key string `yaml:"key"`
	// Indicate that the authorization key is disabled.
	Disabled bool `yaml:"disabled"`
	// Unix time represents when the key will expire.
	ExpiredAt int64 `yaml:"expired_at"`
}

// Represents service settings of the configuration file
type ConfigFileServiceSettings struct {
	// The request limit per seconds
	RateLimit int `yaml:"rate_limit"`
}

// Read a yaml configuration file into ConfigFile struct.
func BuildConfigFile(filepath string) (*ConfigFile, error) {
	rawYaml, err := os.ReadFile(filepath)

	if err != nil {
		return &ConfigFile{}, errors.Join(ErrConfigFileNotFound, err)
	}

	return fromYamlToConfigFile(rawYaml)
}

// Parse a yaml configuration bytes into ConfigFile struct.
func fromYamlToConfigFile(rawYaml []byte) (*ConfigFile, error) {
	var configFile ConfigFile

	err := yaml.Unmarshal(rawYaml, &configFile)

	if err != nil {
		return nil, errors.Join(ErrParseConfigFile, err)
	}

	return &configFile, nil
}
