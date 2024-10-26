package emojify

import (
	"html/template"
	"io"

	"golang.org/x/net/html"
)

// Default configuration using official CDN and SVG images.
var Default = New()

// Replace returns a copy of s with emojis replaced by <img> tags.
// Does NOT sanitize s.
func Replace(s string) string {
	return Default.Replace(s)
}

// HTML escapes text and returns HTML having emojis replaced with twemoji images.
func HTML(text string) template.HTML {
	safe := html.EscapeString(text)
	return template.HTML(Default.Replace(safe))
}

// ReplaceHTML mutates the HTML of root, replacing emojis in text nodes with twemoji images.
func ReplaceHTML(root *html.Node) {
	replaceTextNodes(root, Default.replaceEmojis)
}

// WriteString writes s to w with all emojis replaced by <img> tags.
// Does NOT sanitize s. Use [ReplaceHTML] instead to safely replace HTML text.
func WriteString(w io.Writer, s string) (n int, err error) {
	return Default.WriteString(w, s)
}
