package emojify

//go:generate ./gen.sh

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	officialCDN  = "https://cdn.jsdelivr.net/gh/jdecked/twemoji@15.1.0/assets/"
	defaultClass = "emoji"
)

type Twemoji struct {
	cdn   string
	class string
	fmt   Format
	attrs AttrFunc

	replacer *strings.Replacer
	nodes    map[rune][]emojiMatch
}

type emojiMatch struct {
	str  string
	img  string
	node *html.Node
}

func New(opts ...Option) Twemoji {
	t := Twemoji{
		cdn:   officialCDN,
		fmt:   FormatSVG,
		class: defaultClass,
		nodes: make(map[rune][]emojiMatch),
	}
	for _, opt := range opts {
		opt(&t)
	}
	if err := t.load(); err != nil {
		panic(fmt.Errorf("twemoji failed to load: %w", err))
	}
	return t
}

func (twemoji *Twemoji) load() error {
	repl := make([]string, 0, len(twemojiFiles)*2)
	var buf bytes.Buffer
	for _, base := range twemojiFiles {
		ext := path.Ext(base)
		filename := base[:len(base)-len(ext)]
		hexes := strings.Split(filename, "-")
		runes := make([]rune, len(hexes))
		for i, hex := range hexes {
			n, err := strconv.ParseInt(hex, 16, 64)
			if err != nil {
				return err
			}
			runes[i] = rune(n)
		}
		item := emojiMatch{
			str:  string(runes),
			img:  base,
			node: twemoji.node(string(runes), base),
		}

		buf.Reset()
		if err := html.Render(&buf, item.node); err != nil {
			return err
		}

		elem := buf.String()
		repl = append(repl, item.str, elem)

		head, _ := utf8.DecodeRuneInString(item.str)
		matches := twemoji.nodes[head]
		matches = append(matches, item)
		twemoji.nodes[head] = matches
	}
	twemoji.replacer = strings.NewReplacer(repl...)
	return nil
}

func (tw Twemoji) node(emoji string, src string) *html.Node {
	dir := tw.fmt.Dir()
	if tw.fmt == FormatPNG {
		src = src[:len(src)-len("svg")] + "png"
	}
	img := &html.Node{
		Type:     html.ElementNode,
		Data:     "img",
		DataAtom: atom.Img,
		Attr: []html.Attribute{
			{Key: "draggable", Val: "false"},
			{Key: "class", Val: tw.class},
			{Key: "src", Val: tw.cdn + dir + src},
			{Key: "height", Val: "72"},
			{Key: "alt", Val: emoji},
		},
	}
	if tw.attrs != nil {
		img.Attr = tw.attrs(emoji, img.Attr)
	}
	return img
}

// Option used in [New].
type Option func(*Twemoji)

func WithCDN(href string) Option {
	if href != "" && !strings.HasSuffix(href, "/") {
		href = href + "/"
	}
	return func(t *Twemoji) {
		t.cdn = href
	}
}

func WithClass(class string) Option {
	return func(t *Twemoji) {
		t.class = class
	}
}

// AttrFunc is a function for choosing attributes for Twemoji <img> tags.
type AttrFunc func(emoji string, defaults []html.Attribute) []html.Attribute

// WithAttrs specifies a custom HTML attribute function.
func WithAttrs(fn AttrFunc) Option {
	return func(t *Twemoji) {
		t.attrs = fn
	}
}

// WithFormat specifies the desired image format (default SVG).
func WithFormat(f Format) Option {
	return func(t *Twemoji) {
		t.fmt = f
	}
}

// Format of emoji replacement images.
type Format string

const (
	FormatSVG = "svg"
	FormatPNG = "png"
)

func (f Format) Dir() string {
	if f == FormatSVG {
		return "svg/"
	}
	return "72x72/"
}

// Replace returns a copy of s with all emojis replaced by <img> tags.
// Does NOT sanitize s. Use ReplaceHTML instead to safely replace HTML text.
func (tw Twemoji) Replace(s string) string {
	if tw.replacer == nil {
		return Default.Replace(s)
	}
	return tw.replacer.Replace(s)
}

// WriteString writes s to w with all emojis replaced by <img> tags.
// Does NOT sanitize s. Use ReplaceHTML instead to safely replace HTML text.
func (tw Twemoji) WriteString(w io.Writer, s string) (n int, err error) {
	if tw.replacer == nil {
		return Default.WriteString(w, s)
	}
	return tw.replacer.WriteString(w, s)
}
