package main

import (
	"flag"
	"wikiparser/internal/combine"
)

func main() {
	pagelinks := flag.String("pagelinks", "", "path to the '*-pagelinks.sql.csv' file, required")
	silent := flag.Bool("silent", false, "silent mode, turn off the progress bar, optional")

	flag.Parse()

	if len(*pagelinks) == 0 {
		flag.Usage()
		return
	}

	err := combine.Links(*pagelinks, *silent)
	if err != nil {
		return
	}
}
