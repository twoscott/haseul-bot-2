package dctools

import (
	"fmt"
	"regexp"
	"strings"
)

var markdownEscapeRegex = regexp.MustCompile("([\\|\\*\\`\\~\\_\\]\\)])")

// EscapeMarkdown returns a string with all instances if special Discord
// markdown characters escaped.
func EscapeMarkdown(unescaped string) string {
	return markdownEscapeRegex.ReplaceAllString(unescaped, "\\$1")
}

// MultiEscapeMarkdown returns a slice of the provided strings with their
// Discord markdown escaped.
func MultiEscapeMarkdown(unescapedS ...string) []string {
	escapedS := make([]string, len(unescapedS))
	for i, s := range unescapedS {
		escapedS[i] = EscapeMarkdown(s)
	}
	return escapedS
}

// Hyperlink returns a Discord markdown hyperlink in the form [name](url).
func Hyperlink(name string, url string) string {
	return fmt.Sprintf("[%s](%s)", name, url)
}

// ResizeImage appends the provided size to url.
// For use with Discord image URLs.
func ResizeImage(url string, size int) string {
	if size > 4096 {
		size = 4096
	}

	url = strings.Split(url, "?")[0]
	return fmt.Sprintf("%s?size=%d", url, size)
}
