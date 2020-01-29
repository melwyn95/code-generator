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
		"fmt"
	)
	`
	funcStart = "func (p *%s) MarshalJSON() {\n var json bytes.Buffer\n"
	endBlock  = "}\n"
	ignoreTag = "@@"
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

func processPrimitiveValue(tagValue, fieldName string) {
	if tagValue != ignoreTag {
		buff.WriteString(fmt.Sprintf("json.WriteString(`%s:`)\n", tagValue))
	}
	buff.WriteString(fmt.Sprintf("json.WriteString(fmt.Sprintf(\"\", p.%s) + `,`)\n", fieldName))
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
			buff.WriteString("json.WriteString(`[`)\n")
			buff.WriteString(fmt.Sprintf("for i := range p.%s {\n", field.Names[0]))
			processPrimitiveValue(ignoreTag, field.Names[0].Name+"[i]")
			buff.WriteString("}\n")
			buff.WriteString("json.WriteString(`],`)\n")
		case *ast.MapType:
			buff.WriteString(fmt.Sprintf("for k, v := range p.%s {\n", field.Names[0]))
			buff.WriteString("json.WriteString(`\"` + k + `\"` + `:` + `\"` + v + `\",`)")
			buff.WriteString("}\n")
		default:
			processPrimitiveValue(tagValue, field.Names[0].Name)
		}
	}
	buff.WriteString(`json.WriteString("}")
	`)
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
	processStructs(structDeclarations, structNames)

	fmt.Println(buff.String())

	formattedBytes, err := format.Source(buff.Bytes()) // Handle err latr
	fmt.Println(string(formattedBytes), err)

	ioutil.WriteFile("./generated_json.go", formattedBytes, 0644)
}
