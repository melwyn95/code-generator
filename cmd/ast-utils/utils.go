package ast_utils

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func GetPackages() ([]*ast.Package, error) {
	var packages []*ast.Package
	set := token.NewFileSet()

	packs, err := parser.ParseDir(set, ".", nil, 0)
	for _, p := range packs {
		packages = append(packages, p)
	}

	return packages, err
}

func GetFiles(packages []*ast.Package) []*ast.File {
	var files []*ast.File
	for i := range packages {
		for _, file := range packages[i].Files {
			files = append(files, file)
		}
	}
	return files
}

func GetDeclarations(files []*ast.File) []*ast.GenDecl {
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

func GetTypeDecalations(genericDeclations []*ast.GenDecl) []*ast.TypeSpec {
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

func GetStructDeclarations(typeDeclarations []*ast.TypeSpec) ([]*ast.StructType, []string) {
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
