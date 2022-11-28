package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dustin/go-humanize"
	"golang.org/x/exp/slices"
)

const vowelsRegexString = `[b-df-hj-np-tv-z]`

// The following regular expressions are used for pluralising words according
// to the following rules:
// - https://www.grammarly.com/blog/plural-nouns/
// - https://preply.com/en/blog/simple-rules-for-the-formation-of-plural-nouns-in-english/
var (
	replaceWithIRegex   = regexp.MustCompile(`(?i)us$`)
	replaceWithEsRegex  = regexp.MustCompile(`(?i)is$`)
	replaceWithVesRegex = regexp.MustCompile(`(?i)[^f]fe?$`)
	replaceWithOesRegex = regexp.MustCompile(
		`(?i)(` + vowelsRegexString + `)o$`,
	)
	replaceWithIesRegex = regexp.MustCompile(
		`(?i)(` + vowelsRegexString + `)y$`,
	)

	endsInEsRegex = regexp.MustCompile(`(?i)(s|sh|ch|x|z)$`)
)

// Possessive returns the target string appended with either 's or ' depending
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
// amount is not singular.
func Pluralise(target string, amount int64) string {
	if amount == 1 || target == "" {
		return target
	}

	return getPlural(target)
}

// PluraliseWithCount returns the pluralised version of targe, prefixed with
// the amount itself.
func PluraliseWithCount(target string, amount int64) string {
	return humanize.Comma(amount) + " " + Pluralise(target, amount)
}

// PluraliseSpecial returns the supplied singular string if amount is 1,
// else the supplied plural string is returned.
func PluraliseSpecial(singular, plural string, amount int64) string {
	if amount == 1 {
		return singular
	}

	return plural
}

func getPlural(target string) (plural string) {
	plural = replaceWithIRegex.ReplaceAllString(target, "i")
	if plural != target {
		return plural
	}

	plural = replaceWithEsRegex.ReplaceAllString(target, "es")
	if plural != target {
		return plural
	}

	plural = replaceWithVesRegex.ReplaceAllString(target, "ves")
	if plural != target {
		return plural
	}

	plural = replaceWithOesRegex.ReplaceAllString(target, "${1}oes")
	if plural != target {
		return plural
	}

	plural = replaceWithIesRegex.ReplaceAllString(target, "${1}ies")
	if plural != target {
		return plural
	}

	return target + PluralSuffix(target)
}

// PluralSuffix returns the plural suffix for a word. This misses most of the
// cases that Pluralise covers, as many cases require replace some ending
// characters of the word and don't just append a suffix.
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
			currentPage += fmt.Sprintln(line)
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
