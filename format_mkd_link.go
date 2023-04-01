package main

import (
	"fmt"
	"log"
	"strings"
)

func format_mkd_link(f_name string) string {
	// take a filename, return a markdown_link as a string.

	log.Print("f_name is: ", f_name, "f_type is md")

	file_n_short :=
		strings.Split(
			strings.ReplaceAll(
				f_name,
				"_", " "),
			"/")[0]

	return fmt.Sprintln("[" + file_n_short + "](" + f_name + ".md)\n")

}
