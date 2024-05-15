package parser

import "regexp"

var tableInfo = map[string]struct {
	re      *regexp.Regexp
	columns []int
}{
	"page":      {re: regexp.MustCompile("^(\\d+),0,('.*'),([01]),([01]),([0-9.]+),('\\d+'),('\\d+'|NULL),(\\d+),(\\d+),('.*'|NULL),('.*'|NULL)$"), columns: []int{0, 1, 2}},
	"pagelinks": {re: regexp.MustCompile("^(\\d+),0,('.*'),0,(\\d+|NULL)$"), columns: []int{0, 1}},
	"redirect":  {re: regexp.MustCompile("^(\\d+),0,('.*'),('.*'|NULL),('.*'|NULL)$"), columns: []int{0, 1}},
}
