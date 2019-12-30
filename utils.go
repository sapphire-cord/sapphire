package sapphire

import (
	"regexp"
)

var escapeReg = regexp.MustCompile("@(everyone|here)")

// Utilities to help in bot creation.

// Escape escapes @everyone/@here mentions by adding an invisible character to avoid the ping.
func Escape(input string) string {
	return escapeReg.ReplaceAllString(input, "@\u200b$1")
}
