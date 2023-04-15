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
	//Kludge alert. We need both the file's name and type for images
	//so for images we pass both in and stitch them back together here.

	if strings.Split(f_name, "_")[0] != "nsfw" {
		return_string = "↓ " + file_name_string + "\n" + fmt.Sprintln("!["+file_name_string+"]("+file_name_string+")\n")
		//put the filename for the image below, with an arrow.
		//the image index is most useful for writing your own index_file_name files.
	} else {
		//Sometimes there are images we don't want to auto-load. So instead,
		//we strip the filename of its underscores and create a standard anchor link to it.

		file_n_short :=
			//	strings.Split(
			strings.ReplaceAll(
				f_name,
				"_", " ") //,
			//	"/")[0]
		//Already, I do not recall why we're splitting on /.

		return_string = file_name_string + "→" + fmt.Sprintln("["+file_n_short+"]("+file_name_string+")\n")
	}

	return return_string

}
