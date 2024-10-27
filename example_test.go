package emojify_test

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/guregu/emojify"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ExampleHTML() {
	twemoji := emojify.New(emojify.WithCDN("https://twemoji.example.com/assets/"))

	tmpl := template.New("")
	// Install our emojifying function as "emojify"
	tmpl.Funcs(template.FuncMap{
		"emojify": twemoji.HTML,
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
	// Output: <article><h1>hello <img draggable="false" class="emoji" src="https://twemoji.example.com/assets/svg/1f30e.svg" width="72" height="72" alt="🌎"/></h1><p>no <img draggable="false" class="emoji" src="https://twemoji.example.com/assets/svg/1f6ab.svg" width="72" height="72" alt="🚫"/> javascript for me <img draggable="false" class="emoji" src="https://twemoji.example.com/assets/svg/1f606.svg" width="72" height="72" alt="😆"/></p></article>
}

func ExampleReplaceHTML() {
	twemoji := emojify.New(emojify.WithCDN("https://twemoji.example.com/assets/"))

	// first, we render markdown to html
	rendered := someMarkdownLibrary([]byte("# hello 🌎\n\nno 🚫 javascript for me 😆"), nil, nil)
	// then, parse the html
	body := &html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body}
	frags, err := html.ParseFragment(bytes.NewReader(rendered), body)
	if err != nil {
		panic(err)
	}
	// replace text elements with emoji in them
	var buf bytes.Buffer
	for _, frag := range frags {
		twemoji.ReplaceHTML(frag)
		if err := html.Render(&buf, frag); err != nil {
			panic(err)
		}
	}
	// our emojified HTML:
	fmt.Println(buf.String())
	// Output: <h1><span>hello <img draggable="false" class="emoji" src="https://twemoji.example.com/assets/svg/1f30e.svg" width="72" height="72" alt="🌎"/></span></h1><p><span>no <img draggable="false" class="emoji" src="https://twemoji.example.com/assets/svg/1f6ab.svg" width="72" height="72" alt="🚫"/> javascript for me <img draggable="false" class="emoji" src="https://twemoji.example.com/assets/svg/1f606.svg" width="72" height="72" alt="😆"/></span></p>
}

func someMarkdownLibrary(_ []byte, _, _ any) []byte {
	return []byte(`<h1>hello 🌎</h1><p>no 🚫 javascript for me 😆</p>`)
}
