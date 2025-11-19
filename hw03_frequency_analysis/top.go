package hw03frequencyanalysis

import (
	"slices"
	"sort"
	"strings"
)

const topN = 10

func Top10(text string) []string {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	separators := []rune{' ', '\n', '\t'}

	wordsCount := make(map[string]int)
	var window strings.Builder
	for _, r := range text {
		if !slices.Contains(separators, r) {
			window.WriteRune(r)
		} else if window.String() != "" {
			wordsCount[window.String()]++
			window.Reset()
		}
	}

	topWordsCount := make(map[int]string)
	for key, value := range wordsCount {
		if _, ok := topWordsCount[value]; !ok {
			topWordsCount[value] = key
			continue
		}
		if topWordsCount[value] > key {
			topWordsCount[value] = key
		}
	}

	counts := make([]int, len(topWordsCount))
	i := 0
	for key := range topWordsCount {
		counts[i] = key
		i++
	}
	sort.Ints(counts)

	var top []string
	for i := len(counts) - 1; i >= 0 && i > len(counts)-topN; i-- {
		top = append(top, topWordsCount[counts[i]])
	}

	return top
}
