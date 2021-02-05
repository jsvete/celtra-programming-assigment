package main

import (
	"fmt"
	"sort"
	"strconv"
)

var selectedIds = map[int]struct{}{}

// selectedAccounts parses the existing selectedIds map
// and outputs a sorted slice of already selected account IDs.
func selectedAccounts() []int {
	ids := []int{}
	for id := range selectedIds {
		ids = append(ids, id)
	}

	sort.Slice(ids, func(i int, j int) bool {
		return ids[i] < ids[j]
	})

	return ids
}

// selectAccounts receives an array of strings then parses them then stores them to the selectedIds map.
// It will also return a sorted slice of all selected account IDs.
//
// If an ID is not an integer, it will print out the message with the invalid value.
func selectAccounts(ids ...string) []int {
	for _, strID := range ids {
		id, err := strconv.Atoi(strID)
		if err != nil || id < 1 {
			fmt.Printf(" %q should be an (positive) number\n", strID)

			continue
		}

		selectedIds[id] = struct{}{}
	}

	return selectedAccounts()
}
