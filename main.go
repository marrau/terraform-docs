//go:generate go run assets_generate.go

package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/segmentio/terraform-docs/print"

	"github.com/tj/docopt"
)

var version = "dev"

const usage = `
  Usage:
    terraform-docs [md | markdown | tpl <template-path>] <path>
    terraform-docs -h | --help

  Examples:

    # View inputs and outputs
    $ terraform-docs ./my-module

    # Generate markdown tables of inputs and outputs
	$ terraform-docs md ./my-module

    # Generate templated output
    $ terraform-docs tpl path/to/template-file ./my-module		

  Options:
    -h, --help          Show help information
`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, true)
	if err != nil {
		log.Fatal(err)
	}

	module, diags := tfconfig.LoadModule(args["<path>"].(string))

	if diags.HasErrors() {
		log.Fatal(diags.Error())
	}

	var out string

	switch {
	case args["markdown"].(bool):
		out, err = print.Template("markdown", module)
	case args["md"].(bool):
		out, err = print.Template("markdown", module)
	case args["tpl"].(bool):
		templateName := args["<template-path>"].(string)
		out, err = print.TemplateByFile(templateName, module)
	default:
		out, err = print.Template("pretty", module)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
}
