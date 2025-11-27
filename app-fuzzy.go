package main

import (
	"sort"
	"strings"
	"unicode"
)

func fuzzyScore(query, text string) (int, bool) {
	qi := 0
	score := 0
	lastMatch := -1

	for i, c := range text {
		if qi == len(query) {
			break
		}

		if unicode.ToLower(c) == unicode.ToLower(rune(query[qi])) {
			score += 10

			if lastMatch == i-1 {
				score += 5
			}

			if i == 0 || strings.ContainsRune("_-/ .", rune(text[i-1])) {
				score += 15
			}

			lastMatch = i
			qi++
		}
	}

	if qi != len(query) {
		return 0, false
	}

	score -= (lastMatch - (qi - 1))

	return score, true
}
func fuzzyFindApps(query string, apps []AppEntry) []AppEntry {
	type scored struct {
		score int
		entry AppEntry
	}

	var results []scored

	for _, app := range apps {
		// Cerca solo nel campo Name
		if s, ok := fuzzyScore(query, app.Name); ok {
			results = append(results, scored{
				score: s,
				entry: app,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		// Prima ordina per score (più alto è meglio)
		if results[i].score != results[j].score {
			return results[i].score > results[j].score
		}
		// Se lo score è uguale, ordina per lunghezza del nome (più corto è meglio)
		return len(results[i].entry.Name) < len(results[j].entry.Name)
	})

	out := make([]AppEntry, len(results))
	for i, r := range results {
		out[i] = r.entry
	}
	return out
}
