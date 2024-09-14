package internal

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"gopkg.in/yaml.v3"
	"os"
	"path"
    "reflect"
	"strings"
)

type Config struct {
	Schemas   []string             `yaml:"schema"`
	Documents []string             `yaml:"documents"`
	Overwrite bool                 `yaml:"overwrite"`
	Generates map[string]Generates `yaml:"generates"`
}

type Generates struct {
	Plugins []string `yaml:"plugins"`
	Preset  string   `yaml:"preset"`
}

func (p *Project) GetConfig() Config {
	if !reflect.ValueOf(p.config.Schemas).IsZero() {
		return p.config
	}

	configFilePath := path.Join(p.RootDir, p.ConfigFile)
	dat, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}

	if strings.HasSuffix(p.ConfigFile, ".ts") || strings.HasSuffix(p.ConfigFile, ".js") {
		return ParseTSConfig(string(dat), configFilePath)
	}

	if strings.HasSuffix(p.ConfigFile, ".yml") || strings.HasSuffix(p.ConfigFile, ".yaml") {
		return ParseYAMLConfig(dat)
	}

	return Config{}
}

func ParseYAMLConfig(configData []byte) Config {
	parsedConfig := Config{}

	parseErr := yaml.Unmarshal(configData, &parsedConfig)
	if parseErr != nil {
		panic(parseErr)
	}

	return parsedConfig
}

// Parse dynamic configs (JS and TS)

/*
ParseJSConfig parses a given JS string
*/
func ParseJSConfig(configString string, filePath string) Config {
	bundledConfig, bundleErr := bundleJSConfigFile(filePath)
	if bundleErr != nil {
		panic(bundleErr)
	}

	config, executeErr := executeJSConfigFile(bundledConfig)
	if executeErr != nil {
		panic(executeErr)
	}

	return config
}

/*
ParseTSConfig parses a given TS string
*/
func ParseTSConfig(configString string, filePath string) Config {
	// because types are ignored, TS files can be handled by the JS parser
	return ParseJSConfig(configString, filePath)
}

/*
bundleJSConfigFile bundles a JS file to CJS format so it can be executed
*/
func bundleJSConfigFile(filePath string) (string, error) {
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{filePath},
		Bundle:      true,
		Write:       false,
		LogLevel:    api.LogLevelInfo,
		Format:      api.FormatCommonJS,
		Target:      api.ES2015,
	})

	if len(result.Errors) > 0 {
		return "", errors.New("could not bundle config file")
	}

	return string(result.OutputFiles[0].Contents), nil
}

func executeJSConfigFile(input string) (Config, error) {
	vm := goja.New()

	// Initialize module.exports
	module := vm.NewObject()
	module.Set("exports", vm.NewObject())
	vm.Set("module", module)

	// Run the JavaScript code
	_, err := vm.RunString(input)
	if err != nil {
		return Config{}, err
	}

	// Access module.exports.default
	exports := module.Get("exports").ToObject(vm)
	defaultExport := exports.Get("default")

	// Convert the JavaScript value to a Go value
	var exportResult map[string]interface{}
	err = vm.ExportTo(defaultExport, &exportResult)
	if err != nil {
		return Config{}, err
	}

	config := Config{
		Generates: make(map[string]Generates),
	}

	// Get 'schema' field
	if schemaValue, ok := exportResult["schema"]; ok {
		schemas, err := getStringOrStringSlice(schemaValue)
		if err != nil {
			return Config{}, fmt.Errorf("error parsing 'schema': %v", err)
		}
		config.Schemas = schemas
	} else {
		return Config{}, fmt.Errorf("'schema' field is required")
	}

	// Get 'overwrite' field
	if overwriteValue, ok := exportResult["overwrite"]; ok {
		overwrite, err := getBool(overwriteValue)
		if err != nil {
			return Config{}, fmt.Errorf("error parsing 'overwrite': %v", err)
		}
		config.Overwrite = overwrite
	}

	// Get 'generates' field
	if generatesValue, ok := exportResult["generates"]; ok {
		generatesMap, err := getMapStringInterface(generatesValue)
		if err != nil {
			return Config{}, fmt.Errorf("error parsing 'generates': %v", err)
		}

		for destination, destConfigValue := range generatesMap {
			destConfigMap, err := getMapStringInterface(destConfigValue)
			if err != nil {
				return Config{}, fmt.Errorf("error parsing 'generates[%s]': %v", destination, err)
			}

			generate := Generates{}

			if pluginsValue, ok := destConfigMap["plugins"]; ok {
				pluginsSlice, err := getStringSlice(pluginsValue)
				if err != nil {
					return Config{}, fmt.Errorf("error parsing 'plugins' in 'generates[%s]': %v", destination, err)
				}
				generate.Plugins = pluginsSlice
			}

			config.Generates[destination] = generate
		}
	}

	return config, nil
}

// Helper function to get a string or slice of strings
func getStringOrStringSlice(value interface{}) ([]string, error) {
	switch v := value.(type) {
	case string:
		return []string{v}, nil
	case []interface{}:
		return convertInterfaceSliceToStringSlice(v)
	default:
		return nil, fmt.Errorf("value is not a string or []string")
	}
}

// Helper function to get a boolean value
func getBool(value interface{}) (bool, error) {
	if b, ok := value.(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("value is not a boolean")
}

// Helper function to get map[string]interface{}
func getMapStringInterface(value interface{}) (map[string]interface{}, error) {
	if m, ok := value.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, fmt.Errorf("value is not an object")
}

// Helper function to convert interface{} to []string
func getStringSlice(value interface{}) ([]string, error) {
	if slice, ok := value.([]interface{}); ok {
		return convertInterfaceSliceToStringSlice(slice)
	}
	return nil, fmt.Errorf("value is not an array")
}

// Helper function to convert []interface{} to []string
func convertInterfaceSliceToStringSlice(slice []interface{}) ([]string, error) {
	result := make([]string, len(slice))
	for i, elem := range slice {
		if str, ok := elem.(string); ok {
			result[i] = str
		} else {
			return nil, fmt.Errorf("element at index %d is not a string", i)
		}
	}
	return result, nil
}
