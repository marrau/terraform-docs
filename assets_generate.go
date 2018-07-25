// +build generate

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	var fs http.FileSystem = http.Dir("templates")
	var opts = vfsgen.Options{
		Filename:     "print/template_vfsdata.go",
		PackageName:  "print",
		VariableName: "TemplateDir",
	}

	if err := vfsgen.Generate(fs, opts); err != nil {
		log.Fatalln(err)
	}
}
