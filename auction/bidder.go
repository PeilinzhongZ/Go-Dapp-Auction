package auction

import (
	"../p1"
	"encoding/json"
	"io/ioutil"
	// "log"
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
	Price        int
}

type Bidder struct {
	ID      int
	Address string
	BidList []string
}

func (B *Bidder) ListItem(chainsData [][]p1.MerklePatriciaTrie, auctioneerID, itemID int) []ItemData {
	var list []ItemData
	for _, chainData := range chainsData {
		var itemData ItemData
		var ok bool
		for i := len(chainData) - 1; i > 0; i-- {
			mpt := chainData[i-1]
			if !ok {
				itemData, ok = findItem(mpt, auctioneerID, itemID)
			} else {
				trans, ok2 := findTrans(mpt, auctioneerID, itemID)
				if ok2 {
					itemData.Trans = append(itemData.Trans, trans)
				}
			}
		}
		if ok {
			list = append(list, itemData)
		}
	}
	return list
}

func findItem(mpt p1.MerklePatriciaTrie, auctioneerID, itemID int) (ItemData, bool) {
	if typeName, err := mpt.Get("Type"); err == nil {
		if typeName == "ItemInfo" {
			aucIDString, _ := mpt.Get("AuctioneerID")
			aucID, _ := strconv.Atoi(aucIDString)
			if aucID == auctioneerID {
				IDString, _ := mpt.Get("ItemID")
				ID, _ := strconv.Atoi(IDString)
				if ID == itemID {
					name, _ := mpt.Get("Name")
					desciption, _ := mpt.Get("Description")
					priceString, _ := mpt.Get("Price")
					endString, _ := mpt.Get("End")
					end, _ := strconv.ParseInt(endString, 16, 64)
					price, _ := strconv.Atoi(priceString)
					var transactions []Transaction
					return ItemData{Item{aucID, ID, ItemInfo{name, desciption, price, end}}, transactions}, true
				}
				return ItemData{}, false
			}
			return ItemData{}, false
		}
		return ItemData{}, false
	}
	return ItemData{}, false
}

func findTrans(mpt p1.MerklePatriciaTrie, auctioneerID, itemID int) (Transaction, bool) {
	if typeName, err := mpt.Get("Type"); err == nil {
		if typeName == "Transaction" {
			auctioneerIDStr, _ := mpt.Values["AuctioneerID"]
			aucID, _ := strconv.Atoi(auctioneerIDStr)
			if aucID == auctioneerID {
				itemIDStr, _ := mpt.Values["ItemID"]
				ID, _ := strconv.Atoi(itemIDStr)
				if ID == itemID {
					minerIDStr, _ := mpt.Values["MinerID"]
					minerID, _ := strconv.Atoi(minerIDStr)
					bidderIDStr, _ := mpt.Values["BidderID"]
					bidderID, _ := strconv.Atoi(bidderIDStr)
					bidStr, _ := mpt.Values["Price"]
					bid, _ := strconv.Atoi(bidStr)
					transaction := Transaction{minerID, BidDetail{bidderID, BidInfo{aucID, ID, bid}}}
					return transaction, true
				}
				return Transaction{}, false
			}
			return Transaction{}, false
		}
		return Transaction{}, false
	}
	return Transaction{}, false
}

func (B *Bidder) ListItems(chainsData [][]p1.MerklePatriciaTrie) []map[string]*ItemData {
	var itemDataList []map[string]*ItemData
	for _, chainData := range chainsData {
		itemsData := make(map[string]*ItemData)
		timeNow := time.Now().Unix()
		for i := len(chainData) - 1; i > 0; i-- {
			mpt := chainData[i-1]
			parseMptData(mpt, itemsData, timeNow)
		}
		itemDataList = append(itemDataList, itemsData)
	}
	return itemDataList
}

func (B *Bidder) PostBid(bidderID int32, r *http.Request) (string, string, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return "", "", err
	}
	var bidInfo BidInfo
	if err := json.Unmarshal(body, &bidInfo); err != nil {
		return "", "", err
	}
	bidDetail := BidDetail{int(bidderID), bidInfo}
	bytes, err := json.Marshal(bidDetail)
	if err != nil {
		return "", "", err
	}
	itemIDValue := strconv.Itoa(bidInfo.AuctioneerID) + "-" + strconv.Itoa(bidInfo.ItemID)
	return itemIDValue, string(bytes), nil
}

func parseMptData(mpt p1.MerklePatriciaTrie, itemsData map[string]*ItemData, timeNow int64) {
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
					itemsData[aucIDString+"-"+IDString] = &itemData
					// }
				}
			}
		} else if typeName == "Transaction" {
			minerIDStr, _ := mpt.Values["MinerID"]
			minerID, _ := strconv.Atoi(minerIDStr)
			bidderIDStr, _ := mpt.Values["BidderID"]
			bidderID, _ := strconv.Atoi(bidderIDStr)
			auctioneerIDStr, _ := mpt.Values["AuctioneerID"]
			auctioneerID, _ := strconv.Atoi(auctioneerIDStr)
			itemIDStr, _ := mpt.Values["ItemID"]
			itemID, _ := strconv.Atoi(itemIDStr)
			bidStr, _ := mpt.Values["Price"]
			bid, _ := strconv.Atoi(bidStr)
			if itemdata, ok := itemsData[auctioneerIDStr+"-"+itemIDStr]; ok {
				transaction := Transaction{minerID, BidDetail{bidderID, BidInfo{auctioneerID, itemID, bid}}}
				itemdata.Trans = append(itemdata.Trans, transaction)
			}
		}
	}
}
