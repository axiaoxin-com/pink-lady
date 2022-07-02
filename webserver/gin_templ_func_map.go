// Functions from Go's strings package usable as template actions

package webserver

import (
	"html/template"
	"strings"

	"github.com/chai2010/gettext-go"
	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
)

// TemplFuncs is a template.FuncMap with functions that can be used as template actions.
var TemplFuncs = map[string]interface{}{
	"_text": gettext.Gettext,
	"_safe_url": func(s string) template.URL {
		return template.URL(s)
	},
	"_md2html": func(s string) template.HTML {
		md := []byte(s)
		html := markdown.ToHTML(md, nil, nil)
		html = bluemonday.UGCPolicy().SanitizeBytes(html)
		return template.HTML(html)
	},

	"_int_sum": func(ints ...int) int {
		sum := 0
		for _, i := range ints {
			sum += i
		}
		return sum
	},
	"_str_repeat": func(s string, count int) string {
		return strings.Repeat(s, count)
	},
}
