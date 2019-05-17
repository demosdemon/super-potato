package server

import (
	"fmt"
	"github.com/russross/blackfriday/v2"
	"io"
	"io/ioutil"
	"net/http"
)

const prettyHTMLTemplate = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Result</title>
	<style>
		html {
			font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
			font-weight: 300;
			background-color: #eee;
		}

		h1 {
			font-weight: 100;
		}

		body {
			margin: 3em;
		}

		.logo {
			display: block;
			margin: 10px auto;
			width: 100px;
			height: 100px;
		}
	</style>
	<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.6/styles/default.min.css">
</head>
<body>
	<img class="logo" src="/logo.svg?background=%%23ccc&foreground=%%23eee" alt="">
	<div class="markdown">%s</div>
	<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.6/highlight.min.js"></script>
	<script>hljs.initHighlighting();</script>
</body>
</html>
`

type Markdown struct {
	input io.Reader
}

func NewMarkdown(r io.Reader) *Markdown {
	return &Markdown{
		input: r,
	}
}

const extensions = blackfriday.CommonExtensions | blackfriday.AutoHeadingIDs | blackfriday.HardLineBreak

func (m *Markdown) Render(w http.ResponseWriter) error {
	m.WriteContentType(w)
	md, err := ioutil.ReadAll(m.input)
	if err != nil {
		return err
	}

	html := blackfriday.Run(md, blackfriday.WithExtensions(extensions))
	_, err = fmt.Fprintf(w, prettyHTMLTemplate, string(html))
	return err
}

func (m *Markdown) WriteContentType(w http.ResponseWriter) {
	if v, ok := w.Header()["Content-Type"]; ok && len(v) > 0 {
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
