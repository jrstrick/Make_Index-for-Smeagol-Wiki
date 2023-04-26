package main

import (
	"fmt"
	"make_index/log_wrapper"
	"strings"
)

func format_mkd_link(f_name string) string {
	// take a filename, return a markdown_link as a string.

	log_wrapper.Print("format_mkd_link. f_name is: ", f_name)

	file_n_short :=
		//	strings.Split(
		strings.ReplaceAll(
			f_name,
			"_", " ") //,
	//		"/")[0]
	//strip the underscores from file_name_short.
	//Already I do not recall why we split on /.

	return fmt.Sprintln("[" + file_n_short + "](" + f_name + ".md)\n")

}
