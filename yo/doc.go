package main

// Gor - Fastest Static Blog Engine

import (
	"fmt"
	"os"
)

const (
	HELP = `
Run specified yo tool

Usage:

    yo command [args...]

        -debug      Enable debug mode

    yo new <diretory>                  Init Blog layout
    yo compile                         Compile
    yo http                            Preview Compiled Website
    yo config                          Print Configure
    yo payload                         Print Payload
    yo pprof                           Run pprof (for dev)
    yo post <title> {image diretory}   Post new page
    yo posts                           Show all pages
`
)

func PrintUsage() {
	fmt.Println(HELP)
	os.Exit(1)
}
