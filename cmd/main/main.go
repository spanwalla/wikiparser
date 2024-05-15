package main

import (
	"flag"
	"wikiparser/cmd/parser"
	"wikiparser/cmd/replacer"
)

func main() {
	input := flag.String("input", "", "path to the input file, required")
	page := flag.String("page", "", "path to the '*-page.sql.csv' file; if stated runs replacing titles to id, optional")
	silent := flag.Bool("silent", false, "silent mode, turn off the progress bar, optional")

	flag.Parse()

	if len(*input) == 0 {
		flag.Usage()
		return
	}

	if len(*page) > 0 {
		err := replacer.ReplaceTitleToId(*page, *input, *silent)
		if err != nil {
			return
		}
	} else {
		err := parser.ProcessFile(*input, *silent)
		if err != nil {
			return
		}
	}
}
