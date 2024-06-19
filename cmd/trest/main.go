package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/jasontconnell/trest/process"
)

func main() {
	start := time.Now()
	tfile := flag.String("tfile", "tests.trst", "tests file")
	rootUrl := flag.String("url", "", "root url")
	outfile := flag.String("out", "out.txt", "output file")
	flag.Parse()

	if *tfile == "" || *rootUrl == "" {
		flag.PrintDefaults()
		return
	}

	root, err := process.ReadTests(*tfile)
	if err != nil {
		log.Fatal("reading tests", err)
	}

	root.RootUrl = *rootUrl

	results := process.Run(root)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Duration < results[j].Duration
	})

	fout, err := os.OpenFile(*outfile, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal("can't open output file ", *outfile)
	}
	defer fout.Close()

	for _, r := range results {
		fmt.Fprintln(fout, r.Url, r.Body)
		fmt.Fprintf(fout, " Duration: %v\n", r.Duration)
		fmt.Fprintf(fout, " Status Code: %d\n", r.Status)
		fmt.Fprintf(fout, " Error: %v\n", r.Err)
	}

	log.Println("finished.", time.Since(start))
}
