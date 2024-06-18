package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jasontconnell/trest/process"
)

func main() {
	tfile := flag.String("tfile", "tests.trst", "tests file")
	flag.Parse()

	if *tfile == "" {
		log.Fatal("no file")
	}

	tconfig, err := process.ReadTests(*tfile)
	if err != nil {
		log.Fatal("reading tests", err)
	}

	fmt.Println(tconfig)
}
