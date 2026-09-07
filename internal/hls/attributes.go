package hls

import (
	"iter"
	"strings"
)

func parseAttributes(line string) iter.Seq2[string, string] {
	return func(yield func(k, v string) bool) {
		var start int
		inQuotes := false
		for i := 0; i < len(line); i++ {
			switch line[i] {
			case '"':
				inQuotes = !inQuotes
			case ',':
				if !inQuotes {
					if k, v, ok := strings.Cut(line[start:i], "="); ok {
						if !yield(strings.TrimSpace(k), strings.Trim(strings.TrimSpace(v), `"`)) {
							return
						}
						start = i + 1
					}
				}
			}
		}

		if start < len(line) {
			if k, v, ok := strings.Cut(line[start:], "="); ok {
				yield(strings.TrimSpace(k), strings.Trim(strings.TrimSpace(v), `"`))
			}
		}
	}
}
