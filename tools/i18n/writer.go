package main

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func writeKeysFiles(goDir, tsDir string, allKeys map[string]bool) {
	targetGoFile, err := os.Create(filepath.Join(goDir, "keys_generated.go"))
	if err != nil {
		panic(err)
	}
	defer targetGoFile.Close()

	targetTsFile, err := os.Create(filepath.Join(tsDir, "MessageKey.generated.ts"))
	if err != nil {
		panic(err)
	}
	defer targetTsFile.Close()

	data := keyFileData{
		Items: make([]itemData, 0, len(allKeys)),
	}

	caser := cases.Title(language.English)
	for key, _ := range allKeys {
		data.Items = append(data.Items, itemData{
			ConstName: toPascalCase(key, caser),
			Value:     key,
		})
	}

	goTmpl := template.Must(template.New("keys").Parse(goKeyFileTemplate))
	if err = goTmpl.Execute(targetGoFile, data); err != nil {
		panic(err)
	}

	tsTmpl := template.Must(template.New("keys").Parse(tsKeyFileTemplate))
	if err = tsTmpl.Execute(targetTsFile, data); err != nil {
		panic(err)
	}
}

func writeDictionaryFile(dir string, lang string, props map[string]string) {
	fileName := strings.ReplaceAll(strings.ToLower(lang), "-", "_") + "_generated.go"
	dictFile, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		panic(err)
	}

	defer dictFile.Close()

	langTag := strings.ReplaceAll(lang, "_", "-")
	parts := strings.Split(langTag, "-")

	var funcNameBuilder strings.Builder
	for _, part := range parts {
		funcNameBuilder.WriteString(strings.ToUpper(part))
	}

	keys := make([]itemData, 0, len(props))
	caser := cases.Title(language.English)

	for key, value := range props {
		keys = append(keys, itemData{
			ConstName: toPascalCase(key, caser),
			Value:     strings.ReplaceAll(value, "\"", "\\\""),
		})
	}

	data := dictFileData{
		FuncName: funcNameBuilder.String(),
		LangTag:  langTag,
		Keys:     keys,
	}

	tmpl := template.Must(template.New("dict").Parse(dictionaryFileTemplate))
	if err = tmpl.Execute(dictFile, data); err != nil {
		panic(err)
	}
}
