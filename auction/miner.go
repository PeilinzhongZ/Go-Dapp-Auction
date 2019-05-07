package auction

import (
	"../p1"
	"encoding/json"
	"strconv"
)

type Miner struct {
	ID      int
	Address string
	Trans   []string
	IsMiner bool
}

func (M *Miner) Min(bidJson string) (string, p1.MerklePatriciaTrie, error) {
	var bidDetail BidDetail
	var mpt p1.MerklePatriciaTrie
	if err := json.Unmarshal([]byte(bidJson), &bidDetail); err != nil {
		return "", mpt, err
	}
	mpt.Insert("Type", "Transaction")
	mpt.Insert("MinerID", strconv.Itoa(M.ID))
	mpt.Insert("BidderID", strconv.Itoa(bidDetail.BidderID))
	mpt.Insert("AuctioneerID", strconv.Itoa(bidDetail.BidInfo.AuctioneerID))
	mpt.Insert("ItemID", strconv.Itoa(bidDetail.BidInfo.ItemID))
	mpt.Insert("Price", strconv.Itoa(bidDetail.BidInfo.Price))
	itemIDValue := strconv.Itoa(bidDetail.BidInfo.AuctioneerID) + "-" + strconv.Itoa(bidDetail.BidInfo.ItemID)
	return itemIDValue, mpt, nil
}
