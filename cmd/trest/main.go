package main

import (
	"flag"
	"log"
	"time"

	"github.com/jasontconnell/trest/process"
)

func main() {
	start := time.Now()
	tfile := flag.String("tfile", "tests.trst", "tests file")
	flag.Parse()

	if *tfile == "" {
		log.Fatal("no file")
	}

	root, err := process.ReadTests(*tfile)
	if err != nil {
		log.Fatal("reading tests", err)
	}

	err = process.Run(root, root)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("finished.", time.Since(start))
}
