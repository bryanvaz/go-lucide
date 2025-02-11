package main

import (
	"bytes"
	"fmt"
	"github.com/a-h/templ/cmd/templ/imports"
	"github.com/a-h/templ/generator"
	parser "github.com/a-h/templ/parser/v2"
	"go/format"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const templateContent = `package icons

templ {{ .FuncName }}(attrs ...templ.Attributes) {
<svg
    class={ cn("{{ .LucideClasses }}", attrs) }
    { at(attrs)... }
{{ .Content }}
}
`

const templateTemplFunc = `
templ {{ .FuncName }}(attrs ...templ.Attributes) {
<svg
    class={ cn("{{ .LucideClasses }}", attrs) }
    { at(attrs)... }
{{ .Content }}
}
`
const templateTemplFile = `package icons

{{ .Funcs }}
`

type SVGFile struct {
	LucideClasses string
	FuncName      string
	Content       string
}

type TemplFile struct {
	Funcs string
}

func kebabToCamelCase(input string) string {
	words := strings.Split(input, "-")
	c := cases.Title(language.English)
	for i := range words {
		words[i] = c.String(words[i])
	}
	return strings.Join(words, "")
}

func injectLine(content string, injection string, afterLine int) string {
	lines := strings.Split(content, "\n")
	if len(lines) > afterLine {
		lines = append(lines[:afterLine], append([]string{injection}, lines[afterLine:]...)...)
	} else {
		lines = append(lines, injection)
	}
	return strings.Join(lines, "\n")
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func main() {
	inputDir := "./lucide/icons"
	outputDir := "./templ"

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	// Get all SVG files
	files, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var svgFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".svg" {
			svgFiles = append(svgFiles, file.Name())
		}
	}

	// Parse template
	// tmpl, err := template.New("svgTemplate").Parse(templateContent)
	// if err != nil {
	// 	fmt.Println("Error parsing template:", err)
	// 	return
	// }
	tmplTemplFunc, err := template.New("templTemplate").Parse(templateTemplFunc)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}
	tmplTemplFile, err := template.New("templFile").Parse(templateTemplFile)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// Remove existing .templ files
	existingFiles, err := os.ReadDir(outputDir)
	if err != nil {
		fmt.Println("Error reading output directory:", err)
		return
	}
	for _, file := range existingFiles {
		if filepath.Ext(file.Name()) == ".templ" ||
			filepath.Ext(file.Name()) == ".go" {
			os.Remove(filepath.Join(outputDir, file.Name()))
		}
	}

	funcs := make([]string, 0)

	fmt.Println("Processing", len(svgFiles), "SVG files")
	// Process each SVG file
	for _, filename := range svgFiles {
		inputPath := filepath.Join(inputDir, filename)
		baseFilename := strings.TrimSuffix(filename, ".svg")
		// outputFilename := baseFilename + ".templ"
		camelCaseName := kebabToCamelCase(baseFilename)
		// outputPath := filepath.Join(outputDir, outputFilename)

		content, err := os.ReadFile(inputPath)
		if err != nil {
			fmt.Println("Error reading file:", filename, err)
			continue
		}

		splitStr := strings.Split(string(content), "<svg")

		data := SVGFile{
			LucideClasses: "lucide lucide-" + baseFilename,
			FuncName:      camelCaseName,
			Content:       splitStr[len(splitStr)-1],
		}

		// Write template output to a string first
		var outputBuffer bytes.Buffer
		// if err := tmpl.Execute(&outputBuffer, data); err != nil {
		// 	fmt.Println("Error executing template for:", filename, err)
		// 	continue
		// }
		if err := tmplTemplFunc.Execute(&outputBuffer, data); err != nil {
			fmt.Println("Error executing template for:", filename, err)
			continue
		}

		outputString := outputBuffer.String() // Convert buffer to string

		funcs = append(funcs, outputString)

	}

	fmt.Println("Generating final templ file")

	outputTemplPath := filepath.Join(outputDir, "icons.templ")
	outputGoTemplPath := filepath.Join(outputDir, "icons_templ.go")

	templFileData := TemplFile{
		Funcs: strings.Join(funcs, "\n"),
	}

	var outputFileBuffer bytes.Buffer
	if err := tmplTemplFile.Execute(&outputFileBuffer, templFileData); err != nil {
		fmt.Println("Error executing template for:", "templFileData", err)
		os.Exit(1)
	}

	outputTemplString := outputFileBuffer.String()

	fmt.Println("Formatting final templ file")
	t, err := parser.ParseString(outputFileBuffer.String())
	if err != nil {
		fmt.Println("Error parsing template by templ:", err)
		os.Exit(1)
	}
	t.Filepath = outputTemplPath
	t, err = imports.Process(t)
	if err != nil {
		fmt.Println("Error processing templ imports:", err)
		os.Exit(1)
	}
	w := new(bytes.Buffer)
	if err = t.Write(w); err != nil {
		fmt.Println("formatting error with templ:", err)
	} else {
		outputTemplString = w.String()
	}

	// Write the string output to file
	if err := os.WriteFile(outputTemplPath, []byte(outputTemplString), 0644); err != nil {
		fmt.Println("Error writing to output file:", outputTemplPath, err)
	}
	fmt.Println("Templ file saved to", outputTemplPath)

	fmt.Println("Generating go file")

	// Generate the go file from templ
	t, err = parser.ParseString(outputTemplString)
	if err != nil {
		fmt.Println("Error parsing template by templ:", err)
		os.Exit(1)
	}
	b := new(bytes.Buffer)
	_, err = generator.Generate(t, b)
	if err != nil {
		fmt.Println("Error generating template:", err)
		os.Exit(1)
	}

	fmt.Println("Formatting go file")
	formattedGoCode, err := format.Source(b.Bytes())
	if err != nil {
		fmt.Println("Error formatting generated go code:", err)
		os.Exit(1)
	}

	// Write the string output to file
	if err := os.WriteFile(outputGoTemplPath, formattedGoCode, 0644); err != nil {
		fmt.Println("Error writing to output file:", outputGoTemplPath, err)
	}
	fmt.Println("Go file saved to", outputGoTemplPath)

	// Copy utils.go file
	srcUtilsPath := filepath.Join(outputDir, "utils/utils.go")
	dstUtilsPath := filepath.Join(outputDir, "utils.go")
	if err := copyFile(srcUtilsPath, dstUtilsPath); err != nil {
		fmt.Println("Error copying utils.go:", err)
	}
	fmt.Println("Copied utils.go")

	fmt.Println("Processing complete. Files saved in", outputDir)
}
