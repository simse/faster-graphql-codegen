package internal

import (
    "gopkg.in/yaml.v3"
    "os"
    "reflect"
    "strings"
)

type Config struct {
	Schema []string `yaml:"schema"`
	Documents []string `yaml:"documents"`
	Overwrite bool `yaml:"overwrite"`
	Generates map[string]struct{
		Plugins []string `yaml:"plugins"`
		Preset string `yaml:"preset"`
	} `yaml:"generates"`
}

func (p *Project) GetConfig() Config {
	if !reflect.ValueOf(p.config.Schema).IsZero() {
		return p.config
	}

	dat, err := os.ReadFile(p.RootDir + "/" + p.ConfigFile)
	if err != nil {
		panic(err)
	}

	if strings.HasSuffix(p.ConfigFile, ".ts") {
		return ParseTSConfig(string(dat))
	}

	if strings.HasSuffix(p.ConfigFile, ".yml") || strings.HasSuffix(p.ConfigFile, ".yaml") {
		return ParseYAMLConfig(dat)
	}

	return Config{}
}

func ParseTSConfig(configString string) Config {
	/*vm := goja.New()

	exports := vm.NewObject()
	module := vm.NewObject()
	module.Set("exports", exports)

	_, err := vm.RunString(fmt.Sprintf("(function(exports, module) { %s \n})(module.exports, module);", configString))
	if err != nil {
		panic(err)
	}

	println(exports.Get("config"))*/

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