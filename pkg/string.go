package pkg

import (
	"regexp"
	"strings"
)

func LineToLowCamel(str string) string {
	strSlice := strings.Split(str, "_")
	for k, s := range strSlice {
		if k == 0 {
			continue
		}
		strSlice[k] = strings.ToUpper(s[:1]) + s[1:]
	}
	return strings.Join(strSlice, "")
}

func LineToUpCamel(str string) string {
	if str == "" {
		return ""
	}
	strSlice := strings.Split(str, "_")
	for k, s := range strSlice {
		strSlice[k] = strings.ToUpper(s[:1]) + s[1:]
	}
	return strings.Join(strSlice, "")
}

// Inline 备注变成一行
func Inline(str string) string {
	str = strings.Replace(str, "\n", " ", -1)
	str = strings.Replace(str, "\t", " ", -1)
	return str
}

// ContainsNumber 判断字符串是否包含数字
func ContainsNumber(s string) bool {
	pattern := regexp.MustCompile(`\d`)
	return pattern.MatchString(s)
}

func FormatDescription(desc string) (string, string) {
	pattern := `binding:"([^"]+)" `
	// 编译正则表达式
	reg := regexp.MustCompile(pattern)
	if match := reg.FindStringSubmatch(desc); len(match) == 2 {
		return strings.Trim(match[1], " "), strings.Replace(desc, match[0], "", -1)
	}
	return "", desc
}

func GetRequestName(path string) string {
	return GetFuncName(path) + "Req"
}

func GetFuncName(path string) string {
	paths := strings.Split(path, "/")
	name := paths[len(paths)-1]
	name = LineToUpCamel(name)
	return name
}

func GetResponseName(path string) string {
	return GetFuncName(path) + "Resp"
}

func GetPackageName(path string) string {
	paths := strings.Split(path, "/")
	if len(paths) <= 1 {
		return ""
	}
	return paths[len(paths)-2]
}

func GetTopLevelName(path string) string {
	path = strings.Trim(path, "/")
	paths := strings.Split(path, "/")
	if len(paths) <= 1 {
		return ""
	}
	return paths[0]
}
