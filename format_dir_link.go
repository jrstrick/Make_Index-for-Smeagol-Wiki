package main

import (
	"fmt"
	"make_index/log_wrapper"
	"strings"
)

func format_dir_link(f_name string, index_file_name string) string {
	/*----------------------------------------------------------
	take a directory name, return a markdown_link to the index_file_name
	file in that directory as a string.
	------------------------------------------------------------*/

	log_wrapper.Print("format_dir_link. f_name is: ", f_name)

	file_n_short := strings.ReplaceAll(f_name, "_", " ")
	//file_n_short is the human-readable filename in the
	//head of the link, so we strip out the underscores.

	return fmt.Sprintln("[" + file_n_short + "](" + f_name + "/" + index_file_name + ")\n")

}
