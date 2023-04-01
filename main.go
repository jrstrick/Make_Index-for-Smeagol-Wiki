package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func path_is_dir(path string) bool {
	log.Print("path_is_dir called.")
	a, err := os.Stat(path)
	if err != nil {
		log.Fatal("error statting path ", path)
	}
	return a.IsDir()
}
func main() {
	arg := os.Args

	wg := new(sync.WaitGroup)
	//set up the waitgroup

	if len(arg) >= 2 {
		if path_is_dir(arg[1]) {
			log.Print(arg[1], " looks ok. Running")
			gen_index(arg[1], wg)
		} else {
			log.Fatal(arg[1], "is invalid.")
		}
	} else {
		fmt.Println("Usage: make_index dirname\n\t(where dirname is the directory of your wiki)")
	}

	wg.Wait()
	//wait for the waitgroup to reach zero before exiting.
}
