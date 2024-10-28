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
			svg: `crow: <img draggable="false" class="emoji" src="` + OfficialCDN + `svg/1f426-200d-2b1b.svg" width="72" height="72" alt="🐦‍⬛"/>`,
			png: `crow: <img draggable="false" class="emoji" src="` + OfficialCDN + `72x72/1f426-200d-2b1b.png" width="72" height="72" alt="🐦‍⬛"/>`,
		},
		{
			in:  "🌎, hello! for 🐦🦤",
			svg: `<img draggable="false" class="emoji" src="` + OfficialCDN + `svg/1f30e.svg" width="72" height="72" alt="🌎"/>, hello! for <img draggable="false" class="emoji" src="` + OfficialCDN + `svg/1f426.svg" width="72" height="72" alt="🐦"/><img draggable="false" class="emoji" src="` + OfficialCDN + `svg/1f9a4.svg" width="72" height="72" alt="🦤"/>`,
			png: `<img draggable="false" class="emoji" src="` + OfficialCDN + `72x72/1f30e.png" width="72" height="72" alt="🌎"/>, hello! for <img draggable="false" class="emoji" src="` + OfficialCDN + `72x72/1f426.png" width="72" height="72" alt="🐦"/><img draggable="false" class="emoji" src="` + OfficialCDN + `72x72/1f9a4.png" width="72" height="72" alt="🦤"/>`,
		},
		{
			in:  "6️⃣9️⃣ nice",
			svg: `<img draggable="false" class="emoji" src="` + OfficialCDN + `svg/36-20e3.svg" width="72" height="72" alt="6️⃣"/><img draggable="false" class="emoji" src="` + OfficialCDN + `svg/39-20e3.svg" width="72" height="72" alt="9️⃣"/> nice`,
			png: `<img draggable="false" class="emoji" src="` + OfficialCDN + `72x72/36-20e3.png" width="72" height="72" alt="6️⃣"/><img draggable="false" class="emoji" src="` + OfficialCDN + `72x72/39-20e3.png" width="72" height="72" alt="9️⃣"/> nice`,
		},
	}
	svg := New()
	png := New(WithFormat(PNG))
	for _, try := range table {
		t.Run(try.in, func(t *testing.T) {
			t.Logf("%#v", []byte(try.in))
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
	test := "hello 🐦‍⬛ world 🌎 for 🐦 & 🦤 & 5️⃣!"
	tw := New(WithAttrs(func(emoji string, defaults []html.Attribute) []html.Attribute {
		return append(defaults, html.Attribute{Key: "data-md", Val: emoji})
	}))
	got := tw.Replace(test)
	for _, emoji := range []string{"🐦‍⬛", "🌎", "🐦", "🦤", "5️⃣"} {
		mdattr := fmt.Sprintf(`data-md="%s"`, emoji)
		if !strings.Contains(got, mdattr) {
			t.Error("not found:", emoji, "in:", got)
		}
	}
}

func TestHTML(t *testing.T) {
	text := &html.Node{
		Type: html.TextNode,
		Data: "hello 🐦‍⬛ world 🌎 for 🐦 & 🦤! 5️⃣",
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
	if ct != 5 {
		t.Error("unexpected # of img elements:", ct)
	}
}

func BenchmarkTwemojiReplace(b *testing.B) {
	for n := 0; n < b.N; n++ {
		WriteString(io.Discard, "hello 🐦‍⬛ world 🌎 for 🐦 & 5️⃣!")
	}
}

func BenchmarkTwemojiHTML(b *testing.B) {
	for n := 0; n < b.N; n++ {
		text := &html.Node{
			Type: html.TextNode,
			Data: "hello 🐦‍⬛ world 🌎 for 🐦 & 5️⃣!",
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
