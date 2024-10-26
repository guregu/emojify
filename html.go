package emojify

import (
	"html/template"
	"strings"
	"unicode"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// HTML escapes text and returns HTML having emojis replaced with twemoji images.
// Consider using with [html/template.Template.Funcs].
func (tw Twemoji) HTML(text string) template.HTML {
	if tw.replacer == nil {
		return Default.HTML(text)
	}
	safe := html.EscapeString(text)
	return template.HTML(tw.Replace(safe))
}

// ReplaceHTML mutates the HTML of root, replacing emojis in text nodes with twemoji images.
// Useful for replacing emoji in HTML you've already rendered (e.g. markdown rendering).
func (tw Twemoji) ReplaceHTML(root *html.Node) {
	if tw.replacer == nil {
		Default.ReplaceHTML(root)
		return
	}
	replaceTextNodes(root, tw.replaceEmojis)
}

func (tw Twemoji) replaceEmojis(node *html.Node) *html.Node {
	search := node.Data
	span := &html.Node{
		Type:     html.ElementNode,
		Data:     "span",
		DataAtom: atom.Span,
	}

	var consumed int
	emit := func(next *html.Node, idx int) {
		consumed = idx
		span.AppendChild(next)
	}

	var skip int
	var hit bool
scan:
	for idx, char := range search {
		if (char <= unicode.MaxASCII) ||
			(skip != 0 && skip > idx) {
			continue
		}
		match, ok := tw.nodes[char]
		if !ok {
			continue
		}
		for _, m := range match {
			if strings.HasPrefix(search[idx:], m.str) {
				hit = true
				if idx > 0 && consumed < idx {
					// regular text before the emoji
					emit(&html.Node{
						Type: html.TextNode,
						Data: search[consumed:idx],
					}, idx)
				}
				// actual emoji
				clone := *m.node
				skip = idx + len(m.str)
				emit(&clone, skip)
				continue scan
			}
		}
	}
	if !hit {
		return nil
	}
	// "leftovers"
	if consumed < len(search) {
		emit(&html.Node{
			Type: html.TextNode,
			Data: search[consumed:],
		}, 0)
	}
	return span
}

func replaceTextNodes(root *html.Node, do func(*html.Node) *html.Node) {
	switch {
	case root == nil:
		return
	case root.Type == html.TextNode:
		if rewrite := do(root); rewrite != nil {
			replaceChild(root.Parent, root, rewrite)
		}
		return
	}
	for node := root.FirstChild; node != nil; node = node.NextSibling {
		switch node.Type {
		case html.TextNode:
			if rewrite := do(node); rewrite != nil {
				replaceChild(root, node, rewrite)
				continue
			}
		default:
			if node.FirstChild != nil {
				replaceTextNodes(node, do)
			}
		}
	}
}

func replaceChild(parent, target, replace *html.Node) bool {
	if parent == nil || target == nil || replace == nil {
		return false
	}
	for n := parent.FirstChild; n != nil; n = n.NextSibling {
		if n == target {
			if p := n.Parent; p != nil {
				if p.FirstChild == n {
					p.FirstChild = replace
				}
				if p.LastChild == n {
					p.LastChild = replace
				}
			}
			if prev := n.PrevSibling; prev != nil {
				prev.NextSibling = replace
			}
			replace.PrevSibling = n.PrevSibling
			if next := n.NextSibling; next != nil {
				next.PrevSibling = replace
			}
			replace.NextSibling = n.NextSibling
			return true
		}
	}
	return false
}
