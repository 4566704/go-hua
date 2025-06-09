package common

import "math"

func IgnoreDigits(value, digitsToIgnore int) int {
	divider := int(math.Pow10(digitsToIgnore))
	return (value / divider) * divider
}
