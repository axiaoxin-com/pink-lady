// Functions from Go's strings package usable as template actions

package webserver

import (
	"html/template"
	"strings"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/chai2010/gettext-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
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
	"_int_sub": func(ints ...int) int {
		if len(ints) == 0 {
			return 0
		}
		sub := ints[0]
		for _, i := range ints[1:] {
			sub -= i
		}
		return sub
	},
	"_percent": func(i1, i2 int) int {
		if i2 == 0 {
			return 0
		}
		return int(float64(i1) / float64(i2) * 100)
	},
	"_str_repeat": func(s string, count int) string {
		return strings.Repeat(s, count)
	},
	"_last_page_offset": func(total int, limit int) int {
		return total / limit * limit
	},
	"_str_cut": func(s string, cut int) string {
		return string([]rune(s)[:cut])
	},
	"_pbts_format": func(pbts *timestamp.Timestamp, format string) string {
		gotime, err := ptypes.Timestamp(pbts)
		if err != nil {
			logging.Error(nil, err.Error())
			return ""
		}
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			logging.Error(nil, err.Error())
		} else {
			return gotime.In(loc).Format(format)
		}
		return gotime.Local().Format(format)
	},
	"_time_format": func(t time.Time, format string) string {
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			logging.Error(nil, err.Error())
		} else {
			return t.In(loc).Format(format)
		}
		return t.Local().Format(format)
	},
}
