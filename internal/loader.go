package internal

import (
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)
import (
	"errors"
	"os"
)

func LoadSchema(inputs ...string) (*ast.Schema, error) {
	if len(inputs) == 0 {
		return &ast.Schema{}, errors.New("no inputs given to load")
	}

	inputToLoad := inputs[0]

	// load file
	dat, err := os.ReadFile(inputToLoad)
	if err != nil {
		return &ast.Schema{}, err
	}

	schema, schemaParseError := gqlparser.LoadSchema(&ast.Source{
		BuiltIn: false,
		Input:   string(dat),
		Name:    "test.graphql",
	})
	if schemaParseError != nil {
		return &ast.Schema{}, schemaParseError
	}

	return schema, nil
}
