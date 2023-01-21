package tool

import (
	"net/url"
	"strings"
)

func BuildURL(path string, params map[string][]string) (string, error) {
	paramMap := url.Values{}
	for k, v := range params {
		if len(v) > 0 {
			for _, value := range v {
				paramMap.Set(k, value)
			}
		}
	}
	base, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	base.RawQuery = paramMap.Encode()
	return base.String(), nil
}

func StringBuilder(args ...string) string {
	var build strings.Builder
	for _, value := range args {
		build.WriteString(value)
	}
	return build.String()
}
