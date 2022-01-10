package internal

import "strings"

func CleanUserInput(text string) (cleaned string) {
	cleaned = strings.Replace(text, "\n", "", -1)
	cleaned = strings.Replace(cleaned, "\r", "", -1)
	return
}
