package util

import (
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

var endsInEsRegex = regexp.MustCompile(`\S+(s|sh|ch|x|z)$`)

// Possessive returns the target string appended with either ' or 's depending
// on the string's postfix character.
func Possessive(target string) string {
	if target == "" {
		return target
	}

	return target + PossessiveSuffix(target)
}

// PossessiveSuffix return the possessive suffix for a word.
func PossessiveSuffix(target string) (suffix string) {
	if target == "" {
		return ""
	}

	lastChar := target[len(target)-1]
	if lastChar == 's' || lastChar == 'S' {
		suffix = "'"
	} else {
		suffix = "'s"
	}

	return suffix
}

// Pluralise adds the grammar-correct plural suffix to target if
// amount is not singular
func Pluralise(target string, amount int64) string {
	if amount == 1 || target == "" {
		return target
	}

	return target + PluralSuffix(target)
}

// PluralSuffix returns the plural suffix for a word.
func PluralSuffix(target string) (suffix string) {
	if endsInEsRegex.MatchString(target) {
		suffix = "es"
	} else {
		suffix = "s"
	}

	return suffix
}

// PagedLines returns a slice of pages where the given lines are added to pages
// separated by newlines until they either overflow the character limit or
// the line limit.
func PagedLines(lines []string, limit int, lineLimit int) []string {
	approxPages := (len(lines) / lineLimit) + 1
	pages := make([]string, 0, approxPages)

	linesAdded := 0
	currentPage := ""
	for _, line := range lines {
		if len(currentPage)+len(line) < limit && linesAdded < lineLimit {
			currentPage += line + "\n"
			linesAdded++
		} else {
			pages = append(pages, currentPage)
			currentPage = ""
			linesAdded = 0
		}
	}
	pages = append(pages, currentPage)

	return pages
}

// TrimArgs returns a string where the number of args denoted by limit are
// trimmed from the beginning of the string.
func TrimArgs(content string, limit int) string {
	args := strings.Fields(content)
	for _, arg := range args[:limit] {
		splitIndex := strings.Index(content, arg) + len(arg)
		content = content[splitIndex:]
	}

	return strings.TrimSpace(content)
}

// SearchSort takes a slice of string results and a query to search for, and
// then filters and sorts them accordingly.
func SearchSort(results []string, query string) []string {
	matches := make([]string, 0, len(results))
	for _, r := range results {
		if strings.Contains(r, query) {
			matches = append(matches, r)
		}
	}

	slices.Sort(matches)
	matches = slices.Compact(matches)
	slices.SortStableFunc(matches, func(a, b string) bool {
		return len(a) < len(b)
	})
	slices.SortStableFunc(matches, func(a, b string) bool {
		return strings.Index(a, query) < strings.Index(b, query)
	})

	return matches
}
