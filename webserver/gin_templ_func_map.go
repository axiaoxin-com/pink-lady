// Functions from Go's strings package usable as template actions

package webserver

import "strings"

// TemplFuncs is a template.FuncMap with functions that can be used as template actions.
var TemplFuncs = map[string]interface{}{
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
