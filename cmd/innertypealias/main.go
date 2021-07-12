package main

import (
	"github.com/gostaticanalysis/innertypealias"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(innertypealias.Analyzer) }

