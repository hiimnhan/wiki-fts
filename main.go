package main

import (
	"github.com/charmbracelet/log"
	"github.com/hiimnhan/wiki-fts/common"
	"github.com/kr/pretty"
)

func main() {
	docs, err := common.LoadDocuments(common.WikiDumpZipPath)
	if err != nil {
		log.Error(err)
	}

	pretty.Print(docs[1])

}
