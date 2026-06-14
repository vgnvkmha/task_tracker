package render

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/net/html"
)

func RenderMarkdown(input string) template.HTML {
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Linkify,
		),
	)
	if err := md.Convert([]byte(input), &buf); err != nil {
		return template.HTML("")
	}

	policy := bluemonday.UGCPolicy()
	policy.RequireNoFollowOnLinks(true)
	policy.RequireNoReferrerOnLinks(true)
	policy.AddTargetBlankToFullyQualifiedLinks(true)

	safeHTML := policy.Sanitize(buf.String())
	return template.HTML(ensureLinkAttributes(safeHTML))
}

func ensureLinkAttributes(safeHTML string) string {
	root := &html.Node{Type: html.ElementNode, Data: "div"}
	nodes, err := html.ParseFragment(strings.NewReader(safeHTML), root)
	if err != nil {
		return safeHTML
	}
	for _, node := range nodes {
		root.AppendChild(node)
	}

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			setAttr(node, "target", "_blank")
			setAttr(node, "rel", mergeRel(getAttr(node, "rel")))
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(root)

	var out bytes.Buffer
	for child := root.FirstChild; child != nil; child = child.NextSibling {
		if err := html.Render(&out, child); err != nil {
			return safeHTML
		}
	}
	return out.String()
}

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func setAttr(node *html.Node, key string, value string) {
	for i, attr := range node.Attr {
		if attr.Key == key {
			node.Attr[i].Val = value
			return
		}
	}
	node.Attr = append(node.Attr, html.Attribute{Key: key, Val: value})
}

func mergeRel(value string) string {
	needed := []string{"nofollow", "noopener", "noreferrer"}
	seen := map[string]bool{}
	for _, item := range strings.Fields(value) {
		seen[item] = true
	}
	for _, item := range needed {
		if !seen[item] {
			value = strings.TrimSpace(value + " " + item)
		}
	}
	return value
}
