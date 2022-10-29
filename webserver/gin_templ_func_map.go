// Functions from Go's strings package usable as template actions

package webserver

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"net/url"
	"strings"
	"time"

	mermaid "github.com/abhinav/goldmark-mermaid"
	chromav2html "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/chai2010/gettext-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	// SaltHashids salt for hashids
	SaltHashids = "pHZT6uIvVfxpreeumOhuuaH1m4EpR02el3wJbqcDtsgtIuRY13"
	// MinLenHashids hashid 长度
	MinLenHashids = 8
)

// TemplFuncs is a template.FuncMap with functions that can be used as template actions.
var TemplFuncs = map[string]interface{}{
	"_text": gettext.Gettext,
	"_safe_url": func(s string) template.URL {
		return template.URL(s)
	},
	"_md2html": func(s string, ugc bool) template.HTML {
		md := goldmark.New(
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
				html.WithXHTML(),
			),
			goldmark.WithExtensions(
				extension.GFM,
				highlighting.NewHighlighting(
					highlighting.WithStyle("manni"),
					highlighting.WithFormatOptions(
						chromav2html.WithLineNumbers(true),
					),
				),
				&mermaid.Extender{},
				extension.Footnote,
				extension.DefinitionList,
				extension.CJK,
				extension.Typographer,
			),
		)
		if !ugc {
			md.Renderer().AddOptions(html.WithUnsafe())
		}
		buf := &bytes.Buffer{}
		if err := md.Convert([]byte(s), buf); err != nil {
			logging.Error(nil, "_md2html Convert error:"+err.Error())
			return ""
		}
		return template.HTML(buf.String())
	},
	"_url_query_escape": func(s string) string {
		return url.QueryEscape(s)
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
	"_int_mod": func(i1, i2 int) int {
		return i1 % i2
	},
	"_percent": func(i1, i2 int) int {
		if i2 == 0 {
			return 0
		}
		p := int(float64(i1) / float64(i2) * 100)
		if p == 0 {
			p = 1
		}
		return p
	},
	"_str_repeat": func(s string, count int) string {
		return strings.Repeat(s, count)
	},
	"_str_slice_join": func(s []string, sep string) string {
		return strings.Join(s, sep)
	},
	"_int_hash_encode": func(i int) string {
		return goutils.IntHashEncode(i, SaltHashids, MinLenHashids, "")
	},
	"_int32_hash_encode": func(i int32) string {
		return goutils.IntHashEncode(int(i), SaltHashids, MinLenHashids, "")
	},
	"_last_page_offset": func(total int, limit int) int {
		offset := total / limit * limit
		if offset >= total {
			offset = total - limit
		}
		return offset
	},
	"_str_cut": func(s string, cut int) string {
		r := []rune(s)
		if len(r) < cut {
			return s
		}
		return string(r[:cut]) + "..."
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
	"_time_ago": func(ts int64) string {
		nowts := time.Now().Local().Unix()
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			logging.Error(nil, err.Error())
		} else {
			nowts = time.Now().In(loc).Unix()
		}
		duration := nowts - ts
		if duration < 0 {
			duration = 0
		}
		if duration <= 60 {
			return fmt.Sprintf("%d%s", duration, gettext.Gettext("秒钟前"))
		}

		minutes := int(math.Floor(float64(duration) / 60.0))
		return fmt.Sprintf("%d%s", minutes, gettext.Gettext("分钟前"))
	},
}
