package main

import (
	"flag"
	"slices"

	"log"

	// "runtime"
	"time"

	"github.com/hiimnhan/wiki-fts/common"
	"github.com/hiimnhan/wiki-fts/internal/indexing"
)

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU())
	start := time.Now()
	master := indexing.NewMaster(6)
	master.Run(common.WikiDumpPath)
	log.Printf("All processes took %v", time.Since(start))

	records, err := common.ReadIndexFromFile(common.OutputPath)
	if err != nil {
		log.Fatal(err)
	}
	var query string
	flag.StringVar(&query, "q", "", "search query")
	flag.Parse()

	start = time.Now()
	matchedIds := records.FindIndexes(query)

	slices.Sort(matchedIds)

	// for _, id := range matchedIds {
	// 	doc := (*docsDict)[id]
	// 	fmt.Println(doc.Display())
	// }

	log.Printf("Search took %v, found %d, %v", time.Since(start), len(matchedIds), matchedIds)
}
