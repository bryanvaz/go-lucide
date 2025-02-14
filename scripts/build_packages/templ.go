package main

import (
	"bytes"
	"fmt"
	"go/format"
	"golang.org/x/exp/slices"
	"strings"
	"text/template"

	"github.com/a-h/templ/cmd/templ/imports"
	"github.com/a-h/templ/generator"
	parser "github.com/a-h/templ/parser/v2"
)

const templateTemplFunc = `
// Renders the Lucide icon {{ .KebabCaseName }}.
templ {{ .FuncName }}(attrs ...templ.Attributes) {
<svg
    { at(attrs)... }
    class={ cn("{{ .LucideClasses }}", attrs) }
>
    {{ .Content }}
    { children... }
</svg>
}
`
const templateTemplFile = `package icons

{{ .Funcs }}
`

type TemplFuncTemplateParams struct {
	RootAttributes string
	FuncName       string
	LucideClasses  string
	KebabCaseName  string
	Content        string
}

type TemplFileTemplateParams struct {
	Funcs string
}

var tmplTemplFuncGen *template.Template
var tmplTemplFileGen *template.Template

func generateTemplFunc(icon *LucideIconSvg) (string, error) {
	var err error
	if tmplTemplFuncGen == nil {
		tmplTemplFuncGen, err = template.New("templTemplate").Parse(templateTemplFunc)
		if err != nil {
			return "", err
		}
	}
	splitStr := strings.Split(string(icon.LucideSvgContent), "<svg")
	closingIdx := strings.Index(splitStr[len(splitStr)-1], ">")
	attributes := splitStr[len(splitStr)-1][:closingIdx]
	coreSvgContent := splitStr[len(splitStr)-1][closingIdx+1:]
	content := strings.ReplaceAll(coreSvgContent, "</svg>", "")
	data := TemplFuncTemplateParams{
		RootAttributes: attributes,
		LucideClasses:  icon.LucideClasses(),
		KebabCaseName:  icon.Basename(),
		FuncName:       icon.CamelCaseName(),
		Content:        content,
	}

	var outputBuffer bytes.Buffer
	if err := tmplTemplFuncGen.Execute(&outputBuffer, data); err != nil {
		return "", err
	}

	outputString := outputBuffer.String() // Convert buffer to string
	return outputString, nil
}

func generateTemplFile(funcs ...string) (string, error) {
	var err error
	if tmplTemplFileGen == nil {
		tmplTemplFileGen, err = template.New("templFile").Parse(templateTemplFile)
		if err != nil {
			return "", err
		}
	}
	templFileData := TemplFileTemplateParams{
		Funcs: strings.Join(funcs, "\n"),
	}

	var outputFileBuffer bytes.Buffer
	if err := tmplTemplFileGen.Execute(&outputFileBuffer, templFileData); err != nil {
		return "", err
	}

	outputTemplString := outputFileBuffer.String()
	formattedTemplString, err := formatTemplFile(outputTemplString)
	if err != nil {
		return "", err
	}
	return formattedTemplString, nil
}

func formatTemplFile(templFile string) (string, error) {
	t, err := parser.ParseString(templFile)
	if err != nil {
		return "", err
	}
	t.Filepath = "placeholder.templ"
	t, err = imports.Process(t)
	if err != nil {
		return "", err
	}
	w := new(bytes.Buffer)
	if err = t.Write(w); err != nil {
		return "", err
	}
	return w.String(), nil
}

func generateGoFromTempl(templFile string) (string, error) {
	t, err := parser.ParseString(templFile)
	if err != nil {
		return "", err
	}
	b := new(bytes.Buffer)
	_, err = generator.Generate(t, b)
	if err != nil {
		return "", err
	}
	formattedGoCode, err := format.Source(b.Bytes())
	if err != nil {
		return "", err
	}
	return string(formattedGoCode), nil
}

const rollupFileTemplate = `
package icons

import templFuncs "github.com/bryanvaz/go-templ-lucide-icons/icons"

var (
	{{ .Content }}
)
`

func createRollupFile(icons []*LucideIconSvg) (string, error) {
	tmplRollupFileGen, err := template.New("rollupTemplate").Parse(rollupFileTemplate)
	if err != nil {
		return "", err
	}

	type tmplParams struct{ Content string }
	rollupLines := map[string][]string{}
	funcNames := []string{}
	for _, icon := range icons {
		lines := []string{}
		lines = append(lines, fmt.Sprintf("// Renders the Lucide icon '%s'.", icon.Basename()))
		lines = append(lines, fmt.Sprintf("%s = templFuncs.%s", icon.CamelCaseName(), icon.CamelCaseName()))
		rollupLines[icon.CamelCaseName()] = lines
		funcNames = append(funcNames, icon.CamelCaseName())
	}
	for _, icon := range icons {
		for _, alias := range icon.LucideAliases {
			_, alreadyExists := rollupLines[alias.CamelCaseName()]
			if alreadyExists {
				continue
			}
			aliasLines := []string{}
			aliasLines = append(aliasLines, fmt.Sprintf("// Alias for '%s'(%s).Renders the Lucide icon '%s'", icon.CamelCaseName(), icon.Basename(), alias))
			aliasLines = append(aliasLines, fmt.Sprintf("%s = templFuncs.%s", alias.CamelCaseName(), icon.CamelCaseName()))
			rollupLines[alias.CamelCaseName()] = aliasLines
			funcNames = append(funcNames, alias.CamelCaseName())
		}
	}
	slices.Sort(funcNames)
	content := ""
	for _, funcName := range funcNames {
		content += strings.Join(rollupLines[funcName], "\n") + "\n"
	}
	params := tmplParams{Content: content}
	var outputBuffer bytes.Buffer
	if err := tmplRollupFileGen.Execute(&outputBuffer, params); err != nil {
		return "", err
	}
	formattedOutput, err := format.Source(outputBuffer.Bytes())
	if err != nil {
		return "", err
	}

	return string(formattedOutput), nil
}
