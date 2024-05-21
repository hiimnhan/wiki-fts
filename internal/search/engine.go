package search

import (
	"log"

	"github.com/hiimnhan/wiki-fts/common"
)

type Index map[string][]int

func (idx *Index) FindIndexes(text string) []int {
	log.Printf("Querying %s\n....", text)
	var res []int
	for _, word := range common.TokenizeAndFilter(text) {
		if ids, ok := (*idx)[word]; ok {
			if len(res) == 0 {
				res = ids
			} else {
				res = common.Intersection(res, ids)
			}
		} else {
			return nil
		}
	}

	return res
}
