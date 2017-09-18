package main

import (
	"regexp"
	"testing"
)

func TestVersion(t *testing.T) {
	pattern := regexp.MustCompile(`\d+\.\d+\.\d+`)
	if !pattern.MatchString(Version) {
		t.Errorf("Expected %#v to match \"%s\".", Version, pattern.String())
	}
}
