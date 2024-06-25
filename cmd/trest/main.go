package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/jasontconnell/trest/data"
	"github.com/jasontconnell/trest/process"
)

func main() {
	start := time.Now()
	tfile := flag.String("tfile", "tests.trst", "tests file")
	rootUrl := flag.String("url", "", "root url")
	outfile := flag.String("out", "out.txt", "output file")
	errorsfile := flag.String("errors", "errors.txt", "errors file")
	flag.Parse()

	if *tfile == "" || *rootUrl == "" {
		flag.PrintDefaults()
		return
	}

	root, err := process.ReadTests(*tfile)
	if err != nil || root == nil {
		log.Fatal("reading tests", err)
	}

	root.RootUrl = *rootUrl

	results := process.Run(root)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Duration < results[j].Duration
	})

	werr := writeResults(*outfile, results, func(r data.Result) bool {
		return r.Err == nil && r.Status == 200 && r.HasElement
	})

	eerr := writeResults(*errorsfile, results, func(r data.Result) bool {
		return r.Err != nil || r.Status != 200 || !r.HasElement
	})

	if werr != nil || eerr != nil {
		log.Println("problem writing file", werr, eerr)
	}

	log.Println("finished.", len(results), "Results.", time.Since(start))
}

func writeResults(filename string, results []data.Result, filter func(r data.Result) bool) error {
	fout, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("can't open output file %s. %w", filename, err)
	}
	defer fout.Close()

	for _, r := range results {
		if !filter(r) {
			continue
		}

		fmt.Fprintln(fout, r.Url, r.Body)
		fmt.Fprintf(fout, " Duration: %v\n", r.Duration)
		fmt.Fprintf(fout, " Status Code: %d\n", r.Status)
		if r.Err != nil {
			fmt.Fprintf(fout, " Error: %v\n", r.Err)
			fmt.Fprintf(fout, " Error Response \n%v\n\n", r.ErrResponse)
		}
	}
	return nil
}
