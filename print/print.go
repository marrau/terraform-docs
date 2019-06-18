package print

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"html/template"

	"github.com/Masterminds/sprig"
	"github.com/hashicorp/terraform/config"
)

type templateStruct struct {
	Config        config.Config
	PrintRequired bool
}

// Pretty printer pretty prints a doc.
func Pretty(cfg *config.Config) (string, error) {
	return Template("pretty", cfg, true)
}

// Template uses a txt/template to handle print of the documentation using a template-sample
func Template(templateName string, cfg *config.Config, printRequired bool) (string, error) {
	templateFile, err := TemplateDir.Open(templateName + ".tmpl")
	if err != nil {
		log.Fatalln("Cannot open template", err)
	}
	defer templateFile.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(templateFile)

	templateContent := buf.String()

	return TemplateByString(templateContent, cfg, printRequired)
}

// TemplateByFile uses a txt/template to handle print of the documentation using a file on your disk
func TemplateByFile(templateFile string, cfg *config.Config, printRequired bool) (string, error) {
	dat, err := ioutil.ReadFile(templateFile)
	if err != nil {
		log.Fatalln("Cannot open template-file", err)
	}

	return TemplateByString(string(dat), cfg, printRequired)
}

// TemplateByString uses a txt/template to handle print of the documentation using a string as template
func TemplateByString(templateContent string, cfg *config.Config, printRequired bool) (string, error) {
	buf := new(bytes.Buffer)

	funcMap := sprig.FuncMap()
	funcMap["normalize"] = normalize
	funcMap["humanize"] = humanize
	funcMap["include"] = TemplateByFile
	funcMap["html"] = htmlSafe
	funcMap["tfDocUrl"] = getTerraformDocumentationURL
	funcMap["relPath"] = relPath

	tpl := template.New("printtemplate").Funcs(funcMap)

	printTemplate, err := tpl.Parse(templateContent)
	if err != nil {
		log.Fatalln("Cannot parse template", err)
	}

	buf.Reset()
	err = printTemplate.Execute(buf, templateStruct{
		Config:        *cfg,
		PrintRequired: printRequired,
	})

	return buf.String(), err
}

func htmlSafe(s string) template.HTML {
	return template.HTML(s)
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
func normalize(in interface{}) interface{} {
	if s, ok := in.(string); ok {
		if s == "" {
			return "-"
		}
		return normalizeMarkdownDesc(s)
	}

	return in
}

func getTerraformDocumentationURL(object interface{}) string {
	if provider, ok := object.(*config.ProviderConfig); ok {
		return fmt.Sprintf("https://www.terraform.io/docs/providers/%s/index.html", strings.Replace(provider.Name, "-beta", "", -1))
	}
	if resource, ok := object.(*config.Resource); ok {
		rxp, err := regexp.Compile("[a-z]+")
		if err != nil {
			log.Fatal(err)
		}
		provider := rxp.FindString(resource.Type)
		typeName := strings.Replace(resource.Type, provider+"_", "", -1)
		resourceSubPath := "r"
		if resource.Mode == config.DataResourceMode {
			resourceSubPath = "d"
		}
		return fmt.Sprintf("https://www.terraform.io/docs/providers/%s/%s/%s.html", provider, resourceSubPath, typeName)
	}
	return ""
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

func relPath(dir string) string {
	p, err := filepath.Rel(path.Dir(dir), dir)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
