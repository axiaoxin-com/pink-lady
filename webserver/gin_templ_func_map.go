// Functions from Go's strings package usable as template actions

package webserver

import (
	"bytes"
	"errors"
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
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/tidwall/gjson"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	// SaltHashids salt for hashids
	SaltHashids = "the-salt-hash-ids"
	// MinLenHashids hashid 长度
	MinLenHashids = 8
)

// TemplFuncs is a template.FuncMap with functions that can be used as template actions.
var TemplFuncs = map[string]interface{}{
	"_i18n": LangI18n,
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
	"_url_path_escape": func(s string) string {
		return url.PathEscape(s)
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
		if len(r) <= cut {
			return s
		}
		return string(r[:cut]) + "..."
	},
	"_md_clear": goutils.RemoveMarkdownTags,
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
	"_time_ago": func(ts int64, lang string) string {
		nowts := time.Now().Local().Unix()
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			logging.Error(nil, err.Error())
		} else {
			nowts = time.Now().In(loc).Unix()
		}
		seconds := nowts - ts
		if seconds < 0 {
			seconds = 0
		}
		if seconds <= 60 {
			return fmt.Sprintf("%d%s", seconds, LangI18n(lang, "秒钟前"))
		}

		minutes := int(math.Floor(float64(seconds) / 60.0))
		if minutes <= 60 {
			return fmt.Sprintf("%d%s", minutes, LangI18n(lang, "分钟前"))
		}

		hours := int(math.Floor(float64(minutes) / 60.0))
		if hours <= 24 {
			return fmt.Sprintf("%d%s", hours, LangI18n(lang, "小时前"))
		}
		days := int(math.Floor(float64(hours) / 24.0))
		if days <= 7 {
			return fmt.Sprintf("%d%s", days, LangI18n(lang, "天前"))
		} else if days <= 31 {
			return time.Unix(ts, 0).Local().Format("01-02 15:04")
		}

		return time.Unix(ts, 0).Local().Format("2006-01-02 15:04:05")
	},
	"_gjson_get": func(js string, key string) interface{} {
		return gjson.Get(js, key).Value()
	},
	"_dict": func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, errors.New("invalid dict call")
		}
		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, errors.New("dict keys must be strings")
			}
			dict[key] = values[i+1]
		}
		return dict, nil
	},
}
