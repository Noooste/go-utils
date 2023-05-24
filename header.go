package utils

import (
	http "github.com/Noooste/fhttp"
	"net/url"
)

type HeaderType uint8

const (
	DefaultHeaders HeaderType = iota
	scriptGetHeaders
	scriptPostHeaders
)

func GetHeaderAndOrder(t HeaderType, referer *url.URL, additional ...[2]string) (header http.Header, order []string) {
	switch t {
	case DefaultHeaders:
		header = http.Header{
			"sec-ch-ua":          {SecChUa},
			"sec-ch-ua-mobile":   {"?0"},
			"user-agent":         {UserAgent},
			"sec-ch-ua-platform": {"\"Windows\""},
			"accept":             {"*/*"},
			"sec-fetch-site":     {"none"},
			"sec-fetch-mode":     {"navigate"},
			"sec-fetch-user":     {"?1"},
			"sec-fetch-dest":     {"document"},
			"accept-encoding":    {"gzip, deflate, br"},
			"accept-language":    {"en-US;q=0.8,en;q=0.7"},
		}

		order = []string{
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"user-agent",
			"sec-ch-ua-platform",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"sec-fetch-dest",
			"accept-encoding",
			"accept-language",
		}

	case scriptGetHeaders:
		header = http.Header{
			"sec-ch-ua":          {SecChUa},
			"sec-ch-ua-mobile":   {"?0"},
			"user-agent":         {UserAgent},
			"sec-ch-ua-platform": {"\"Windows\""},
			"accept":             {"*/*"},
			"sec-fetch-site":     {"same-origin"},
			"sec-fetch-mode":     {"no-cors"},
			"sec-fetch-dest":     {"script"},
			"referer":            {referer.String()},
			"accept-encoding":    {"gzip, deflate, br"},
			"accept-language":    {"en-US;q=0.8,en;q=0.7"},
		}

		order = []string{
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"user-agent",
			"sec-ch-ua-platform",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
		}

	case scriptPostHeaders:
		header = http.Header{
			"sec-ch-ua":          {SecChUa},
			"content-type":       {"text/plain;charset=UTF-8"},
			"sec-ch-ua-mobile":   {"?0"},
			"user-agent":         {UserAgent},
			"sec-ch-ua-platform": {"\"Windows\""},
			"accept":             {"*/*"},
			"origin":             {referer.Scheme + "://" + referer.Host},
			"sec-fetch-site":     {"same-origin"},
			"sec-fetch-mode":     {"cors"},
			"sec-fetch-dest":     {"empty"},
			"referer":            {referer.String()},
			"accept-encoding":    {"gzip, deflate, br"},
			"accept-language":    {"en-US;q=0.8,en;q=0.7"},
		}

		order = []string{
			"content-length",
			"sec-ch-ua",
			"content-type",
			"sec-ch-ua-mobile",
			"user-agent",
			"sec-ch-ua-platform",
			"accept",
			"origin",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
		}

	default:
		header = http.Header{}
		order = []string{}
	}

	for _, v := range additional {
		header[v[0]] = []string{v[1]}
		order = append(order, v[0])
	}

	return header, order
}
