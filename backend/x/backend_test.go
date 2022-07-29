package x

import (
	"testing"
)

func TestParseDisplay(t *testing.T) {
	var cases = []struct {
		Spec    string
		Host    string
		Display int
		Screen  int
	}{
		{":0", "", 0, 0},
		{"example.com:0.20", "example.com", 0, 20},
	}

	for _, c := range cases {
		host, display, screen, err := parseDisplay(c.Spec)
		if err != nil {
			t.Fatal(err)
		}
		if host != c.Host ||
			display != c.Display ||
			screen != c.Screen {
			t.Fatalf("wrong result host=%s display=%d screen=%d for %v", host, display, screen, c)
		}
	}
}
