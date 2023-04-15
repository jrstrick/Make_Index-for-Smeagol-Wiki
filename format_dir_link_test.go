package main

import "testing"

func Test_format_dir_link(t *testing.T) {
	index_file_name := "Test_File_Name.md"
	got := format_dir_link("Simple_test_link", index_file_name)
	want := "[Simple test link](Simple_test_link" + "/" + index_file_name + ")\n\n"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
