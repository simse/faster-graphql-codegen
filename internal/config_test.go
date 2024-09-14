package internal

import (
    "reflect"
    "testing"
)

// Assuming your main code is in a file named `config.go`
// We'll write the tests in `config_test.go`

// TestExecuteJSConfigFile tests the executeJSConfigFile function
func TestExecuteJSConfigFile(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected Config
        wantErr  bool
    }{
        {
            name: "ValidConfigWithStringSchema",
            input: `
            var config = {
                schema: "schema.graphql",
                overwrite: true,
                generates: {
                    "output.ts": {
                        plugins: ["typescript"]
                    }
                }
            };
            module.exports = { default: config };
            `,
            expected: Config{
                Schemas:   []string{"schema.graphql"},
                Overwrite: true,
                Generates: map[string]Generates{
                    "output.ts": {
                        Plugins: []string{"typescript"},
                    },
                },
            },
            wantErr: false,
        },
        {
            name: "ValidConfigWithArraySchema",
            input: `
            var config = {
                schema: ["schema1.graphql", "schema2.graphql"],
                overwrite: false,
                generates: {
                    "output.ts": {
                        plugins: ["typescript", "graphql-codegen"]
                    }
                }
            };
            module.exports = { default: config };
            `,
            expected: Config{
                Schemas:   []string{"schema1.graphql", "schema2.graphql"},
                Overwrite: false,
                Generates: map[string]Generates{
                    "output.ts": {
                        Plugins: []string{"typescript", "graphql-codegen"},
                    },
                },
            },
            wantErr: false,
        },
        {
            name: "MissingSchemaField",
            input: `
            var config = {
                overwrite: true,
                generates: {
                    "output.ts": {
                        plugins: ["typescript"]
                    }
                }
            };
            module.exports = { default: config };
            `,
            expected: Config{},
            wantErr:  true,
        },
        {
            name: "InvalidSchemaType",
            input: `
            var config = {
                schema: 123,
                overwrite: true,
                generates: {
                    "output.ts": {
                        plugins: ["typescript"]
                    }
                }
            };
            module.exports = { default: config };
            `,
            expected: Config{},
            wantErr:  true,
        },
        {
            name: "InvalidOverwriteType",
            input: `
            var config = {
                schema: "schema.graphql",
                overwrite: "yes",
                generates: {
                    "output.ts": {
                        plugins: ["typescript"]
                    }
                }
            };
            module.exports = { default: config };
            `,
            expected: Config{},
            wantErr:  true,
        },
        {
            name: "InvalidGeneratesType",
            input: `
            var config = {
                schema: "schema.graphql",
                overwrite: true,
                generates: "invalid"
            };
            module.exports = { default: config };
            `,
            expected: Config{},
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := executeJSConfigFile(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("executeJSConfigFile() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("executeJSConfigFile() = %+v, expected %+v", result, tt.expected)
            }
        })
    }
}

// TestGetStringOrStringSlice tests the getStringOrStringSlice function
func TestGetStringOrStringSlice(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected []string
        wantErr  bool
    }{
        {
            name:     "StringInput",
            input:    "single string",
            expected: []string{"single string"},
            wantErr:  false,
        },
        {
            name:     "StringSliceInput",
            input:    []interface{}{"string1", "string2"},
            expected: []string{"string1", "string2"},
            wantErr:  false,
        },
        {
            name:    "InvalidInputType",
            input:   123,
            wantErr: true,
        },
        {
            name:    "SliceWithNonString",
            input:   []interface{}{"string", 123},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := getStringOrStringSlice(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("getStringOrStringSlice() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("getStringOrStringSlice() = %v, expected %v", result, tt.expected)
            }
        })
    }
}

// TestGetBool tests the getBool function
func TestGetBool(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected bool
        wantErr  bool
    }{
        {
            name:     "BoolTrue",
            input:    true,
            expected: true,
            wantErr:  false,
        },
        {
            name:     "BoolFalse",
            input:    false,
            expected: false,
            wantErr:  false,
        },
        {
            name:    "InvalidType",
            input:   "true",
            wantErr: true,
        },
        {
            name:    "IntegerInput",
            input:   1,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := getBool(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("getBool() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && result != tt.expected {
                t.Errorf("getBool() = %v, expected %v", result, tt.expected)
            }
        })
    }
}

// TestGetMapStringInterface tests the getMapStringInterface function
func TestGetMapStringInterface(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected map[string]interface{}
        wantErr  bool
    }{
        {
            name: "ValidMap",
            input: map[string]interface{}{
                "key1": "value1",
                "key2": 2,
            },
            expected: map[string]interface{}{
                "key1": "value1",
                "key2": 2,
            },
            wantErr: false,
        },
        {
            name:    "InvalidType",
            input:   []interface{}{"key1", "value1"},
            wantErr: true,
        },
        {
            name:    "NilInput",
            input:   nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := getMapStringInterface(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("getMapStringInterface() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("getMapStringInterface() = %v, expected %v", result, tt.expected)
            }
        })
    }
}

// TestGetStringSlice tests the getStringSlice function
func TestGetStringSlice(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected []string
        wantErr  bool
    }{
        {
            name:     "ValidStringSlice",
            input:    []interface{}{"string1", "string2"},
            expected: []string{"string1", "string2"},
            wantErr:  false,
        },
        {
            name:    "InvalidType",
            input:   "string",
            wantErr: true,
        },
        {
            name:    "SliceWithNonString",
            input:   []interface{}{"string1", 2},
            wantErr: true,
        },
        {
            name:    "EmptySlice",
            input:   []interface{}{},
            expected: []string{},
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := getStringSlice(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("getStringSlice() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("getStringSlice() = %v, expected %v", result, tt.expected)
            }
        })
    }
}

// TestConvertInterfaceSliceToStringSlice tests the convertInterfaceSliceToStringSlice function
func TestConvertInterfaceSliceToStringSlice(t *testing.T) {
    tests := []struct {
        name     string
        input    []interface{}
        expected []string
        wantErr  bool
    }{
        {
            name:     "ValidConversion",
            input:    []interface{}{"string1", "string2"},
            expected: []string{"string1", "string2"},
            wantErr:  false,
        },
        {
            name:    "NonStringElement",
            input:   []interface{}{"string1", 2},
            wantErr: true,
        },
        {
            name:     "EmptySlice",
            input:    []interface{}{},
            expected: []string{},
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := convertInterfaceSliceToStringSlice(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("convertInterfaceSliceToStringSlice() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("convertInterfaceSliceToStringSlice() = %v, expected %v", result, tt.expected)
            }
        })
    }
}
