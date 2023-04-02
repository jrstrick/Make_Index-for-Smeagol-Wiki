package main

import "testing"

func Test_simple_format_img_link(t *testing.T) {

	got := format_img_link("Simple_test", "jpg")
	want := "↓ Simple_test.jpg\n\n![Simple_test.jpg\n](Simple_test.jpg\n)\n\n"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func Test_nsfw_format_img_link(t *testing.T) {
	got := format_img_link("nsfw_test", "jpg")
	want := "nsfw_test.jpg\n→[nsfw test](nsfw_test.jpg\n)\n\n"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}
