package util

import (
	"fmt"
	"regexp"
)

// Example 1 - input:  "'driver'    => 'mysql',",
//                     "driver",
//                     false
//             output: "mysql"
// Example 2 - input:  "'host'      => env('DB_HOST', 'localhost'),",
//                     "host",
//                     true
//             output: "DB_HOST"
func KeyValueParser1(str, key string, isEnv bool) string {
	var rgxStr string
	if !isEnv {
		rgxStr = fmt.Sprintf(`'%s' *?=> *?\'(.*)\'`, key)
	} else {
		rgxStr = fmt.Sprintf(`'%s' *?=> *?env\( *?\'(.*)\' *?,`, key)
	}
	var rgx = regexp.MustCompile(rgxStr)
	rs := rgx.FindAllStringSubmatch(str, 1)
	if len(rs) > 0 && len(rs[0]) > 1 {
		return rs[0][1]
	}
	return ""
}
