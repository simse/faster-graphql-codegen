package plugins

import (
	"errors"
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

type PluginTask struct {
	Schema *ast.Schema
	Output *strings.Builder
	Config interface{}
}

/*
VerifyPlugin checks if a plugin is executable by faster-graphql-codegen
*/
func VerifyPlugin(pluginName string, config interface{}) error {
	if pluginName != "typescript" && pluginName != "introspection" {
		return errors.New("unknown plugin")
	}

	return nil
}
