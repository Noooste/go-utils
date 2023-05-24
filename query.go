package utils

import (
	"net/url"
	"strings"
)

type OrderedQuery struct {
	query []string
}

func (oq *OrderedQuery) Add(key string, value string) {
	oq.query = append(oq.query, key+"="+url.QueryEscape(value))
}

func (oq *OrderedQuery) String() string {
	return strings.Join(oq.query, "&")
}
