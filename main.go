//go:generate go run assets_generate.go

package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/config"
	"github.com/segmentio/terraform-docs/print"

	"github.com/tj/docopt"
)

var version = "dev"

const usage = `
  Usage:
    terraform-docs [--no-required] [md | markdown | tpl <template-path>] <path>
    terraform-docs [--sort-by-required] [md | markdown | tpl <template-path>] <path>
    terraform-docs -h | --help

  Examples:

    # View inputs and outputs
    $ terraform-docs ./my-module

    # View inputs and outputs for variables.tf and outputs.tf only
    $ terraform-docs variables.tf outputs.tf

    # Generate markdown tables of inputs and outputs
	$ terraform-docs md ./my-module

    # Generate markdown tables of inputs and outputs, but don't print "Required" column
    $ terraform-docs --no-required md ./my-module

    # Generate markdown tables of inputs and outputs for the given module and ../config.tf
	$ terraform-docs md ./my-module ../config.tf
	
    # Generate templated output
    $ terraform-docs tpl path/to/template-file ./my-module		

  Options:
    -h, --help          Show help information
    --no-required       Generate markdown tables of inputs and outputs, but don't print "Required" column
    --sort-by-required  Sort the required inputs to the top of the output table
`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, true)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.LoadDir(args["<path>"].(string))
	if err != nil {
		log.Fatal(err)
	}

	sortByRequired := args["--sort-by-required"].(bool)
	print.Sort(cfg, sortByRequired)

	printRequired := !args["--no-required"].(bool)

	var out string

	switch {
	case args["markdown"].(bool):
		out, err = print.Template("markdown", cfg, printRequired)
	case args["md"].(bool):
		out, err = print.Template("markdown", cfg, printRequired)
	case args["tpl"].(bool):
		templateName := args["<template-path>"].(string)
		out, err = print.TemplateByFile(templateName, cfg, printRequired)
	default:
		out, err = print.Template("pretty", cfg, printRequired)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
}
