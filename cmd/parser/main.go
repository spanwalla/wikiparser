package main

import (
	"flag"
	"wikiparser/internal/parse"
)

func main() {
	input := flag.String("input", "", "path to the .sql.gz file, required")
	silent := flag.Bool("silent", false, "silent mode, turn off the progress bar, optional")

	flag.Parse()

	if len(*input) == 0 {
		flag.Usage()
		return
	}

	err := parse.ProcessFile(*input, *silent)
	if err != nil {
		return
	}
}
