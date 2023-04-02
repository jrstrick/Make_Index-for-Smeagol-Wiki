package main

import "testing"

func Test_format_dir_link(t *testing.T) {

	got := format_dir_link("Simple_test_link")
	want := "[Simple test link](Simple_test_link/Index.md)\n\n"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
