package emojify

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestTwemoji(t *testing.T) {
	table := []struct {
		in  string
		svg string
		png string
	}{
		{
			in:  "crow: 🐦‍⬛",
			svg: `crow: <img draggable="false" class="emoji" src="https://cdn.jsdelivr.net/gh/jdecked/twemoji@15.1.0/assets/svg/1f426-200d-2b1b.svg" height="72" alt="🐦‍⬛"/>`,
			png: `crow: <img draggable="false" class="emoji" src="https://cdn.jsdelivr.net/gh/jdecked/twemoji@15.1.0/assets/72x72/1f426-200d-2b1b.png" height="72" alt="🐦‍⬛"/>`,
		},
		{
			in:  "🌎, hello! for 🐦🦤",
			svg: `<img draggable="false" class="emoji" src="https://cdn.jsdelivr.net/gh/jdecked/twemoji@15.1.0/assets/svg/1f30e.svg" height="72" alt="🌎"/>, hello! for <img draggable="false" class="emoji" src="https://cdn.jsdelivr.net/gh/jdecked/twemoji@15.1.0/assets/svg/1f426.svg" height="72" alt="🐦"/><img draggable="false" class="emoji" src="https://cdn.jsdelivr.net/gh/jdecked/twemoji@15.1.0/assets/svg/1f9a4.svg" height="72" alt="🦤"/>`,
			png: `<img draggable="false" class="emoji" src="https://cdn.jsdelivr.net/gh/jdecked/twemoji@15.1.0/assets/72x72/1f30e.png" height="72" alt="🌎"/>, hello! for <img draggable="false" class="emoji" src="https://cdn.jsdelivr.net/gh/jdecked/twemoji@15.1.0/assets/72x72/1f426.png" height="72" alt="🐦"/><img draggable="false" class="emoji" src="https://cdn.jsdelivr.net/gh/jdecked/twemoji@15.1.0/assets/72x72/1f9a4.png" height="72" alt="🦤"/>`,
		},
	}
	svg := New()
	png := New(WithFormat(FormatPNG))
	for _, try := range table {
		t.Run(try.in, func(t *testing.T) {
			if got := svg.Replace(try.in); try.svg != got {
				t.Errorf("svg(%q) →\n got: %q\nwant: %q", try.in, got, try.svg)
			}
			if got := png.Replace(try.in); try.png != got {
				t.Errorf("png(%q) →\n got: %q\nwant: %q", try.in, got, try.png)
			}
		})
	}
}

func TestWithAttr(t *testing.T) {
	test := "hello 🐦‍⬛ world 🌎 for 🐦 & 🦤"
	repl := New(WithAttrs(func(emoji string, defaults []html.Attribute) []html.Attribute {
		return append(defaults, html.Attribute{Key: "data-md", Val: emoji})
	}))
	got := repl.Replace(test)
	for _, emoji := range []string{"🐦‍⬛", "🌎", "🐦", "🦤"} {
		mdattr := fmt.Sprintf(`data-md="%s"`, emoji)
		t.Log(got)
		if !strings.Contains(got, mdattr) {
			t.Error("not found:", emoji, "in:", got)
		}
	}
}

func TestHTML(t *testing.T) {
	text := &html.Node{
		Type: html.TextNode,
		Data: "hello 🐦‍⬛ world 🌎 for 🐦 & 🦤!",
	}
	doc := &html.Node{
		Type:       html.ElementNode,
		Data:       "p",
		DataAtom:   atom.P,
		FirstChild: text,
		LastChild:  text,
	}
	text.Parent = doc
	ReplaceHTML(doc)

	span := doc.FirstChild
	if span == nil || span.Data != "span" {
		t.Fatal("not a span")
	}

	var ct int
	for node := span.FirstChild; node != nil; node = node.NextSibling {
		if node.Data == "img" {
			ct++
		}
	}
	if ct != 4 {
		t.Error("unexpected # of img elements:", ct)
	}
}

func BenchmarkTwemojiReplace(b *testing.B) {
	for n := 0; n < b.N; n++ {
		WriteString(io.Discard, "hello 🐦‍⬛ world 🌎 for 🐦 & 🦤!")
	}
}

func BenchmarkTwemojiHTML(b *testing.B) {
	for n := 0; n < b.N; n++ {
		text := &html.Node{
			Type: html.TextNode,
			Data: "hello 🐦‍⬛ world 🌎 for 🐦 & 🦤!",
		}
		doc := &html.Node{
			Type:       html.ElementNode,
			Data:       "p",
			DataAtom:   atom.P,
			FirstChild: text,
			LastChild:  text,
		}
		text.Parent = doc
		ReplaceHTML(doc)
	}
}
