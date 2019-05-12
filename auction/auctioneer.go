package auction

import (
	"../p1"
	"encoding/json"
	"io/ioutil"
	// "log"
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

type Result struct {
	Finalized      bool
	MinerID        int
	BidderID       int
	TransactionNum int
	Price          int
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
	mpt.Insert("ItemID", strconv.Itoa(A.ItemNum+1))
	mpt.Insert("Name", itemInfo.Name)
	mpt.Insert("Description", itemInfo.Description)
	mpt.Insert("Price", strconv.Itoa(itemInfo.Price))
	mpt.Insert("End", strconv.FormatInt(itemInfo.End, 16))
	return mpt, nil
}

func (A *Auctioneer) DetermineWinner(itemData ItemData) p1.MerklePatriciaTrie {
	maxTx := Transaction{Detail: BidDetail{BidInfo: BidInfo{Price: itemData.Iteminfo.Info.Price}}}
	for _, tx := range itemData.Trans {
		if tx.Detail.BidInfo.Price > maxTx.Detail.BidInfo.Price {
			maxTx = tx
		}
	}
	var mpt p1.MerklePatriciaTrie
	mpt.Insert("Type", "Result")
	mpt.Insert("MinerID", strconv.Itoa(maxTx.MinerID))
	mpt.Insert("BidderID", strconv.Itoa(maxTx.Detail.BidderID))
	mpt.Insert("Num", strconv.Itoa(maxTx.Detail.Num))
	mpt.Insert("AuctioneerID", strconv.Itoa(maxTx.Detail.BidInfo.AuctioneerID))
	mpt.Insert("ItemID", strconv.Itoa(maxTx.Detail.BidInfo.ItemID))
	mpt.Insert("Price", strconv.Itoa(maxTx.Detail.BidInfo.Price))
	return mpt
}
