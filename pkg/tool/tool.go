package tool

import "net/url"

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
