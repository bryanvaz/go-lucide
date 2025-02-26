package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	filepathPkg "path/filepath"
	"strings"
)

func kebabToCamelCase(input string) string {
	words := strings.Split(input, "-")
	c := cases.Title(language.English)
	for i := range words {
		words[i] = c.String(words[i])
	}
	return strings.Join(words, "")
}

type LucideIconAlias string

type LucideIconSvg struct {
	LucideIconSvgPath string
	LucideSvgContent  string
	LucideAliases     []LucideIconAlias
}

func (i *LucideIconSvg) KebabName() string {
	return i.Basename()
}

func (i *LucideIconSvg) CamelCaseName() string {
	return kebabToCamelCase(i.KebabName())
}
func (i *LucideIconSvg) Basename() string {
	return strings.TrimSuffix(filepathPkg.Base(i.LucideIconSvgPath), ".svg")
}
func (i *LucideIconSvg) LucideClasses() string {
	return "lucide lucide-" + i.Basename()
}

func (a *LucideIconAlias) CamelCaseName() string {
	return kebabToCamelCase(string(*a))
}

func injestIcons(lucideRepoPath string) ([]*LucideIconSvg, error) {
	iconsPath := filepathPkg.Join(lucideRepoPath, "icons")
	files, err := os.ReadDir(iconsPath)
	if err != nil {
		return nil, err
	}

	var svgFiles []*LucideIconSvg
	for _, file := range files {
		if filepathPkg.Ext(file.Name()) == ".svg" {
			filepath := filepathPkg.Join(iconsPath, file.Name())
			content, err := os.ReadFile(filepath)
			if err != nil {
				return nil, err
			}
			jsonFilePath := filepathPkg.Join(iconsPath, strings.TrimSuffix(file.Name(), ".svg")+".json")
			aliases := []LucideIconAlias{}
			if _, err = os.Stat(jsonFilePath); err == nil {
				jsonFile, err := os.ReadFile(jsonFilePath)
				if err != nil {
					return nil, err
				}
				type LucideIconJsonAlias struct {
					Name string `json:"name"`
				}
				type LucideIconJson struct {
					Aliases []json.RawMessage `json:"aliases"`
				}
				var iconJson LucideIconJson
				err = json.Unmarshal(jsonFile, &iconJson)
				if err != nil {
					fmt.Printf("Error unmarshalling json file: %s\n", jsonFilePath)
					return nil, err
				}
				for _, alias := range iconJson.Aliases {
					var aliasStr string
					var aliasObj LucideIconJsonAlias
					if err := json.Unmarshal(alias, &aliasStr); err == nil {
						aliases = append(aliases, LucideIconAlias(aliasStr))
						continue
					}
					if err := json.Unmarshal(alias, &aliasObj); err == nil {
						aliases = append(aliases, LucideIconAlias(aliasObj.Name))
						continue
					}
					fmt.Errorf("Error unmarshalling alias '%s' in file %s", alias, jsonFilePath)
				}
			}
			svgFiles = append(svgFiles, &LucideIconSvg{
				LucideIconSvgPath: filepath,
				LucideSvgContent:  string(content),
				LucideAliases:     aliases,
			})

		}
	}

	return svgFiles, nil
}
