package main

import (
	"log"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{
			"./src/app/app.js",
		},
		EntryNames:        "[name].[hash]",
		Bundle:            true,
		Splitting:         true,
		Format:            api.FormatESModule,
		Write:             true,
		Outdir:            "../static/js",
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
	})

	if len(result.Errors) > 0 {
		log.Println(result)
		os.Exit(1)
	}
}
