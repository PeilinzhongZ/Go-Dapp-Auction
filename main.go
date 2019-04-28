package main

import (
	// "./auction"
	"./p3"
	// "fmt"
	"log"
	"net/http"
	"os"
	// "strconv"
	// "time"
)

func main() {
	router := p3.NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":6686", router))
	}
}
