package search

import (
	"github.com/hiimnhan/wiki-fts/common"
)

type SearchEngine struct {
	Index    common.Index
	DocsDict common.DocumentDict
}
