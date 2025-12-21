package hw03frequencyanalysis

import (
	"slices"
	"sort"
	"strings"
	"unicode"
)

type topWordsCountType map[int][]string

const topN = 10

func Top10(text string) []string {
	if strings.TrimSpace(text) == "" {
		return []string{}
	}

	wordsCount := getWordsCount(text)
	topWordsCount := getTopWordsCount(wordsCount)
	counts := getCountsOfWords(topWordsCount)

	top := []string{}
	left := topN
	for i := len(counts) - 1; i >= 0 && i > len(counts)-topN; i-- {
		toTake := min(left, len(topWordsCount[counts[i]]))
		top = append(top, topWordsCount[counts[i]][0:toTake]...)
		left -= toTake
	}

	return top
}

func getWordsCount(text string) map[string]int {
	// I process text manually for the sake of memory efficiency
	separators := []rune{' ', '\n', '\r', '\t'}
	excludedString := []string{"", "-"}
	var window strings.Builder
	wordsCount := make(map[string]int)
	for _, r := range text {
		if !slices.Contains(separators, r) {
			window.WriteRune(r)
		} else if !slices.Contains(excludedString, window.String()) {
			normalizedWord := normalizeWord(window.String())
			if normalizedWord == "" {
				continue
			}
			wordsCount[normalizedWord]++
			window.Reset()
		}
	}
	wordsCount[normalizeWord(window.String())]++
	return wordsCount
}

func getTopWordsCount(wordsCount map[string]int) topWordsCountType {
	topWordsCount := make(map[int][]string)
	for key, value := range wordsCount {
		topWordsCount[value] = append(topWordsCount[value], key)
	}
	for key := range topWordsCount {
		sort.Strings(topWordsCount[key])
	}
	return topWordsCount
}

func getCountsOfWords(topWordsCount topWordsCountType) []int {
	counts := make([]int, len(topWordsCount))
	i := 0
	for key := range topWordsCount {
		counts[i] = key
		i++
	}
	sort.Ints(counts)
	return counts
}

func normalizeWord(word string) string {
	if word == "" {
		return ""
	}

	wordSlice := []rune(word)
	if len(wordSlice) == 1 {
		return strings.ToLower(word)
	}

	checkLeft := !unicode.IsLetter(wordSlice[0]) && !unicode.IsNumber(wordSlice[0])
	for checkLeft && len(wordSlice) > 1 {
		wordSlice = wordSlice[1:]
		checkLeft = !unicode.IsLetter(wordSlice[0]) && !unicode.IsNumber(wordSlice[0])
	}

	checkRight := !unicode.IsLetter(wordSlice[len(wordSlice)-1]) && !unicode.IsNumber(wordSlice[len(wordSlice)-1])
	for checkRight && len(wordSlice) > 1 {
		wordSlice = wordSlice[:len(wordSlice)-1]
		checkRight = !unicode.IsLetter(wordSlice[len(wordSlice)-1]) && !unicode.IsNumber(wordSlice[len(wordSlice)-1])
	}

	return strings.ToLower(string(wordSlice))
}
