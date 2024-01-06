package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func SafeSub[T constraints.Integer](a, b T) T {
	if a > b {
		return a - b
	}
	return 0
}

func FormatSize(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
