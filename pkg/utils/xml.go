package utils

import "html"

func EscapeXML(input string) string {
	return html.EscapeString(input)
}
