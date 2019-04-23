package auction

import (
	"../p2"
	// "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ItemData struct {
	It    Item
	Trans []Transaction
}

type Transaction struct {
	// Bidder transaction
}

type Bidder struct {
	ID      int
	Address string
}

func (B *Bidder) Items() {
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

	for _, chian := range chains {
		parseChain(chian)
	}
}

func parseChain(chain []string) {
	timeNow := time.Now().Unix()
	itemsData := make(map[string]ItemData)
	for i := len(chain); i > 0; i-- {
		blockString := chain[i-1]
		var block p2.Block
		if err := json.Unmarshal([]byte(blockString), &block); err != nil {
			fmt.Println("Parse JSON failed")
			break
		}
		parseBlockData(block, itemsData, timeNow)
	}
	printChain(itemsData)
}

func parseBlockData(block p2.Block, itemsData map[string]ItemData, timeNow int64) {
	if typeName, err := block.Value.Get("Type"); err == nil {
		if typeName == "ItemInfo" {
			if endString, err := block.Value.Get("End"); err == nil {
				if end, err := strconv.ParseInt(endString, 16, 64); err == nil {
					if end > timeNow {
						aucIDString, _ := block.Value.Get("Auctioneer")
						aucID, _ := strconv.Atoi(aucIDString)
						IDString, _ := block.Value.Get("ID")
						ID, _ := strconv.Atoi(IDString)
						name, _ := block.Value.Get("Name")
						desciption, _ := block.Value.Get("Description")
						priceString, _ := block.Value.Get("Price")
						price, _ := strconv.Atoi(priceString)
						var transactions []Transaction
						itemData := ItemData{Item{aucID, ID, name, desciption, price, end}, transactions}
						itemsData[aucIDString+"--"+IDString] = itemData
					}
				}
			}
		} else if typeName == "Transcation" {

		}
	}
}

func printChain(itemsData map[string]ItemData) {

}
