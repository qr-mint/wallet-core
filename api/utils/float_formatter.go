package utils

import (
	"fmt"
	"strings"
)

func FormatFloat(value float64) string {
	if value == 0 {
		return "0.0"
	}
	formatted := fmt.Sprintf("%.9f", value)
	formatted = strings.TrimRight(formatted, "0")
	if formatted[len(formatted)-1] == '.' {
		formatted += "0"
	}
	return formatted
}

func FormatFloatFiat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}
