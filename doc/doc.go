package doc

import (
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
)

// Module represents a terraform module block.
type Module struct {
	Name        string
	Description string
	Source      string
}

// Provider represents a terraform provider block.
type Provider struct {
	Name          string
	Description   string
	Documentation string
	Version       string
}

// Resource represents a terraform resource block.
type Resource struct {
	Name          string
	Type          string
	Description   string
	Documentation string
}

// Input represents a terraform input variable.
type Input struct {
	Name        string
	Description string
	Default     string
	Type        string
	Required    bool
}

// Value returns the default value as a string.
func value(v *Value) string {
	if v != nil {
		switch v.Type {
		case "string":
			return v.Literal
		case "map":
			return "<map>"
		case "list":
			return "<list>"
		}
	}

	return "-"
}

// Value represents a terraform value.
type Value struct {
	Type    string
	Literal string
}

// Output represents a terraform output.
type Output struct {
	Name        string
	Description string
}

// Doc represents a terraform module doc.
type Doc struct {
	Version   string
	Comment   string
	Modules   []Module
	Providers []Provider
	Resources []Resource
	Inputs    []Input
	Outputs   []Output
}

type modulesByName []Module

func (a modulesByName) Len() int           { return len(a) }
func (a modulesByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a modulesByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type resourcesByName []Resource

func (a resourcesByName) Len() int           { return len(a) }
func (a resourcesByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a resourcesByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type providersByName []Provider

func (a providersByName) Len() int           { return len(a) }
func (a providersByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a providersByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type inputsByName []Input

func (a inputsByName) Len() int           { return len(a) }
func (a inputsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a inputsByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type outputsByName []Output

func (a outputsByName) Len() int           { return len(a) }
func (a outputsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a outputsByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type inputsByRequired []Input

func (a inputsByRequired) Len() int      { return len(a) }
func (a inputsByRequired) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a inputsByRequired) Less(i, j int) bool {
	switch {
	// i required, j not: i gets priority
	case a[i].Required && !a[j].Required:
		return true
	// j required, i not: i does not get priority
	case !a[i].Required && a[j].Required:
		return false
	// Otherwise, sort by name
	default:
		return a[i].Name < a[j].Name
	}
}

// Create creates a new *Doc from the supplied map
// of filenames and *ast.File.
func Create(files map[string]*ast.File, sortByRequired bool) Doc {
	doc := new(Doc)

	for name, f := range files {
		list := f.Node.(*ast.ObjectList)

		requiredVersion := version(list)
		if len(requiredVersion) > 0 {
			doc.Version = requiredVersion
		}

		doc.Providers = append(doc.Providers, providers(list)...)
		doc.Resources = append(doc.Resources, resources(list)...)
		doc.Inputs = append(doc.Inputs, inputs(list)...)
		doc.Outputs = append(doc.Outputs, outputs(list)...)
		doc.Modules = append(doc.Modules, modules(list)...)

		filename := path.Base(name)
		comments := f.Comments

		if filename == "main.tf" && len(comments) > 0 {
			doc.Comment = header(comments[0])
		}
	}

	switch {
	case sortByRequired:
		sort.Sort(inputsByRequired(doc.Inputs))
	default:
		sort.Sort(inputsByName(doc.Inputs))
	}
	sort.Sort(outputsByName(doc.Outputs))
	sort.Sort(modulesByName(doc.Modules))
	sort.Sort(providersByName(doc.Providers))
	sort.Sort(resourcesByName(doc.Resources))
	return *doc
}

func modules(list *ast.ObjectList) []Module {
	var ret []Module

	for _, item := range list.Items {
		if is(item, "module") {
			name, _ := strconv.Unquote(item.Keys[1].Token.Text)
			if name == "" {
				name = item.Keys[1].Token.Text
			}
			items := item.Val.(*ast.ObjectType).List.Items
			var desc string
			switch {
			case description(items) != "":
				desc = description(items)
			case item.LeadComment != nil:
				desc = comment(item.LeadComment.List)
			}
			source := get(items, "source")
			ret = append(ret, Module{
				Name:        name,
				Description: desc,
				Source:      strings.TrimSpace(source.Literal),
			})
		}
	}

	return ret
}

// Version returns the terraform version_required string from 'list'.
func version(list *ast.ObjectList) string {
	var ret string

	for _, item := range list.Items {
		if is(item, "terraform") && item.Val.(*ast.ObjectType).List.Items[0].Keys[0].Token.Text == "required_version" {
			version := item.Val.(*ast.ObjectType).List.Items[0].Val.(*ast.LiteralType).Token.Text
			version = strings.Trim(version, "\"")

			if len(version) > 0 {
				ret = version
			}
		}
	}

	return ret
}

// Providers returns all providers from 'list' along with links
// to their Terraform documentation.
func providers(list *ast.ObjectList) []Provider {
	var ret []Provider
	var version = "Latest"

	for _, item := range list.Items {
		if is(item, "provider") {
			name := item.Keys[1].Token.Text
			name = strings.Trim(name, "\"")
			link := fmt.Sprintf("https://www.terraform.io/docs/providers/%s", name)

			items := item.Val.(*ast.ObjectType).List.Items
			var desc string
			switch {
			case description(items) != "":
				desc = description(items)
			case item.LeadComment != nil:
				desc = comment(item.LeadComment.List)
			}
			if v := get(items, "description"); v != nil {
				version = v.Literal
			}

			ret = append(ret, Provider{
				Name:          name,
				Description:   desc,
				Documentation: link,
				Version:       strings.TrimSpace(version),
			})
		}
	}

	return ret
}

// Resources returns all resources from 'list' along with links
// to their Terraform documentation.
func resources(list *ast.ObjectList) []Resource {
	var ret []Resource

	for _, item := range list.Items {
		if is(item, "resource") {
			name := item.Keys[2].Token.Text
			name = strings.Trim(name, "\"")

			resourceType := item.Keys[1].Token.Text
			resourceType = strings.Trim(resourceType, "\"")

			resourceTypes := strings.SplitN(resourceType, "_", 2)
			namespace := resourceTypes[0]
			typestr := resourceTypes[1]
			link := fmt.Sprintf("https://www.terraform.io/docs/providers/%s/r/%s.html", namespace, typestr)

			items := item.Val.(*ast.ObjectType).List.Items
			var desc string
			switch {
			case description(items) != "":
				desc = description(items)
			case item.LeadComment != nil:
				desc = comment(item.LeadComment.List)
			}

			ret = append(ret, Resource{
				Name:          name,
				Type:          resourceType,
				Documentation: link,
				Description:   desc,
			})
		}
	}

	return ret
}

// Inputs returns all variables from `list`.
func inputs(list *ast.ObjectList) []Input {
	var ret []Input

	for _, item := range list.Items {
		if is(item, "variable") {
			name, _ := strconv.Unquote(item.Keys[1].Token.Text)
			if name == "" {
				name = item.Keys[1].Token.Text
			}
			items := item.Val.(*ast.ObjectType).List.Items
			var desc string
			switch {
			case description(items) != "":
				desc = description(items)
			case item.LeadComment != nil:
				desc = comment(item.LeadComment.List)
			}

			var itemsType = get(items, "type")
			var itemType string

			if itemsType == nil || itemsType.Literal == "" {
				itemType = "string"
			} else {
				itemType = itemsType.Literal
			}

			def := get(items, "default")
			ret = append(ret, Input{
				Name:        name,
				Description: desc,
				Default:     value(def),
				Type:        itemType,
				Required:    def == nil,
			})
		}
	}

	return ret
}

// Outputs returns all outputs from `list`.
func outputs(list *ast.ObjectList) []Output {
	var ret []Output

	for _, item := range list.Items {
		if is(item, "output") {
			name, _ := strconv.Unquote(item.Keys[1].Token.Text)
			if name == "" {
				name = item.Keys[1].Token.Text
			}
			items := item.Val.(*ast.ObjectType).List.Items
			var desc string
			switch {
			case description(items) != "":
				desc = description(items)
			case item.LeadComment != nil:
				desc = comment(item.LeadComment.List)
			}

			ret = append(ret, Output{
				Name:        name,
				Description: strings.TrimSpace(desc),
			})
		}
	}

	return ret
}

// Get `key` from the list of object `items`.
func get(items []*ast.ObjectItem, key string) *Value {
	for _, item := range items {
		if is(item, key) {
			v := new(Value)

			if lit, ok := item.Val.(*ast.LiteralType); ok {
				if value, ok := lit.Token.Value().(string); ok {
					v.Literal = value
				} else {
					v.Literal = lit.Token.Text
				}
				v.Type = "string"
				return v
			}

			if _, ok := item.Val.(*ast.ObjectType); ok {
				v.Type = "map"
				return v
			}

			if _, ok := item.Val.(*ast.ListType); ok {
				v.Type = "list"
				return v
			}

			return nil
		}
	}

	return nil
}

// description returns a description from items or an empty string.
func description(items []*ast.ObjectItem) string {
	if v := get(items, "description"); v != nil {
		return v.Literal
	}

	return ""
}

// Is returns true if `item` is of `kind`.
func is(item *ast.ObjectItem, kind string) bool {
	if len(item.Keys) > 0 {
		return item.Keys[0].Token.Text == kind
	}

	return false
}

// Unquote the given string.
func unquote(s string) string {
	s, _ = strconv.Unquote(s)
	return s
}

// Comment cleans and returns a comment.
func comment(l []*ast.Comment) string {
	var line string
	var ret string

	for _, t := range l {
		line = strings.TrimSpace(t.Text)
		line = strings.TrimPrefix(line, "#")
		line = strings.TrimPrefix(line, "//")
		ret += strings.TrimSpace(line) + "\n"
	}

	return ret
}

// Header returns the header comment from the list
// or an empty comment. The head comment must start
// at line 1 and start with `/**`.
func header(c *ast.CommentGroup) (comment string) {
	if len(c.List) == 0 {
		return comment
	}

	if c.Pos().Line != 1 {
		return comment
	}

	cm := strings.TrimSpace(c.List[0].Text)

	if strings.HasPrefix(cm, "/**") {
		lines := strings.Split(cm, "\n")

		if len(lines) < 2 {
			return comment
		}

		lines = lines[1 : len(lines)-1]
		for _, l := range lines {
			l = strings.TrimSpace(l)
			switch {
			case strings.TrimPrefix(l, "* ") != l:
				l = strings.TrimPrefix(l, "* ")
			default:
				l = strings.TrimPrefix(l, "*")
			}
			comment += l + "\n"
		}
	}

	return comment
}
