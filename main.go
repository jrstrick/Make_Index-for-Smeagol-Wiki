package main

import (
	"flag"
	"fmt"
	"make_index/log_wrapper"
	"os"
	"sync"

	"github.com/pelletier/go-toml"
)

// --------------------------------------
// Global Variables

const VERSION = 0.3

func path_is_dir(path string) bool {
	log_wrapper.Print("path_is_dir called.")
	a, err := os.Stat(path)
	if err != nil {
		fmt.Println(path, "not found.")
		log_wrapper.Fatal("error statting path ", path)
	}
	return a.IsDir()
}

func main() {

	var wiki_root string
	var index_file_name string
	var config_file string
	//set variables up for all our flags, to be loaded below.

	fmt.Printf("make_index v%v running...\n", VERSION)

	flag.StringVar(&wiki_root, "wiki_root", "./", "Path to the wiki. Default: ./")
	flag.StringVar(&index_file_name, "index_file", "README.md", "Index file name. (Default: README.md)")
	flag.StringVar(&config_file, "config_file", "smeagol.toml", "Config toml file name. (Default: smeagol.toml)")
	debug := flag.Bool("debug", false, "Enable all the debug messages. Very chatty. (Default: false)")
	flag.Parse()
	log_wrapper.Set_DEBUG(debug)
	//load values into all our flag strings, even if it's only the defaults.
	//these are unix-style flags, passed at the command line with a - sign.

	wg := new(sync.WaitGroup)
	//set up the waitgroup

	if path_is_dir(wiki_root) {
		log_wrapper.Print(wiki_root, "wiki root directory looks ok. Running")
		smeagol_config_bytes, file_open_error := os.ReadFile(wiki_root + "/" + config_file)
		if file_open_error == nil {
			log_wrapper.Print("Found a smeagol.toml")

			smeagol_toml, _ := toml.Load(string(smeagol_config_bytes))

			index_file_name = smeagol_toml.Get("index-page").(string)
			index_file_name += ".md"
			log_wrapper.Print("using " + index_file_name + " as our index-page name per smeagol.toml")
		}

		gen_index(wiki_root, index_file_name, wg)
	} else {
		log_wrapper.Fatal(wiki_root, "is invalid.")
	}

	wg.Wait()
	//wait for the waitgroup to reach zero before exiting.
	fmt.Println("make_index done.")
}
