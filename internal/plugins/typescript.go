package plugins

import (
    "errors"
    "github.com/vektah/gqlparser/v2/ast"
    "log/slog"
    "strings"
    "time"
)

func (p* PluginTask) Typescript() {
    ConvertSchema(p.Schema, p.Output)
}

/*
ConvertSchema converts a graphql schema to Typescript output
*/
func ConvertSchema(schema *ast.Schema, output *strings.Builder) {
    output.WriteString("/* Generated by faster-graphql-codegen on " + time.Now().Format(time.DateTime) + " */\n")

    AddBaseTypes(output)
    knownScalars := AddScalars(schema, output)

    for _, definition := range schema.Types {
        if definition.BuiltIn {
            continue
        }

        err := ConvertDefinition(definition, output, knownScalars)
        if err != nil {
            slog.Error(err.Error(), "kind", definition.Kind, "name", definition.Name)
        } else {
            output.WriteString("\n")
        }
    }
}

var builtInScalars = map[string]string{
    "String": "string",
}

/*
AddScalars parses a schema and outputs a Scalars type, it also returns a list of scalars it found
*/
func AddScalars(schema *ast.Schema, output *strings.Builder) []*ast.Definition {
    output.WriteString("/** All built-in and custom scalars, mapped to their actual values */\n")
    output.WriteString("export type Scalars = {\n")

    var scalars []*ast.Definition

    for _, definition := range schema.Types {
        if definition.Kind == ast.Scalar {
            scalars = append(scalars, definition)

            output.WriteString("\t" + definition.Name + ": ")

            if knownScalarType, ok := builtInScalars[definition.Name]; ok {
                output.WriteString(knownScalarType)
            } else {
                output.WriteString("any")
            }

            output.WriteString(";\n")
        }
    }

    output.WriteString("};\n")

    return scalars
}

func AddBaseTypes(output *strings.Builder) {
    output.WriteString("export type Maybe<T> = T | null;\nexport type InputMaybe<T> = Maybe<T>;\nexport type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };\nexport type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };\nexport type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };\nexport type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };\nexport type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };")
    output.WriteString("\n")
}

func ConvertDefinition(definition *ast.Definition, output *strings.Builder, knownScalars []*ast.Definition) error {
    switch definition.Kind {
    case ast.Enum:
        ConvertEnum(definition, output)
    case ast.Union:
        ConvertUnion(definition, output)
    case ast.Interface:
        ConvertInterface(definition, output, knownScalars)
    case ast.Object:
        ConvertObject(definition, output, knownScalars)
    case ast.InputObject:
        ConvertInputObject(definition, output, knownScalars)
    case ast.Scalar:
    default:
        return errors.New("unknown definition kind")
    }

    return nil
}

func ConvertEnum(definition *ast.Definition, output *strings.Builder) {
    enumName := ToCamel(definition.Name)

    WriteComment(definition, output)

    output.WriteString("export enum " + enumName + " {\n")

    for i, enumValue := range definition.EnumValues {
        enumName := enumValue.Name
        enumKey := ToUpper(enumValue.Name)
        output.WriteString("\t" + enumKey + " = '" + enumName + "'")

        if i != len(definition.EnumValues)-1 {
            output.WriteString(",")
        }
        output.WriteString("\n")
    }
    output.WriteString("}\n")
}

func WriteComment(definition *ast.Definition, output *strings.Builder) {
    comment := definition.Description

    if comment != "" {
        output.WriteString("/* " + comment + " */\n")
    }
}

func WriteFieldComment(definition *ast.FieldDefinition, output *strings.Builder) {
    comment := definition.Description

    if comment != "" {
        output.WriteString("\t/* " + comment + " */\n")
    }
}

func WriteArgumentComment(definition *ast.ArgumentDefinition, output *strings.Builder) {
    comment := definition.Description

    if comment != "" {
        output.WriteString("\t/* " + comment + " */\n")
    }
}

func ConvertUnion(definition *ast.Definition, output *strings.Builder) {
    unionName := ToCamel(definition.Name)

    WriteComment(definition, output)

    output.WriteString("export type " + unionName + " = ")
    for i, alias := range definition.Types {
        output.WriteString(alias)

        if i != len(definition.Types)-1 {
            output.WriteString(" | ")
        } else {
            output.WriteString(";")
        }
    }
    output.WriteString("\n")
}

func ConvertInterface(definition *ast.Definition, output *strings.Builder, knownScalars []*ast.Definition) {
    interfaceName := ToCamel(definition.Name)

    WriteComment(definition, output)

    output.WriteString("export type " + interfaceName + " = {\n")
    for _, field := range definition.Fields {
        WriteFieldComment(field, output)

        fieldName := field.Name

        output.WriteString("\t" + fieldName)
        AddFieldType(field, output, "Maybe", knownScalars)
    }

    output.WriteString("}\n")

    for _, field := range definition.Fields {
        WriteFieldArguments(field, output, knownScalars, interfaceName)
    }
}

