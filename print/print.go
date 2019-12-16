package print

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"html/template"

	"github.com/Masterminds/sprig"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

// Template uses a txt/template to handle print of the documentation using a template-sample
func Template(templateName string, cfg *tfconfig.Module) (string, error) {
	templateFile, err := TemplateDir.Open(templateName + ".tmpl")
	if err != nil {
		log.Fatalln("Cannot open template", err)
	}
	defer templateFile.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(templateFile)

	templateContent := buf.String()

	return TemplateByString(templateContent, cfg)
}

// TemplateByFile uses a txt/template to handle print of the documentation using a file on your disk
func TemplateByFile(templateFile string, cfg *tfconfig.Module) (string, error) {
	dat, err := ioutil.ReadFile(templateFile)
	if err != nil {
		log.Fatalln("Cannot open template-file", err)
	}

	return TemplateByString(string(dat), cfg)
}

// TemplateByString uses a txt/template to handle print of the documentation using a string as template
func TemplateByString(templateContent string, cfg *tfconfig.Module) (string, error) {
	buf := new(bytes.Buffer)

	funcMap := sprig.FuncMap()
	funcMap["normalize"] = normalize
	funcMap["humanize"] = humanize
	funcMap["include"] = TemplateByFile
	funcMap["html"] = htmlSafe
	funcMap["tfDocUrl"] = getTerraformDocumentationURL
	funcMap["fileExists"] = fileExists
	funcMap["relPath"] = relPath
	funcMap["guessType"] = guessType
	funcMap["sortByRequired"] = sortByRequired

	tpl := template.New("printtemplate").Funcs(funcMap)

	printTemplate, err := tpl.Parse(templateContent)
	if err != nil {
		log.Fatalln("Cannot parse template", err)
	}

	buf.Reset()
	err = printTemplate.Execute(buf, cfg)

	return buf.String(), err
}

func sortByRequired(v map[string]*tfconfig.Variable) (sorted []*tfconfig.Variable) {
	for _, val := range v {
		sorted = append(sorted, val)
	}
	sort.Sort(inputsByRequired(sorted))
	return sorted
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
	return strings.Replace(strings.TrimSpace(s), "\n", "<br>", -1)
}

// normalize prints out "-" for empty strings else does the same as normalizeMarkdownDesc
func normalize(in interface{}) interface{} {
	if s, ok := in.(string); ok {
		if s == "" {
			return "-"
		}
		return template.HTML(normalizeMarkdownDesc(s))
	}

	return in
}

func getTerraformDocumentationURL(object interface{}) string {
	if provider, ok := object.(*tfconfig.ProviderRef); ok {
		return fmt.Sprintf("https://www.terraform.io/docs/providers/%s/index.html", strings.Replace(provider.Name, "-beta", "", -1))
	}
	if resource, ok := object.(*tfconfig.Resource); ok {
		rxp, err := regexp.Compile("[a-z]+")
		if err != nil {
			log.Fatal(err)
		}
		provider := rxp.FindString(resource.Type)
		typeName := strings.Replace(resource.Type, provider+"_", "", -1)
		resourceSubPath := "r"
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

type typeGuess struct {
	Name     string
	Value    interface{}
	Required bool
	Type     string
	RawType  string
	Nested   []typeGuess
}

var objReplaceRGX = regexp.MustCompile(`(?ms)^object\({(.*)}\)$`)
var keyTypeRGX = regexp.MustCompile(`(?m)([\S]+)[\s]*=[\s]*([\S]+)`)

func inspectObjectType(obj string, def interface{}) typeGuess {
	nested := []typeGuess{}

	newObj := objReplaceRGX.ReplaceAllString(obj, "$1")

	for _, match := range keyTypeRGX.FindAllStringSubmatch(newObj, -1) {
		var defValue interface{}
		if dv, ok := def.(map[string]interface{}); ok {
			defValue = dv[match[1]]
		}
		nested = append(nested, typeGuess{
			Name:     match[1],
			Value:    defValue,
			Type:     match[2],
			Required: defValue == nil,
		})
	}

	return typeGuess{
		Value:    def,
		Type:     "object",
		RawType:  obj,
		Nested:   nested,
		Required: def == nil,
	}
}

func guessType(v interface{}) typeGuess {
	if variable, ok := v.(*tfconfig.Variable); ok {
		if variable.Type != "" {
			if strings.HasPrefix(variable.Type, "object") {
				return inspectObjectType(variable.Type, variable.Default)
			}
			return typeGuess{
				Value:    variable.Default,
				Type:     variable.Type,
				Required: variable.Default == nil,
			}
		}

		if _, ok := variable.Default.(map[string]interface{}); ok {
			return typeGuess{
				Value:    variable.Default,
				Type:     "map",
				Required: variable.Default == nil,
			}
		}
		if _, ok := variable.Default.([]interface{}); ok {
			return typeGuess{
				Value:    variable.Default,
				Type:     "list",
				Required: variable.Default == nil,
			}
		}

		return typeGuess{
			Value:    variable.Default,
			Type:     fmt.Sprint(reflect.TypeOf(variable.Default)),
			Required: variable.Default == nil,
		}
	} else if _, ok := v.(string); ok {
		return typeGuess{
			Value:    v,
			Type:     "string",
			Required: v == nil,
		}
	}

	return typeGuess{
		Value:    v,
		Type:     fmt.Sprint(reflect.TypeOf(v)),
		Required: v == nil,
	}
}
