package main

import (
	"flag"
	"wikiparser/internal/replace"
)

func main() {
	input := flag.String("input", "", "path to the '*-[pagelinks/redirect].sql.csv' file, required")
	page := flag.String("page", "", "path to the '*-page.sql.csv' file, required")
	silent := flag.Bool("silent", false, "silent mode, turn off the progress bar, optional")

	flag.Parse()

	if len(*input) == 0 || len(*page) == 0 {
		flag.Usage()
		return
	}

	err := replace.TitleToId(*page, *input, *silent)
	if err != nil {
		return
	}
}
