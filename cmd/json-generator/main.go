package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"io/ioutil"
	"strings"

	astutils "github.com/melwyn95/code-generator/cmd/ast-utils"
)

const (
	importStmnt = `package common
	import (
		"bytes"
		"strconv"
	)

	//DO NOT EDIT. This code is auto-generated using go:generate json-generator

	`
	helpers = `func truncateLastComma(buff *bytes.Buffer) {
		if buff != nil {
			bufflen := buff.Len()
			if bufflen > 0 && string(buff.Bytes()[bufflen-1:]) == "," {
				buff.Truncate(bufflen - 1)
			}
		}
	}
	`
	funcStart   = "func (p *%s) MarshalJSON() ([]byte, error) {\n var json bytes.Buffer\n"
	returnStmnt = "return json.Bytes(), nil"
	endBlock    = "}\n"
	ignoreTag   = "@@"
	INT_TYPE    = "int"
	STRING_TYPE = "string"
)

var buff bytes.Buffer

func getValueFromTag(tag string) (string, error) {
	t := strings.Trim(tag, "\\`")
	values := strings.Split(t, ":")
	if len(values) > 0 {
		return values[1], nil
	}
	return "", errors.New("error: tag value not present")
}

func processPrimitiveValue(tagValue, fieldName, primitiveType string) {
	if tagValue != ignoreTag {
		buff.WriteString(fmt.Sprintf("json.WriteString(`%s:`)\n", tagValue))
	}

	switch primitiveType {
	case INT_TYPE:
		buff.WriteString(fmt.Sprintf("json.WriteString(strconv.Itoa(p.%s))\n", fieldName))
	case STRING_TYPE:
		buff.WriteString(fmt.Sprintf("json.WriteString(`\"` + p.%s + `\"`)\n", fieldName))
	}

	buff.WriteString(`json.WriteString(",")
	`)
}

func processStruct(st *ast.StructType, structName string) {
	buff.WriteString(fmt.Sprintf(funcStart, structName))
	buff.WriteString(`json.WriteString("{")
	`)
	for i := range st.Fields.List {
		field := st.Fields.List[i]
		tagValue, _ := getValueFromTag(field.Tag.Value)
		switch field.Type.(type) {
		case *ast.ArrayType:
			buff.WriteString(fmt.Sprintf("json.WriteString(`%s:`)\n", tagValue))
			buff.WriteString("json.WriteString(`[`)\n")
			buff.WriteString(fmt.Sprintf("for i := range p.%s {\n", field.Names[0]))
			primitiveType := field.Type.(*ast.ArrayType).Elt.(*ast.Ident).Name
			processPrimitiveValue(ignoreTag, field.Names[0].Name+"[i]", primitiveType)
			buff.WriteString("}\n")
			buff.WriteString("truncateLastComma(&json)\n")
			buff.WriteString("json.WriteString(`],`)\n")
		case *ast.MapType:
			buff.WriteString(fmt.Sprintf("for k, v := range p.%s {\n", field.Names[0]))
			buff.WriteString("json.WriteString(`\"` + k + `\"` + `:` + `\"` + v + `\",`)")
			buff.WriteString("}\n")
			buff.WriteString("truncateLastComma(&json)\n")
		case *ast.Ident:
			primitiveType := field.Type.(*ast.Ident).Name
			processPrimitiveValue(tagValue, field.Names[0].Name, primitiveType)
		}
	}
	buff.WriteString(`json.WriteString("}")
	`)
	buff.WriteString("truncateLastComma(&json)\n")
	buff.WriteString(returnStmnt)
	buff.WriteString(endBlock)
}

func processStructs(structs []*ast.StructType, structNames []string) {
	for i := range structs {
		processStruct(structs[i], structNames[i])
	}
}

func main() {
	packages, _ := astutils.GetPackages() // Handle err latr
	files := astutils.GetFiles(packages)
	declarations := astutils.GetDeclarations(files)
	typeDeclarations := astutils.GetTypeDecalations(declarations)
	structDeclarations, structNames := astutils.GetStructDeclarations(typeDeclarations)

	buff.WriteString(importStmnt)
	buff.WriteString(helpers)
	processStructs(structDeclarations, structNames)

	formattedBytes, err := format.Source(buff.Bytes())
	if err != nil {
		errors.New(err.Error())
	}

	ioutil.WriteFile("./generated_json.go", formattedBytes, 0644)
}

// Usefull link: https://stackoverflow.com/questions/20234342/get-a-simple-string-representation-of-a-struct-field-s-type
