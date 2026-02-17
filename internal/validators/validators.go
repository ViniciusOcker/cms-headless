package validators

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func SanitizeHTML(input string) string {
	p := bluemonday.UGCPolicy() // Protege contra XSS
	return p.Sanitize(input)
}

func GenerateSlug(title string) string {
	return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}
