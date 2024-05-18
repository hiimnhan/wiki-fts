package main

import (
	"github.com/hiimnhan/wiki-fts/common"
	"github.com/hiimnhan/wiki-fts/internal/indexing"
)

func main() {
	master := indexing.NewMaster(10)
	master.Run(common.WikiDumpZipPath)
}
