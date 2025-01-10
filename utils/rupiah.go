package utils

import (
	"fmt"
	"strings"
)

func FormatToRupiah(amount int) string {
	// Convert the integer to a string
	amountStr := fmt.Sprintf("%d", amount)

	// Add dots every three digits from the right
	var result strings.Builder
	length := len(amountStr)
	for i, digit := range amountStr {
		if (length-i)%3 == 0 && i != 0 {
			result.WriteRune('.')
		}
		result.WriteRune(digit)
	}

	// Prefix with "Rp"
	return "Rp" + result.String()
}
