package utils

import (
	"testing"
)

func TestQuery(t *testing.T) {
	q := OrderedQuery{}
	q.Add("a", "b")
	q.Add("c", "d")
	if q.String() != "a=b&c=d" {
		t.Error("OrderedQuery.Add() failed")
	}
}

type q struct {
	A int    `json:"a"`
	B string `json:"b"`
	C string `json:"c,omitempty"`
}

func TestQuery2(t *testing.T) {
	dumped := UrlEncode(q{
		A: 1,
		B: "b",
	})
	if dumped != "a=a&b=b" {
		t.Error("UrlEncode() failed, expected a=1&b=b, got", dumped)
	}
}

func TestQuery3(t *testing.T) {
	dumped := UrlEncode(map[string]any{
		"a": "a",
		"b": "b",
	})
	if dumped != "a=a&b=b" {
		t.Error("UrlEncode() failed, expected a=a&b=b, got", dumped)
	}
}
