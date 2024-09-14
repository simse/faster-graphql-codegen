package internal

import (
	"errors"
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

	_, executeErr := executeJSConfigFile(bundledConfig)
	if executeErr != nil {
		panic(executeErr)
	}

	os.Exit(0)

	return Config{}
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
		Target: 	 api.ES2015,
	})

	if len(result.Errors) > 0 {
		return "", errors.New("could not bundle config file")
	}

	return string(result.OutputFiles[0].Contents), nil
}

func executeJSConfigFile(input string) (Config, error) {
	vm := goja.New()

	// Step 1: Initialize module.exports
	module := vm.NewObject()
	module.Set("exports", vm.NewObject())
	vm.Set("module", module)

	_, err := vm.RunString(input)
	if err != nil {
		panic(err)
	}

	// Step 3: Access module.exports.default
	exports := module.Get("exports").ToObject(vm)
	defaultExport := exports.Get("default")

	// Step 4: Convert the JavaScript value to a Go value
	var exportResult map[string]interface{}
	err = vm.ExportTo(defaultExport, &exportResult)
	if err != nil {
		panic(err)
	}

	config := Config{}

	// get schema
	switch v := reflect.ValueOf(exportResult["schema"]); v.Kind() {
	case reflect.String:
		config.Schemas = []string{exportResult["schema"].(string)}
	case reflect.Array:
		config.Schemas = exportResult["schema"].([]string)
	default:
		panic("schema unknown type")
	}

	return config, nil
}
