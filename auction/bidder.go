package auction

import (
	"../p1"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ItemData struct {
	Iteminfo Item
	Trans    []Transaction
}

type Transaction struct {
	MinerID int
	Detail  BidDetail
}

type BidDetail struct {
	BidderID int
	BidInfo  BidInfo
}

type BidInfo struct {
	AuctioneerID int
	ItemID       int
	Bid          int
}

type Bidder struct {
	ID      int
	Address string
}

func (B *Bidder) ListItems(chainsData [][]p1.MerklePatriciaTrie) []map[string]ItemData {
	var itemDataList []map[string]ItemData
	for _, chainData := range chainsData {
		itemsData := make(map[string]ItemData)
		timeNow := time.Now().Unix()
		for i := len(chainData) - 1; i > 0; i-- {
			mpt := chainData[i-1]
			parseMptData(mpt, itemsData, timeNow)
		}
		itemDataList = append(itemDataList, itemsData)
	}
	return itemDataList
}

func (B *Bidder) PostBid(bidderID int32, r *http.Request) (string, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return "", err
	}
	var bidInfo BidInfo
	if err := json.Unmarshal(body, &bidInfo); err != nil {
		return "", err
	}
	bidDetail := BidDetail{int(bidderID), bidInfo}
	bytes, err := json.Marshal(bidDetail)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func parseMptData(mpt p1.MerklePatriciaTrie, itemsData map[string]ItemData, timeNow int64) {
	if typeName, err := mpt.Get("Type"); err == nil {
		if typeName == "ItemInfo" {
			if endString, err := mpt.Get("End"); err == nil {
				if end, err := strconv.ParseInt(endString, 16, 64); err == nil {
					// if end > timeNow {
					aucIDString, _ := mpt.Get("AuctioneerID")
					aucID, _ := strconv.Atoi(aucIDString)
					IDString, _ := mpt.Get("ItemID")
					ID, _ := strconv.Atoi(IDString)
					name, _ := mpt.Get("Name")
					desciption, _ := mpt.Get("Description")
					priceString, _ := mpt.Get("Price")
					price, _ := strconv.Atoi(priceString)
					var transactions []Transaction
					itemData := ItemData{Item{aucID, ID, ItemInfo{name, desciption, price, end}}, transactions}
					itemsData[aucIDString+"-"+IDString] = itemData
					// }
				}
			}
		} else if typeName == "Transcation" {

		}
	}
}
