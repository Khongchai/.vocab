package vocab_testing

import "strings"

func TrimLines(text string) string {
	split := strings.Split(text, "\n")
	filtered := []string{}
	for i := range split {
		trimmed := strings.TrimSpace(split[i])
		if trimmed == "" {
			continue
		}
		filtered = append(filtered, trimmed)
	}
	joined := strings.Join(filtered, "\n")
	return joined
}
