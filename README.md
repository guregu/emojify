# emojify [![GoDoc](https://godoc.org/github.com/guregu/emojify?status.svg)](https://godoc.org/github.com/guregu/emojify)

Server-side rendering helpers for [Twemoji](https://github.com/jdecked/twemoji).

### Motivation

Many operating systems tie their emoji updates to major editions (e.g. Windows 11), leaving some users unable to display newer emoji.
Twemoji replaces emoji text with SVG or PNG images, but the official JS library does this on the client, leading to undesirable pop-in or hacks to avoid showing native emojis.
This library helps you render them server-side instead.

## Usage

### Configuring

You can change the URL and such, useful if you're self hosting.

```go
var Twemoji = emojify.New(
	emojify.WithCDN("https://selfhosted.example.com/static/twemoji/"),
	emojify.WithClass("twemoji"),
	emojify.WithFormat(emojify.SVG),
	emojify.WithAttrs(func(emoji string, defaults []html.Attribute) []html.Attribute {
		return append(defaults, html.Attribute{Key: "data-md", Val: emoji})
	}),
)
```

### `html/template`

You can use this library as a handy template function.

```go
func ExampleHTML() {
	tmpl := template.New("")
	// Install our emojifying function as "emojify"
	tmpl.Funcs(template.FuncMap{
		"emojify": emojify.HTML,
	})

	// Simple example template where we'll emojify the title and body text of an article
	tmpl = template.Must(tmpl.Parse(
		`<article><h1>{{.Title | emojify}}</h1><p>{{.Msg | emojify}}</p></article>`,
	))

	data := struct {
		Title, Msg string
	}{
		Title: "hello 🌎",
		Msg:   "no 🚫 javascript for me 😆",
	}
	if err := tmpl.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
	// Output: <article><h1>hello <img draggable="false" class="emoji" src="https://cdn.jsdelivr.net/gh/jdecked/twemoji@15.1.0/assets/svg/1f30e.svg" width="72" height="72" alt="🌎"/></h1><p>no <img draggable="false" class="emoji" src="https://cdn.jsdelivr.net/gh/jdecked/twemoji@15.1.0/assets/svg/1f6ab.svg" width="72" height="72" alt="🚫"/> javascript for me <img draggable="false" class="emoji" src="https://cdn.jsdelivr.net/gh/jdecked/twemoji@15.1.0/assets/svg/1f606.svg" width="72" height="72" alt="😆"/></p></article>
}
```

### Mutating HTML

Safely modify HTML by parsing it and replacing relevant text elements.
In this example we use Markdown rendering output.

```go
import (
	"bytes"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/gomarkdown/markdown"
	"github.com/guregu/emojify"
)

func renderMarkdown(text []byte) ([]byte, error) {
	// first, we render markdown to html
	rendered := markdown.ToHTML(text, nil, nil)
	// then, parse the html
	body := &html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body}
	frags, err := html.ParseFragment(bytes.NewReader(rendered), body)
	if err != nil {
		return nil, err
	}
	// replace text elements with emoji in them
	var buf bytes.Buffer
	for _, frag := range frags {
		emojify.ReplaceHTML(frag)
		if err := html.Render(&buf, frag); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
```

## Development

To update Twemoji and regenerate `twemoji.go`:

```bash
git submodule update --init --recursive
go generate
```
