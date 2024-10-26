# emojify [![GoDoc](https://godoc.org/github.com/guregu/emojify?status.svg)](https://godoc.org/github.com/guregu/emojify)

Server-side rendering helpers for [Twemoji](https://github.com/jdecked/twemoji).

### Motivation

Many operating systems tie their emoji updates to major editions (e.g. Windows 11), leaving some users unable to display newer emoji.
Twemoji replaces emoji text with SVG or PNG images, but the official JS library does this on the client, leading to undesirable pop-in or hacks to avoid showing native emojis.
This library helps you render them server-side instead.

## Development

To update Twemoji and regenerate `twemoji.go`:

```bash
git submodule update --init --recursive
go generate
```
