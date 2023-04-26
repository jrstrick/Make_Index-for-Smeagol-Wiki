package log_wrapper

import "log"

var DEBUG bool

func Set_DEBUG(pointer *bool) {
	if *pointer {
		DEBUG = true
	} else {
		DEBUG = false
	}
}

func init() {
	DEBUG = false
}
func Print(v ...any) {
	if DEBUG {
		log.Print(v)
	}
}
func Fatal(v ...any) {
	log.Fatal(v)
	//always log fatals.
}
