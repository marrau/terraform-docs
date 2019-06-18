package print

import (
	"sort"

	"github.com/hashicorp/terraform/config"
)

type modulesByName []*config.Module

func (a modulesByName) Len() int           { return len(a) }
func (a modulesByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a modulesByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type resourcesByName []*config.Resource

func (a resourcesByName) Len() int           { return len(a) }
func (a resourcesByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a resourcesByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type providersByName []*config.ProviderConfig

func (a providersByName) Len() int           { return len(a) }
func (a providersByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a providersByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type inputsByName []*config.Variable

func (a inputsByName) Len() int           { return len(a) }
func (a inputsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a inputsByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type outputsByName []*config.Output

func (a outputsByName) Len() int           { return len(a) }
func (a outputsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a outputsByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type localsByName []*config.Local

func (a localsByName) Len() int           { return len(a) }
func (a localsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a localsByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type inputsByRequired []*config.Variable

func (a inputsByRequired) Len() int      { return len(a) }
func (a inputsByRequired) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a inputsByRequired) Less(i, j int) bool {
	switch {
	// i required, j not: i gets priority
	case a[i].Required() && !a[j].Required():
		return true
	// j required, i not: i does not get priority
	case !a[i].Required() && a[j].Required():
		return false
	// Otherwise, sort by name
	default:
		return a[i].Name < a[j].Name
	}
}

// Sort sorts the attributes to match  documentation order
func Sort(config *config.Config, sortByRequired bool) {
	switch {
	case sortByRequired:
		sort.Sort(inputsByRequired(config.Variables))
	default:
		sort.Sort(inputsByName(config.Variables))
	}

	sort.Sort(localsByName(config.Locals))
	sort.Sort(outputsByName(config.Outputs))
	sort.Sort(modulesByName(config.Modules))
	sort.Sort(providersByName(config.ProviderConfigs))
	sort.Sort(resourcesByName(config.Resources))
}
