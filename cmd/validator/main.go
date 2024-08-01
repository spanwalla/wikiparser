package main

import (
	"flag"
	"wikiparser/internal/validate"
)

func main() {
	pagelinks := flag.String("pagelinks", "", "path to the '*-pagelinks.sql.csv' file, required")
	redirect := flag.String("redirect", "", "path to the '*-redirect.sql.csv' file, required")
	silent := flag.Bool("silent", false, "silent mode, turn off the progress bar, optional")

	flag.Parse()

	if len(*pagelinks) == 0 || len(*redirect) == 0 {
		flag.Usage()
		return
	}

	err := validate.Redirect(*pagelinks, *redirect, *silent)
	if err != nil {
		return
	}
}
