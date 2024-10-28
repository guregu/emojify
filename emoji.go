package emojify

//go:generate ./script/gen.sh

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	// OfficialCDN is the official CDN hosting Twemoji assets.
	OfficialCDN  = "https://cdn.jsdelivr.net/gh/jdecked/twemoji@" + Version + "/assets/"
	defaultClass = "emoji"
)

// Twemoji is a configuration/cache of emoji replacements.
// The zero value will use [Default].
type Twemoji struct {
	cdn   string
	class string
	fmt   Format
	attrs AttrFunc

	replacer *strings.Replacer
	nodes    map[rune][]resource
}

type resource struct {
	str  string     // unicode text
	img  string     // filename
	node *html.Node // <img> element
}

// New creates a new [Twemoji] with the given set of [Option].
func New(opts ...Option) Twemoji {
	t := Twemoji{
		cdn:   OfficialCDN,
		fmt:   SVG,
		class: defaultClass,
		nodes: make(map[rune][]resource),
	}
	for _, opt := range opts {
		opt(&t)
	}
	if err := t.load(); err != nil {
		panic(fmt.Errorf("twemoji failed to load: %w", err))
	}
	return t
}

func (tw *Twemoji) load() error {
	keyvals := make([]string, 0, len(twemojiData)*2)
	var buf bytes.Buffer
	for _, item := range twemojiData {
		item.node = tw.node(item.str, item.img)

		buf.Reset()
		if err := html.Render(&buf, item.node); err != nil {
			return err
		}

		elem := buf.String()
		keyvals = append(keyvals, item.str, elem)

		head, _ := utf8.DecodeRuneInString(item.str)
		matches := tw.nodes[head]
		matches = append(matches, item)
		tw.nodes[head] = matches
	}
	tw.replacer = strings.NewReplacer(keyvals...)
	return nil
}

func (tw Twemoji) node(emoji string, src string) *html.Node {
	dir := tw.fmt.dir()
	if tw.fmt == PNG {
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
			{Key: "width", Val: "72"},
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

// WithCDN specifies the CDN (i.e. URL root) for the emoji image assets.
// Default value is the official (jsDelivr) CDN.
func WithCDN(href string) Option {
	if href != "" && !strings.HasSuffix(href, "/") {
		href = href + "/"
	}
	return func(t *Twemoji) {
		t.cdn = href
	}
}

// WithClass specifies the class given to emoji replacement <img> elements.
// Default is "emoji".
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
	// SVG images.
	SVG Format = "svg"
	// PNG (72x72 px) images.
	PNG Format = "png"
)

func (f Format) dir() string {
	switch f {
	case SVG:
		return "svg/"
	case PNG:
		return "72x72/"
	}
	return "/"
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
