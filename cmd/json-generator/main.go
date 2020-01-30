package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
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

func getPackages() ([]*ast.Package, error) {
	var packages []*ast.Package
	set := token.NewFileSet()

	packs, err := parser.ParseDir(set, ".", nil, 0)
	for _, p := range packs {
		packages = append(packages, p)
	}

	return packages, err
}

func getFiles(packages []*ast.Package) []*ast.File {
	var files []*ast.File
	for i := range packages {
		for _, file := range packages[i].Files {
			files = append(files, file)
		}
	}
	return files
}

func getDeclarations(files []*ast.File) []*ast.GenDecl {
	var genericDeclarations []*ast.GenDecl
	for _, file := range files {
		for _, declaration := range file.Decls {
			if genericDeclaration, ok := declaration.(*ast.GenDecl); ok {
				genericDeclarations = append(genericDeclarations, genericDeclaration)
			}
		}
	}
	return genericDeclarations
}

func getTypeDecalations(genericDeclations []*ast.GenDecl) []*ast.TypeSpec {
	var typeDeclarations []*ast.TypeSpec
	for _, genericDeclation := range genericDeclations {
		for _, specs := range genericDeclation.Specs {
			if typespec, ok := specs.(*ast.TypeSpec); ok {
				typeDeclarations = append(typeDeclarations, typespec)
			}
		}
	}
	return typeDeclarations
}

func getStructDeclarations(typeDeclarations []*ast.TypeSpec) ([]*ast.StructType, []string) {
	var structDeclarations []*ast.StructType
	var structNames []string
	for _, typeDeclaration := range typeDeclarations {
		if structDeclaration, ok := typeDeclaration.Type.(*ast.StructType); ok {
			structDeclarations = append(structDeclarations, structDeclaration)
			structNames = append(structNames, typeDeclaration.Name.Name)
		}
	}
	return structDeclarations, structNames
}

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
	packages, _ := getPackages() // Handle err latr
	files := getFiles(packages)
	declarations := getDeclarations(files)
	typeDeclarations := getTypeDecalations(declarations)
	structDeclarations, structNames := getStructDeclarations(typeDeclarations)

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
