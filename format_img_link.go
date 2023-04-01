package main

import (
	"fmt"
	"log"
	"strings"
)

func format_img_link(f_name string, f_type string) string {
	/*----------------------------------------------------------
	take an image name, return a markdown_link to the image in that
	directory as a string.
	------------------------------------------------------------*/

	log.Print("format_img_link. f_name is: ", f_name)

	var return_string string
	file_name_string := fmt.Sprintln(f_name + "." + f_type)

	if strings.Split(f_name, "_")[0] != "nsfw" {
		return_string = "↓ " + file_name_string + "\n" + fmt.Sprintln("!["+file_name_string+"]("+file_name_string+")\n")
	} else {
		file_n_short :=
			strings.Split(
				strings.ReplaceAll(
					f_name,
					"_", " "),
				"/")[0]
		return_string = file_name_string + "→" + fmt.Sprintln("["+file_n_short+"]("+file_name_string+")\n")
	}

	return return_string

}
