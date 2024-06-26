package main

import (
	"flag"
	"wikiparser/cmd/combiner"
	"wikiparser/cmd/parser"
	"wikiparser/cmd/replacer"
	"wikiparser/cmd/validator"
)

func main() {
	input := flag.String("input", "", "path to the input file, required")
	page := flag.String("page", "", "path to the '*-page.sql.csv' file, required for 'replace' mode")
	redirect := flag.String("redirect", "", "path to the '*-redirect.sql.csv' file, required for 'validate' mode")
	mode := flag.String("mode", "", "mode to use [parse/replace/combine/validate], required")
	silent := flag.Bool("silent", false, "silent mode, turn off the progress bar, optional")

	flag.Parse()

	if len(*input) == 0 {
		flag.Usage()
		return
	}

	if *mode == "parse" {
		err := parser.ProcessFile(*input, *silent)
		if err != nil {
			return
		}
	} else if *mode == "replace" {
		if len(*page) == 0 {
			flag.Usage()
			return
		}
		err := replacer.ReplaceTitleToId(*page, *input, *silent)
		if err != nil {
			return
		}
	} else if *mode == "validate" {
		if len(*redirect) == 0 {
			flag.Usage()
			return
		}
		err := validator.ReplaceRedirect(*input, *redirect, *silent)
		if err != nil {
			return
		}
	} else if *mode == "combine" {
		err := combiner.CombineLinks(*input, *silent)
		if err != nil {
			return
		}
	} else {
		flag.Usage()
		return
	}
}
