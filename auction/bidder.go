package auction

import (
	// "../p2"
	// "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "strconv"
)

type bidder struct {
	ID      int
	Address string
}

func (B *bidder) Items() {
	resp, err := http.Get(B.Address + "/canonical")
	if err != nil {
		fmt.Println("Server failed")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read failed")
	}
	var chains [][]string
	if err := json.Unmarshal(body, &chains); err != nil {
		fmt.Println("Parse failed")
	}

}
