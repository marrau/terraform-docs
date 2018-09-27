package print

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"html/template"

	"github.com/Masterminds/sprig"
	"github.com/segmentio/terraform-docs/doc"
)

type templateStruct struct {
	Doc           doc.Doc
	PrintRequired bool
}

// Pretty printer pretty prints a doc.
func Pretty(d doc.Doc) (string, error) {
	var buf bytes.Buffer

	if len(d.Comment) > 0 {
		buf.WriteString(fmt.Sprintf("\n%s\n", d.Comment))
	}

	if len(d.Version) > 0 {
		format := "  \033[36mterraform.required_version\033[0m (%s)\n\n\n"
		buf.WriteString(fmt.Sprintf(format, d.Version))
	}

	if len(d.Providers) > 0 {
		buf.WriteString("\n")

		for _, i := range d.Providers {
			format := "  \033[36mprovider.%s\033[0m\n Version: %s\n  \033[90m%s\033[0m\n  \033[90m%s\033[0m\n\n"
			s := fmt.Sprintf(format, i.Name, i.Version, i.Documentation, strings.TrimSpace(i.Description))
			buf.WriteString(s)
		}

		buf.WriteString("\n")
	}

	if len(d.Modules) > 0 {
		buf.WriteString("\n")

		for _, i := range d.Modules {
			format := "  \033[36mmodule.%s\033[0m\n  \033[90m%s\033[0m\n\n"
			buf.WriteString(fmt.Sprintf(format, i.Name, i.Description))
		}

		buf.WriteString("\n")
	}

	if len(d.Resources) > 0 {
		buf.WriteString("\n")

		for _, i := range d.Resources {
			format := "  \033[36mresource.%s.%s\033[0m\n  \033[90m%s\033[0m\n\n"
			buf.WriteString(fmt.Sprintf(format, i.Type, i.Name, i.Documentation))
		}

		buf.WriteString("\n")
	}

	if len(d.Inputs) > 0 {
		buf.WriteString("\n")

		for _, i := range d.Inputs {
			format := "  \033[36mvar.%s\033[0m (%s)\n  \033[90m%s\033[0m\n\n"
			desc := i.Description

			if desc == "" {
				desc = "-"
			}

			buf.WriteString(fmt.Sprintf(format, i.Name, i.Default, desc))
		}

		buf.WriteString("\n")
	}

	if len(d.Outputs) > 0 {
		buf.WriteString("\n")

		for _, i := range d.Outputs {
			format := "  \033[36moutput.%s\033[0m\n  \033[90m%s\033[0m\n\n"
			s := fmt.Sprintf(format, i.Name, strings.TrimSpace(i.Description))
			buf.WriteString(s)
		}

		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// Template uses a txt/template to handle print of the documentation using a template-sample
func Template(templateName string, d doc.Doc, printRequired bool) (string, error) {
	templateFile, err := TemplateDir.Open(templateName + ".tmpl")
	if err != nil {
		log.Fatalln("Cannot open template", err)
	}
	defer templateFile.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(templateFile)

	templateContent := buf.String()

	return TemplateByString(templateContent, d, printRequired)
}

// TemplateByFile uses a txt/template to handle print of the documentation using a file on your disk
func TemplateByFile(templateFile string, d doc.Doc, printRequired bool) (string, error) {
	dat, err := ioutil.ReadFile(templateFile)
	if err != nil {
		log.Fatalln("Cannot open template-file", err)
	}

	return TemplateByString(string(dat), d, printRequired)
}

// TemplateByString uses a txt/template to handle print of the documentation using a string as template
func TemplateByString(templateContent string, d doc.Doc, printRequired bool) (string, error) {
	buf := new(bytes.Buffer)

	funcMap := sprig.FuncMap()
	funcMap["normalize"] = normalize
	funcMap["humanize"] = humanize
	funcMap["include"] = TemplateByFile

	tpl := template.New("printtemplate").Funcs(funcMap)

	printTemplate, err := tpl.Parse(templateContent)
	if err != nil {
		log.Fatalln("Cannot parse template", err)
	}

	buf.Reset()
	err = printTemplate.Execute(buf, templateStruct{
		Doc:           d,
		PrintRequired: printRequired,
	})

	return buf.String(), err
}

// JSON prints the given doc as json.
func JSON(d doc.Doc) (string, error) {
	s, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "", err
	}

	// <, >, and & are printed as code points by the json package.
	// The brackets are needed to pretty-print required_version.
	// Convert the brackets back into printable chars, limiting the
	// number of printed brackets to 1 each, which is enough to
	// prevent HTML injection (json's concern - why they encode).
	jsonString := strings.Replace(string(s), "\\u003c", "<", 1)
	jsonString = strings.Replace(jsonString, "\\u003e", ">", 1)

	return jsonString, nil
}

// Humanize the given `v`.
func humanize(def bool) string {
	if def {
		return "**yes**"
	}

	return "no"
}

// normalizeMarkdownDesc fixes line breaks in descriptions for markdown:
//
//  * Double newlines are converted to <br><br>
//  * A second pass replaces all other newlines with spaces
func normalizeMarkdownDesc(s string) string {
	return strings.Replace(strings.Replace(strings.TrimSpace(s), "\n\n", "<br><br>", -1), "\n", " ", -1)
}

// normalize prints out "-" for empty strings else does the same as normalizeMarkdownDesc
func normalize(s string) string {
	if s == "" {
		return "-"
	}
	return normalizeMarkdownDesc(s)
}

// Exists reports whether the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
