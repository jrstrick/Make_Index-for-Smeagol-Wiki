package main

import (
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode"

	"make_index/log_wrapper"

	"github.com/adrg/frontmatter"
)

func on_error_die(err error, text string) {
	/*-----------------------------------------------------------------
	Just a quick function that checks the value of a standard error
	type and if it's not nil, logs a reason string and panics.
	-----------------------------------------------------------------*/

	if err != nil {
		//if we have an error...

		log_wrapper.Fatal("PANIC: "+text+"\n\t Reason: ", err)
		//call log_wrapper.Fatal with the text from the function call
		//and the error itself. The log_wrapper.Fatal function, in turn,
		//calls a panic.
	}

}
func gen_index_preflight(path string, index_file_name string) bool {

	/*----------------------------------------------------------------
		Called by gen_index, this function checks to see if an index_file_name file
		exists in the current path.

		If it does not:
		we pass preflight

		If index_file_name exists, but the first YAML header tag is AUTOGEN:
		we pass preflight

		Otherwise, we have a user-generated index_file_name, and our user will
		not likely appreciate it being clobbered.
	------------------------------------------------------------------*/

	var front_matter struct {
		/*------------------------------------------------------------
			Still sussing this datastructure out. I need a better understanding
			of how frontmatter works.
		--------------------------------------------------------------*/
		Name string   `yaml:"name"`
		Tags []string `yaml:"tags"`
	}

	log_wrapper.Print("gen_index_preflight called at: ", path)
	// Log that we're running. This logic has been a PITA.

	file_bytes, file_open_error := os.ReadFile(path + "/" + index_file_name)
	//Try to open index_file_name and read the whole thing into file_bytes.

	if file_open_error == nil {
		// If we didn't get an error reading the file, it exists.

		log_wrapper.Print(index_file_name+" exists on ", path)
		// Log that fact.

		frontmatter.Parse(strings.NewReader(string(file_bytes)), &front_matter)
		//Parse the YAML frontmatter from the index_file_name file into the front_matter struct.

		if len(front_matter.Tags) == 0 {
			// This is a sanity check, for file with no tags: tags or no valid YAML
			// header. If we hit that, we're in unknown waters, so fail the preflight.

			log_wrapper.Print("Failed preflight: " + index_file_name + " found, but it contains no tags.")
			// Log that fact

			return false
			//return false right now.
		}

		log_wrapper.Print("Name:", front_matter.Name, "\nTags:", front_matter.Tags[0])
		//Log the first tag in the front_matter.Tags slice. It's the one we care about.

		if front_matter.Tags[0] != "AUTOGEN" {
			log_wrapper.Print("No AUTOGEN tag in header")
			//If there's an index_file_name file without an AUTOGEN tag at the beginning
			//we fail preflight.

			log_wrapper.Print("Failed preflight: " + index_file_name + " with no AUTOGEN flag.")

			return false
			//short-circuit the rest of the function and return false right now.

		} else {
			//If we're here, there's an index_file_name file WITH an AUTOGEN tag at the beginning
			//so we pass preflight.

			log_wrapper.Print("Passed Preflight. There is an " + index_file_name + ", but it AUTOGEN. Clobber away.")
		}

	} else {
		//If we're here, there was no index_file_name file (or some other file error happend. eek)
		//so we pass preflight. Just fall through to the main return.

		log_wrapper.Print("No " + index_file_name + " found.")
		//Log that fact.

	}

	return true
	//Either there is no index_file_name file on this path, or it has the AUTOGEN tag.
	//Return true.
}

