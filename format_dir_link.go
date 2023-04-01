package main

import (
	"fmt"
	"log"
	"strings"
)

func format_dir_link(f_name string) string {
	/*----------------------------------------------------------
	take a directory name, return a markdown_link to the Index.md file in that
	directory as a string.
	------------------------------------------------------------*/

	log.Print("format_dir_link. f_name is: ", f_name)

	file_n_short := strings.ReplaceAll(f_name, "_", " ")

	return fmt.Sprintln("[" + file_n_short + "](" + f_name + "/Index.md)\n")

}
