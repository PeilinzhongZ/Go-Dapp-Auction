package auction

import (
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
	Result   Result
}

type Transaction struct {
	MinerID int
	Detail  BidDetail
}

type BidDetail struct {
	BidderID int
	Num      int
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

func (B *Bidder) ListItem(chainsData [][]TrieWithTime, auctioneerID, itemID int) []ItemData {
	var list []ItemData
	for _, chainData := range chainsData {
		var itemData ItemData
		var ok bool
		for i := len(chainData) - 1; i > 0; i-- {
			twt := chainData[i-1]
			if !ok {
				itemData, ok = findItem(twt, auctioneerID, itemID)
			} else {
				trans, result, isTX, ok2 := findTrans(twt, auctioneerID, itemID)
				if ok2 {
					if isTX {
						// if itemData.Iteminfo.Info.End > twt.Timestamp {
						itemData.Trans = append(itemData.Trans, trans)
						// }
					} else {
						itemData.Result = result
					}
				}
			}
		}
		if ok {
			list = append(list, itemData)
		}
	}
	return list
}

func findItem(twt TrieWithTime, auctioneerID, itemID int) (ItemData, bool) {
	if typeName, err := twt.Trie.Get("Type"); err == nil {
		if typeName == "ItemInfo" {
			aucIDString, _ := twt.Trie.Values["AuctioneerID"]
			aucID, _ := strconv.Atoi(aucIDString)
			if aucID == auctioneerID {
				IDString, _ := twt.Trie.Values["ItemID"]
				ID, _ := strconv.Atoi(IDString)
				if ID == itemID {
					name, _ := twt.Trie.Values["Name"]
					desciption, _ := twt.Trie.Values["Description"]
					priceString, _ := twt.Trie.Values["Price"]
					endString, _ := twt.Trie.Values["End"]
					end, _ := strconv.ParseInt(endString, 16, 64)
					price, _ := strconv.Atoi(priceString)
					var transactions []Transaction
					return ItemData{Iteminfo: Item{aucID, ID, ItemInfo{name, desciption, price, end}}, Trans: transactions}, true
				}
				return ItemData{}, false
			}
			return ItemData{}, false
		}
		return ItemData{}, false
	}
	return ItemData{}, false
}

func findTrans(twt TrieWithTime, auctioneerID, itemID int) (Transaction, Result, bool, bool) {
	if typeName, err := twt.Trie.Get("Type"); err == nil {
		auctioneerIDStr, _ := twt.Trie.Values["AuctioneerID"]
		aucID, _ := strconv.Atoi(auctioneerIDStr)
		if aucID == auctioneerID {
			itemIDStr, _ := twt.Trie.Values["ItemID"]
			ID, _ := strconv.Atoi(itemIDStr)
			if ID == itemID {
				minerIDStr, _ := twt.Trie.Values["MinerID"]
				minerID, _ := strconv.Atoi(minerIDStr)
				bidderIDStr, _ := twt.Trie.Values["BidderID"]
				bidderID, _ := strconv.Atoi(bidderIDStr)
				bidStr, _ := twt.Trie.Values["Price"]
				bid, _ := strconv.Atoi(bidStr)
				numStr, _ := twt.Trie.Values["Num"]
				num, _ := strconv.Atoi(numStr)
				if typeName == "Transaction" {
					transaction := Transaction{minerID, BidDetail{bidderID, num, BidInfo{aucID, ID, bid}}}
					return transaction, Result{}, true, true
				} else if typeName == "Result" {
					result := Result{true, minerID, bidderID, num, bid}
					return Transaction{}, result, false, true
				}
			}
			return Transaction{}, Result{}, true, false
		}
		return Transaction{}, Result{}, true, false
	}
	return Transaction{}, Result{}, true, false
}

func (B *Bidder) ListItems(chainsData [][]TrieWithTime) []map[string]*ItemData {
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
	bidDetail := BidDetail{int(bidderID), len(B.BidList), bidInfo}
	bytes, err := json.Marshal(bidDetail)
	if err != nil {
		return "", "", err
	}
	itemIDValue := strconv.Itoa(bidInfo.AuctioneerID) + "-" + strconv.Itoa(bidInfo.ItemID)
	return itemIDValue, string(bytes), nil
}

func parseMptData(twt TrieWithTime, itemsData map[string]*ItemData, timeNow int64) {
	if typeName, err := twt.Trie.Get("Type"); err == nil {
		if typeName == "ItemInfo" {
			if endString, err := twt.Trie.Get("End"); err == nil {
				if end, err := strconv.ParseInt(endString, 16, 64); err == nil {
					// if end > timeNow {
					aucIDString, _ := twt.Trie.Values["AuctioneerID"]
					aucID, _ := strconv.Atoi(aucIDString)
					IDString, _ := twt.Trie.Values["ItemID"]
					ID, _ := strconv.Atoi(IDString)
					name, _ := twt.Trie.Values["Name"]
					desciption, _ := twt.Trie.Values["Description"]
					priceString, _ := twt.Trie.Values["Price"]
					price, _ := strconv.Atoi(priceString)
					var transactions []Transaction
					itemData := ItemData{Iteminfo: Item{aucID, ID, ItemInfo{name, desciption, price, end}}, Trans: transactions}
					itemsData[aucIDString+"-"+IDString] = &itemData
					// }
				}
			}
		} else {
			minerIDStr, _ := twt.Trie.Values["MinerID"]
			minerID, _ := strconv.Atoi(minerIDStr)
			bidderIDStr, _ := twt.Trie.Values["BidderID"]
			bidderID, _ := strconv.Atoi(bidderIDStr)
			auctioneerIDStr, _ := twt.Trie.Values["AuctioneerID"]
			auctioneerID, _ := strconv.Atoi(auctioneerIDStr)
			itemIDStr, _ := twt.Trie.Values["ItemID"]
			numStr, _ := twt.Trie.Values["Num"]
			num, _ := strconv.Atoi(numStr)
			itemID, _ := strconv.Atoi(itemIDStr)
			bidStr, _ := twt.Trie.Values["Price"]
			bid, _ := strconv.Atoi(bidStr)
			if itemdata, ok := itemsData[auctioneerIDStr+"-"+itemIDStr]; ok {
				if typeName == "Transaction" {
					// if itemdata.Iteminfo.Info.End > twt.Timestamp {
					transaction := Transaction{minerID, BidDetail{bidderID, num, BidInfo{auctioneerID, itemID, bid}}}
					itemdata.Trans = append(itemdata.Trans, transaction)
					// }
				} else if typeName == "Result" {
					result := Result{true, minerID, bidderID, num, bid}
					itemdata.Result = result
				}
			}
		}
	}
}
