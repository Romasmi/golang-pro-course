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

	wordsCount := getWordsCount(text)
	topWordsCount := getTopWordsCount(wordsCount)
	counts := getCountsOfWords(topWordsCount)

	var top []string
	for i := len(counts) - 1; i >= 0 && i > len(counts)-topN; i-- {
		top = append(top, topWordsCount[counts[i]])
	}

	return top
}

func getWordsCount(text string) map[string]int {
	separators := []rune{' ', '\n', '\t'}
	var window strings.Builder
	wordsCount := make(map[string]int)
	for _, r := range text {
		if !slices.Contains(separators, r) {
			window.WriteRune(r)
		} else if window.String() != "" {
			wordsCount[window.String()]++
			window.Reset()
		}
	}
	return wordsCount
}

func getTopWordsCount(wordsCount map[string]int) map[int]string {
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
	return topWordsCount
}

func getCountsOfWords(topWordsCount map[int]string) []int {
	counts := make([]int, len(topWordsCount))
	i := 0
	for key := range topWordsCount {
		counts[i] = key
		i++
	}
	sort.Ints(counts)
	return counts
}
