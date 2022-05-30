package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const MaxLength = 10

var regex = regexp.MustCompile(`[a-zA-Zа-яА-Я]+-?[a-zA-Zа-яА-Я]*`)

func Top10(text string) (res []string) {
	matches, top := regex.FindAllString(text, -1), MaxLength

	if matches == nil {
		return nil
	}

	wordFrequencies := make(map[string]int, len(matches))
	for _, word := range matches {
		wordFrequencies[strings.ToLower(word)]++
	}

	frequencyWords := make(map[int][]string, len(wordFrequencies))
	for word, frequency := range wordFrequencies {
		frequencyWords[frequency] = append(frequencyWords[frequency], word)
	}

	frequencies := make([]int, 0, len(frequencyWords))
	for frequency := range frequencyWords {
		frequencies = append(frequencies, frequency)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(frequencies)))
	for _, frequency := range frequencies {
		sort.Strings(frequencyWords[frequency])
		res = append(res, frequencyWords[frequency]...)
	}

	if len(wordFrequencies) < top {
		top = len(wordFrequencies)
	}

	return res[:top]
}
