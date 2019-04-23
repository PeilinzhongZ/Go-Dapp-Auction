package main

import (
	"./p3"
	"fmt"
	// "log"
	// "./auction"
	"net/http"
	"os"
)

func main() {
	router := p3.NewRouter()
	if len(os.Args) > 1 {
		go http.ListenAndServe(":"+os.Args[1], router)
		_, err := http.Get("http://localhost:" + os.Args[1] + "/start")
		if err != nil {
			fmt.Println("Failed")
			return
		}
		// test
		// A := auction.auctioneer{6687, "http://localhost:6687", 1}
		// A.PostItem(auction.item{"a", "b", 1, int64(1)})
	} else {
		go http.ListenAndServe(":6686", router)
		_, err := http.Get("http://localhost:6686/start?first=true")
		if err != nil {
			fmt.Println("Failed")
			return
		}
	}
	// for true {

	// }
	// start server
}
