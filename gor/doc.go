// Gor - Fastest Static Blog Engine
package main

import (
	"fmt"
	"os"
)

const (
	HELP = `
Run specified gor tool

Usage:

    gor command [args...]

        -debug      Enable debug mode

    gor new <diretory>                  Init Blog layout
    gor compile                         Compile
    gor http                            Preview Compiled Website
    gor config                          Print Configure
    gor payload                         Print Payload
    gor pprof                           Run pprof (for dev)
    gor post <title> {image diretory}   Post new page
    gor posts                           Show all pages
`
)

func PrintUsage() {
	fmt.Println(HELP)
	os.Exit(1)
}
