package utils

import (
	"fmt"
	"math"
)

// FormatNumber formats large numbers with K, M, B, T suffixes.
func FormatNumber(n float64) string {
	if n < 0 {
		return "-" + FormatNumber(-n)
	}

	suffixes := []string{"", "K", "M", "B", "T", "Q"}

	if n < 1000 {
		if n == float64(int64(n)) {
			return fmt.Sprintf("%d", int64(n))
		}
		return fmt.Sprintf("%.1f", n)
	}

	exp := int(math.Log10(n) / 3)
	if exp >= len(suffixes) {
		exp = len(suffixes) - 1
	}

	value := n / math.Pow(1000, float64(exp))

	if value >= 100 {
		return fmt.Sprintf("%.0f%s", value, suffixes[exp])
	} else if value >= 10 {
		return fmt.Sprintf("%.1f%s", value, suffixes[exp])
	}
	return fmt.Sprintf("%.2f%s", value, suffixes[exp])
}

// FormatMana formats mana with appropriate suffix.
func FormatMana(mana float64) string {
	return FormatNumber(mana) + " mana"
}

// FormatPercent formats a decimal as a percentage.
func FormatPercent(value float64) string {
	return fmt.Sprintf("%.0f%%", value*100)
}

// FormatMultiplier formats a multiplier value.
func FormatMultiplier(value float64) string {
	return fmt.Sprintf("%.2fx", value)
}

// Clamp constrains a value between min and max.
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampInt constrains an integer between min and max.
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Min returns the minimum of two floats.
func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two floats.
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// MinInt returns the minimum of two integers.
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxInt returns the maximum of two integers.
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
