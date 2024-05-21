package main

import (
	"log"
	"time"

	"github.com/hiimnhan/wiki-fts/common"
	"github.com/hiimnhan/wiki-fts/internal/indexing"
)

func main() {
	start := time.Now()
	master := indexing.NewMaster(10)
	master.Run(common.WikiDumpZipPath)
	log.Printf("All processes took %v", time.Since(start))
}
