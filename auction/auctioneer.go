package auction

import (
	"../p1"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Item struct {
	AuctioneerID int
	ItemID       int
	Info         ItemInfo
}

type ItemInfo struct {
	Name        string
	Description string
	Price       int
	End         int64
}

type Auctioneer struct {
	ID      int
	Address string
	ItemNum int
}

func (A *Auctioneer) PostItem(r *http.Request) (p1.MerklePatriciaTrie, error) {
	body, err := ioutil.ReadAll(r.Body)
	var mpt p1.MerklePatriciaTrie
	defer r.Body.Close()
	if err != nil {
		return mpt, err
	}
	var itemInfo ItemInfo
	if err := json.Unmarshal(body, &itemInfo); err != nil {
		return mpt, err
	}
	mpt.Insert("Type", "ItemInfo")
	mpt.Insert("AuctioneerID", strconv.Itoa(A.ID))
	mpt.Insert("ItemID", strconv.Itoa(A.ItemNum))
	mpt.Insert("Name", itemInfo.Name)
	mpt.Insert("Description", itemInfo.Description)
	mpt.Insert("Price", strconv.Itoa(itemInfo.Price))
	mpt.Insert("End", strconv.FormatInt(itemInfo.End, 16))
	return mpt, nil
}
