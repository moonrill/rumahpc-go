package utils

import (
	"regexp"
	"strings"
)

func Slugify(s string) string {
	// Convert to lowercase
	slug := strings.ToLower(s)

	// Replace non-alphanumeric characters (excluding spaces) with a (-)
	re := regexp.MustCompile(`[^a-z0-9\s]+`)
	slug = re.ReplaceAllString(slug, "-")

	// Replace spaces with a (-)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Replace multiple (-) with a single (-)
	re = regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")

	// Remove leading and trailing (-)
	slug = strings.Trim(slug, "-")

	return slug
}
