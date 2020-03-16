package main

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

type tcase struct {
	name          string
	regexMethod   *regexp.Regexp
	regexPath     *regexp.Regexp
	expectedValue bool
}

func TestMatchOnlyForMethod(t *testing.T) {
	cases := []tcase{
		{name: "literal match", regexMethod: regexp.MustCompile("DELETE"), expectedValue: true},
		{name: "two posibilities", regexMethod: regexp.MustCompile("DELETE|OPTIONS"), expectedValue: true},
		{name: "not matching", regexMethod: regexp.MustCompile("GET"), expectedValue: false},
	}
	request, _ := http.NewRequest("DELETE", "http://mytest", nil)
	for _, c := range cases {
		matcher := makeLogRequestMatcher(c.regexMethod, nil)
		assert.Equal(t, c.expectedValue, matcher(request), "failed for case %s", c.name)
	}
}

func TestMatchOnlyForPath(t *testing.T) {
	cases := []tcase{
		{name: "literal match", regexPath: regexp.MustCompile("^/hello"), expectedValue: true},
		{name: "two posibilities", regexPath: regexp.MustCompile("^/(hello|/bye)"), expectedValue: true},
		{name: "prefix match", regexPath: regexp.MustCompile("^/hello*"), expectedValue: true},
		{name: "prefix not matching", regexPath: regexp.MustCompile("^/hello/"), expectedValue: false},
	}
	errorMsgs := map[bool]string{
		true:  "failed test for case %q: %q did not match regex %q",
		false: "failed test for case %q: %q matched regex %q",
	}

	request, _ := http.NewRequest("DELETE", "http://mytest/hello?something", nil)
	for _, c := range cases {
		matcher := makeLogRequestMatcher(nil, c.regexPath)
		assert.Equal(t, c.expectedValue, matcher(request), errorMsgs[c.expectedValue], c.name, request.URL.Path, c.regexPath)
	}
}

func TestMatchBothMethodAndPath(t *testing.T) {
	cases := []tcase{
		{name: "literal match", regexPath: regexp.MustCompile("/hello"), regexMethod: regexp.MustCompile("DELETE"), expectedValue: true},
		{name: "method does not match", regexPath: regexp.MustCompile("/hello"), regexMethod: regexp.MustCompile("GET"), expectedValue: false},
		{name: "path does not match", regexPath: regexp.MustCompile("/bye"), regexMethod: regexp.MustCompile("DELETE"), expectedValue: false},
	}

	request, _ := http.NewRequest("DELETE", "http://mytest/hello?something", nil)
	for _, c := range cases {
		matcher := makeLogRequestMatcher(c.regexMethod, c.regexPath)
		assert.Equal(t, c.expectedValue, matcher(request), "failed test for case %q", c.name)
	}
}
