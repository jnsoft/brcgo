package main

import (
	"github.com/brcgo/src/pipelines"
	"github.com/brcgo/src/util"
)

const (
	NO_OF_PARSER_WORKERS     = 2
	NO_OF_AGGREGATOR_WORKERS = 5
	ERROR                    = "❌ Error reading file"
	WARNING                  = "⚠️ Warning"
	DONE                     = "✅ Done"
)

func main() {
	fname := "testfile_10000.tmp"

	if false {
		util.GenerateFile(1000000, 1500, fname)
		return
	}

	pipelines.WorkerpoolPipeline(fname, NO_OF_AGGREGATOR_WORKERS, false)
	pipelines.ReadParseAggregatePipeline(fname, NO_OF_PARSER_WORKERS, NO_OF_AGGREGATOR_WORKERS, false)

}
