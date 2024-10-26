package emojify_test

import (
	"html/template"
	"os"

	"github.com/guregu/emojify"
)

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