func ConvertObject(definition *ast.Definition, output *strings.Builder, knownScalars []*ast.Definition) {
    interfaceName := ToCamel(definition.Name)

    WriteComment(definition, output)

    output.WriteString("export type " + interfaceName + " = ")

    // check implements
    for _, implementsInterface := range definition.Interfaces {
        output.WriteString(implementsInterface + " & ")
    }
    output.WriteString("{\n")

    output.WriteString("\t__typename: '" + interfaceName + "';\n")
    for _, field := range definition.Fields {
        if field.Name == "__type" || field.Name == "__schema" {
            continue
        }

        WriteFieldComment(field, output)

        fieldName := field.Name

        output.WriteString("\t" + fieldName)
        AddFieldType(field, output, "Maybe", knownScalars)
    }

    output.WriteString("}\n")

    for _, field := range definition.Fields {
        WriteFieldArguments(field, output, knownScalars, interfaceName)
    }
}

func WriteFieldArguments(definition *ast.FieldDefinition, output *strings.Builder, knownScalars []*ast.Definition, rootName string) {
    if len(definition.Arguments) == 0 {
        return
    }

    fieldName := ToCamel(definition.Name)

    WriteFieldComment(definition, output)
    output.WriteString("export type " + rootName + fieldName + "Args = {\n")
    for _, argument := range definition.Arguments {
        WriteArgumentComment(argument, output)
        output.WriteString("\t" + argument.Name)
        AddArgumentType(argument, output, "Maybe", knownScalars)
    }
    output.WriteString("}\n")
}

func ConvertInputObject(definition *ast.Definition, output *strings.Builder, knownScalars []*ast.Definition) {
    interfaceName := ToCamel(definition.Name)

    WriteComment(definition, output)

    output.WriteString("export type " + interfaceName + " = {\n")
    for _, field := range definition.Fields {
        WriteFieldComment(field, output)

        fieldName := field.Name

        output.WriteString("\t" + fieldName)
        AddFieldType(field, output, "InputMaybe", knownScalars)
    }

    output.WriteString("}\n")
}

func AddFieldType(definition *ast.FieldDefinition, output *strings.Builder, maybeType string, knownScalars []*ast.Definition) {
    isNullable := !definition.Type.NonNull

    if isNullable {
        output.WriteString("?: " + maybeType + "<")
    } else {
        output.WriteString(": ")
    }

    isArray := definition.Type.Elem != nil
    isElemNullable := true

    if isArray {
        output.WriteString("Array<")

        isElemNullable = !definition.Type.Elem.NonNull
    }

    if isArray && isElemNullable {
        output.WriteString(maybeType + "<")
    }

    // check if scalar is known
    isScalarKnown := false
    for _, scalar := range knownScalars {
        if scalar.Name == definition.Type.Name() {
            isScalarKnown = true
        }
    }

    if isScalarKnown {
        output.WriteString("Scalars['" + definition.Type.Name() + "']")
    } else {
        output.WriteString(ToCamel(definition.Type.Name()))
    }

    if isNullable {
        output.WriteString(">")
    }

    if isArray {
        output.WriteString(">")
    }

    if isArray && isElemNullable {
        output.WriteString(">")
    }

    output.WriteString(";\n")
}

func AddArgumentType(definition *ast.ArgumentDefinition, output *strings.Builder, maybeType string, knownScalars []*ast.Definition) {
    isNullable := !definition.Type.NonNull

    if isNullable {
        output.WriteString("?: " + maybeType + "<")
    } else {
        output.WriteString(": ")
    }

    isArray := definition.Type.Elem != nil
    isElemNullable := true

    if isArray {
        output.WriteString("Array<")

        isElemNullable = !definition.Type.Elem.NonNull
    }

    if isArray && isElemNullable {
        output.WriteString(maybeType + "<")
    }

    // check if scalar is known
    isScalarKnown := false
    for _, scalar := range knownScalars {
        if scalar.Name == definition.Type.Name() {
            isScalarKnown = true
        }
    }

    if isScalarKnown {
        output.WriteString("Scalars['" + definition.Type.Name() + "']")
    } else {
        output.WriteString(ToCamel(definition.Type.Name()))
    }

    if isNullable {
        output.WriteString(">")
    }

    if isArray {
        output.WriteString(">")
    }

    if isArray && isElemNullable {
        output.WriteString(">")
    }

    output.WriteString(";\n")
}