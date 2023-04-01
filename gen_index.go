package main

import (
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode"

	"github.com/adrg/frontmatter"
)

func on_error_die(err error, text string) {
	/*-----------------------------------------------------------------
	Just a quick function that checks the value of a standard error
	type and if it's not nil, logs a reason string and panics.
	-----------------------------------------------------------------*/

	if err != nil {
		//if we have an error...

		log.Fatal("PANIC: "+text+"\n\t Reason: ", err)
		//call log.Fatal with the text from the function call
		//and the error itself. The log.Fatal function, in turn,
		//calls a panic.
	}

}
func gen_index_preflight(path string) bool {

	/*----------------------------------------------------------------
		Called by gen_index, this function checks to see if an Index.md file
		exists in the current path.

		If it does not:
		we pass preflight

		If Index.md exists, but the first YAML header tag is AUTOGEN:
		we pass preflight

		Otherwise, we have a user-generated Index.md, and our user will
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

	log.Print("gen_index_preflight called at: ", path)
	// Log that we're running. This logic has been a PITA.

	file_bytes, file_open_error := os.ReadFile(path + "/Index.md")
	//Try to open Index.md and read the whole thing into file_bytes.

	if file_open_error == nil {
		// If we didn't get an error reading the file, it exists.

		log.Print("Index.md exists on ", path)
		// Log that fact.

		frontmatter.Parse(strings.NewReader(string(file_bytes)), &front_matter)
		//Parse the YAML frontmatter from the Index.md file into the front_matter struct.

		if len(front_matter.Tags) == 0 {
			// This is a sanity check, for file with no tags: tags or no valid YAML
			// header. If we hit that, we're in unknown waters, so fail the preflight.

			log.Print("Failed preflight: Index.md found, but it contains no tags.")
			// Log that fact

			return false
			//return false right now.
		}

		log.Print("Name:", front_matter.Name, "\nTags:", front_matter.Tags[0])
		//Log the first tag in the front_matter.Tags slice. It's the one we care about.

		if front_matter.Tags[0] != "AUTOGEN" {
			log.Print("No AUTOGEN tag in header")
			//If there's an Index.md file without an AUTOGEN tag at the beginning
			//we fail preflight.

			log.Print("Failed preflight: Index.md with no AUTOGEN flag.")

			return false
			//short-circuit the rest of the function and return false right now.

		} else {
			//If we're here, there's an Index.md file WITH an AUTOGEN tag at the beginning
			//so we pass preflight.

			log.Print("Passed Preflight. There is an Index.md, but it AUTOGEN. Clobber away.")
		}

	} else {
		//If we're here, there was no Index.md file (or some other file error happend. eek)
		//so we pass preflight. Just fall through to the main return.

		log.Print("No Index.md found.")
		//Log that fact.

	}

	return true
	//Either there is no Index.md file on this path, or it has the AUTOGEN tag.
	//Return true.
}

func gen_index(path string, wg *sync.WaitGroup) {
	/*-----------------------------------------------------------------
	This is the function that does all the heavy lifting.

	We are probably in our own thread, so increase the waitgroup value by 1.

	If we pass the preflight check:
		Create the Index.md file.
		Write a YAML header with at least a tags: entry and the AUTOGEN tag.

		Iterate through the files in this path.
			If it's a directory:
				kick a new go thread with path+our filename.
				process the directory name into our CURRENT Index.md
				as a link to the subdirectory's Index.md.

			Else
				Process the filename by type into a link in our current
				Index.md file. Index.md should be ignored, since we handle
				it separately. Kludgy.

		Close Index.md.

	Decrease the waitgroup value by 1.
	Aaaand we're done.
	--------------------------------------------------------------------*/

	log.Print("gen_index called at path ", path)
	//Log that we're running. It gets hinky when we're multithreaded.

	wg.Add(1)
	//Increase the waitgroup by 1 for this thread.

	/*------------------------------------------------------------------
	Pre-flight check the directory in path
	-------------------------------------------------------------------*/

	if !gen_index_preflight(path) {
		// If path contains an Index.md file that's not tagged AUTOGEN,

		log.Print("Preflight Check Failed:", path+"/Index.md")
		// Log that fact

		wg.Done()
		return
		// Lower the wait group by one and return.
	}

	//If we get to this point in a given run, we passed the preflight check.

	/*-----------------------------------------------------------
	Create the Index.md file and add the AUTOGEN header to it.
	------------------------------------------------------------*/
	output_file, err := os.Create(path + "/Index.md")
	//output_file := os.Stdout

	on_error_die(err, "Unable to create Index.md on "+path)

	_, err = io.WriteString(output_file, "---\ntags:\n- AUTOGEN\n---\n\n")
	on_error_die(err, "Failed writing the header to Index.md on path "+path)

	log.Print("opened Index.md in", path)

	/*----------------------------------------------------------
	Populate the Index.md file
	-----------------------------------------------------------*/

	log.Print("Passed Preflight Check")
	//Log that fact.

	dir_file_handle, err := os.Open(path)
	on_error_die(err, "Unable to open "+path)
	//open path. Panic if we fail

	directory, err := dir_file_handle.ReadDir(0)
	on_error_die(err, "Unable to read in "+path)

	//read in the whole directory at path. Panic if we fail

	sort.Slice(directory,
		func(i, j int) bool {
			return directory[i].Name() < directory[j].Name()
		})

	//Yanno, according to the API, readdir should return the directory already sorted.
	//We have to sort it by name anyway. Go fig.

	for index := range directory {
		file := directory[index]
		file_name_type := strings.Split(file.Name(), ".")

		log.Print("Processing ", file.Name())

		if file.IsDir() {
			log.Print(file.Name(), " is a directory")

			if (file_name_type[0] != "") &&
				((unicode.IsUpper(rune(file_name_type[0][0]))) ||
					(unicode.IsNumber(rune(file_name_type[0][0])))) {
				go gen_index(path+"/"+file_name_type[0], wg)

				_, err := io.WriteString(output_file, format_dir_link(file_name_type[0]))
				on_error_die(err, "Unable to write Directory Link")
			}

		} else {

			switch file_name_type[1] {

			case "md":
				{

					log.Print(file.Name(), " is a Markdown file.")
					_, err := io.WriteString(output_file, format_mkd_link(file_name_type[0]))
					on_error_die(err, "Unable to write markdown Link")

				}
			case "jpg", "gif", "png":
				{
					log.Print(file.Name(), " is an image file.")
					_, err := io.WriteString(output_file, format_img_link(file_name_type[0], file_name_type[1]))
					on_error_die(err, "Unable to write image Link")
				}

			default:
				{

					log.Print(file.Name(), " is in an unsupported format.")
				}

			}

		}

	}

	wg.Done()
	output_file.Close()
	return
	// Lower the wait group by one and return.

}
