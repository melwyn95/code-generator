package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"

	astutils "github.com/melwyn95/code-generator/cmd/ast-utils"
)

func main() {
	packages, _ := astutils.GetPackages() // Handle err latr
	files := astutils.GetFiles(packages)
	declarations := astutils.GetDeclarations(files)
	typeDeclarations := astutils.GetTypeDecalations(declarations)
	_, structNames := astutils.GetStructDeclarations(typeDeclarations)

	templateString, err := ioutil.ReadFile("template.txt")
	if err != nil {
		errors.New(err.Error())
	}

	codeTemplate := template.Must(template.New("code").Parse(string(templateString)))

	for _, structName := range structNames {
		var output bytes.Buffer
		err := codeTemplate.Execute(&output, struct {
			Type string
		}{
			Type: structName,
		})
		if err != nil {
			errors.New(err.Error())
		}

		fileName := fmt.Sprintf("./%sMapper.go", structName)
		ioutil.WriteFile(fileName, output.Bytes(), 0644)
	}
}
