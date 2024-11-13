//go:build tools
// +build tools

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"

	"github.com/openshift/installer-aro-wrapper/data"
)

func main() {
	err := vfsgen.Generate(data.Assets, vfsgen.Options{
		PackageName:  "data",
		BuildTags:    "release",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
