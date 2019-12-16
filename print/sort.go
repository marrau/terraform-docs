package print

import (
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

type resourcesByName []*tfconfig.Resource

func (a resourcesByName) Len() int           { return len(a) }
func (a resourcesByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a resourcesByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type inputsByName []*tfconfig.Variable

func (a inputsByName) Len() int           { return len(a) }
func (a inputsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a inputsByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type outputsByName []*tfconfig.Output

func (a outputsByName) Len() int           { return len(a) }
func (a outputsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a outputsByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type inputsByRequired []*tfconfig.Variable

func (a inputsByRequired) Len() int      { return len(a) }
func (a inputsByRequired) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a inputsByRequired) Less(i, j int) bool {
	ai := a[i].Default != nil
	aj := a[j].Default != nil

	switch {
	// i required, j not: i gets priority
	case ai && !aj:
		return false
	// j required, i not: i does not get priority
	case !ai && aj:
		return true
	// Otherwise, sort by name
	default:
		return a[i].Name < a[j].Name
	}
}