func gen_index(path string, index_file_name string, wg *sync.WaitGroup) {
	/*-----------------------------------------------------------------
	This is the function that does all the heavy lifting.

	We are probably in our own thread, so increase the waitgroup value by 1.

	If we pass the preflight check:
		Create the index_file_name file.
		Defer closing it.

		Write a YAML header with at least a tags: entry and the AUTOGEN tag.

		Iterate through the files in this path.
			If it's a directory:
				kick a new go thread with path+our filename.
				process the directory name into our CURRENT index_file_name
				as a link to the subdirectory's index_file_name.

			Else
				Process the filename by type into a link in our current
				index_file_name file. index_file_name should be ignored, since we handle
				it separately. Kludgy.


	Decrease the waitgroup value by 1.
	Aaaand we're done.
	--------------------------------------------------------------------*/

	log_wrapper.Print("gen_index called at path ", path)
	//Log that we're running. It gets hinky when we're multithreaded.
	log_wrapper.Print("index_file_name presently: ", index_file_name)

	wg.Add(1)
	//Increase the waitgroup by 1 for this thread.

	/*------------------------------------------------------------------
	Pre-flight check the directory in path
	-------------------------------------------------------------------*/

	if !gen_index_preflight(path, index_file_name) {
		// If path contains an index_file_name file that's not tagged AUTOGEN,

		log_wrapper.Print("Preflight Check Failed:", path+"/"+index_file_name)
		// Log that fact

		wg.Done()
		return
		// Lower the wait group by one and return.
	}

	//If we get to this point in a given run, we passed the preflight check.

	/*-----------------------------------------------------------
	Create the index_file_name file and add the AUTOGEN header to it.
	------------------------------------------------------------*/
	output_file, err := os.Create(path + "/" + index_file_name)
	on_error_die(err, "Unable to create "+index_file_name+" on "+path)
	//open the output file. Panic if we fail.

	defer output_file.Close()
	//close the output_file when this function exits.

	_, err = io.WriteString(output_file, "---\ntags:\n- AUTOGEN\n---\n\n")
	on_error_die(err, "Failed writing the header to "+index_file_name+" on path "+path)
	//write the output file's tags and AUTOGEN header. Panic if we fail.

	log_wrapper.Print("opened "+index_file_name+" in", path)

	/*----------------------------------------------------------
	Populate the index_file_name file
	-----------------------------------------------------------*/

	log_wrapper.Print("Passed Preflight Check")
	//Log that fact.

	dir_file_handle, err := os.Open(path)
	on_error_die(err, "Unable to open "+path)
	//open path. Panic if we fail

	directory, err := dir_file_handle.ReadDir(0)
	on_error_die(err, "Unable to read in "+path)
	//read in the whole directory at path. Panic if we fail
	//This could get ugly if the directories are big enough.

	sort.Slice(directory,
		func(i, j int) bool {
			return directory[i].Name() < directory[j].Name()
		})
	//Yanno, according to the API, readdir should return the directory already sorted.
	//We have to sort it by name anyway. Go fig.

	for index := range directory {
		file := directory[index]
		file_name_type := strings.Split(file.Name(), ".")
		//split the filename from the type extention

		log_wrapper.Print("Processing ", file.Name())

		if file.IsDir() {
			//do all this if we're dealing with a directory.

			log_wrapper.Print(file.Name(), " is a directory")

			if (file_name_type[0] != "") &&
				((unicode.IsUpper(rune(file_name_type[0][0]))) ||
					(unicode.IsNumber(rune(file_name_type[0][0])))) {
				//filter the directory name and type for sanity.
				//the name can not be an empty string.
				//the name must be capitalized or begin with a number.

				go gen_index(path+"/"+file_name_type[0], index_file_name, wg)
				//new thread with gen_index at this new directory.

				_, err := io.WriteString(output_file, format_dir_link(file_name_type[0], index_file_name))
				on_error_die(err, "Unable to write Directory Link")
				//write the link to the directory's index_file_name. Panic if we fail.
			}

		} else {
			//So we're not a directory if we're here.

			switch file_name_type[1] {
			//file_name_type is a slice with the file name in [0]
			//and the type extension in [1]. So we're switching on
			//the type extension.

			case "md":
				{
					log_wrapper.Print(file.Name(), " is a Markdown file.")
					_, err := io.WriteString(output_file, format_mkd_link(file_name_type[0]))
					on_error_die(err, "Unable to write markdown Link")
					//write the markdown file link using the fmt_mkd_link function on the file name
					//without the extension. Panic if the write fails.

				}
			case "jpg", "jpeg", "gif", "png":
				{
					log_wrapper.Print(file.Name(), " is an image file.")
					_, err := io.WriteString(output_file, format_img_link(file_name_type[0], file_name_type[1]))
					on_error_die(err, "Unable to write image Link")
					//Write the image link using format_image_link and
					//*both* elements of the file_name_type slice.
					//Kludge alert.
				}

			default:
				{
					log_wrapper.Print(file.Name(), " is in an unsupported format.")
					//If we get here, we didn't match any known types.
					//So we don't write any links. Ignore it and move on.
				}

			}

		}

	}

	wg.Done()
	return
	// Lower the wait group by one and return.

}
