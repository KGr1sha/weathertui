package main

import (
	"strings"
	"slices"
)

type filterItem struct {
	word string
	distance int // levenshtein
}

func filter(words []string, filter string) []string {
	if len(words) == 0 {
		return []string{}
	}
	filter = strings.ToLower(filter)

	contains := startsWithFilter(words, filter)
	if len(contains) > 0 {
		return contains
	}

	withCommonRunes := withCommonRunes(words, filter)
	if len(withCommonRunes) == 0 {
		return []string{}
	}

	withDistances := make([]filterItem, len(withCommonRunes))
	for i, w := range withCommonRunes {
		withDistances[i] = filterItem{word: w, distance: levenshteinDistance(w, filter)}
	}

	return bestMatches(withDistances)
}

func containsFilter(words []string, word string) []string {
	var ans []string
	word = strings.ToLower(word)
	for _, w := range words {
		l := strings.ToLower(w)
		if strings.Contains(l, word) {
			ans = append(ans, w)
		}
	}
	return ans
}

func startsWithFilter(words []string, prefix string) []string {
	var ans []string
	prefix = strings.ToLower(prefix)
	for _, w := range words {
		lower  := strings.ToLower(w)
		if strings.HasPrefix(lower, prefix) {
			ans = append(ans, w)
		}
	}
	return ans
}

func bestMatches(items []filterItem) []string {
	slices.SortFunc(items, func(a, b filterItem) int {
		if a.distance < b.distance {
			return -1
		} else if a.distance > b.distance {
			return 1
		}
		return 0
	})

	var ans []string
	leastDistance := items[0].distance
	for _, item := range items {
		if item.distance == leastDistance {
			ans = append(ans, item.word)
		}
	}
	return ans
}


func levenshteinDistance(a, b string) int {
    la, lb := len(a), len(b)
    if la == 0 {
        return lb
    }
    if lb == 0 {
        return la
    }

    prev := make([]int, lb+1)
    curr := make([]int, lb+1)

    for j := 0; j <= lb; j++ {
        prev[j] = j
    }

    for i := 1; i <= la; i++ {
        curr[0] = i
        for j := 1; j <= lb; j++ {
            cost := 0
            if a[i-1] != b[j-1] {
                cost = 1
            }
            curr[j] = min(
                curr[j-1]+1,
                prev[j]+1,
                prev[j-1]+cost,
            )
        }
        prev, curr = curr, prev
    }

    return prev[lb]
}

func withCommonRunes(words []string, word string) []string {
	var ans []string
	word = strings.ToLower(word)
	targetSet := runeSet(word)
	for _, w := range words {
		if runeSet(strings.ToLower(w)) & targetSet != 0 {
			ans = append(ans, w)
		}
	}
	return ans
}

func runeSet(word string) uint64 {
	var set uint64
	for _, r := range word {
        if r >= 'a' && r <= 'z' {
            set |= 1 << (r - 'a')
        }
    }
    return set
}
