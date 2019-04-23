package auction

import (
	"../p1"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type item struct {
	Name        string
	Description string
	Price       int
	End         int64
}

type auctioneer struct {
	ID      int
	Address string
	ItemNum int
}

func (A *auctioneer) PostItem(it item) {
	A.ItemNum = A.ItemNum + 1
	var mpt p1.MerklePatriciaTrie
	mpt.Insert("Type", "Item")
	mpt.Insert("Auctioneer", strconv.Itoa(A.ID))
	mpt.Insert("ID", strconv.Itoa(A.ItemNum))
	mpt.Insert("Name", it.Name)
	mpt.Insert("Description", it.Description)
	mpt.Insert("Price", strconv.Itoa(it.Price))
	mpt.Insert("End", strconv.FormatInt(it.End, 16))
	mptJSON, err := json.Marshal(mpt)
	if err != nil {
		// handle error
	}
	body := bytes.NewBuffer(mptJSON)
	_, err = http.Post(A.Address+"/post", "application/json", body)
	if err != nil {
		fmt.Println("Post Failed!")
	}
}
