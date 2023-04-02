package main

import "testing"

func Test_format_mkd_link(t *testing.T) {

	got := format_mkd_link("Simple_test_link")
	want := "[Simple test link](Simple_test_link.md)\n\n"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
