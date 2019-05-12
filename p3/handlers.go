package p3

import (
	"../p1"
	"../p2"
	"./data"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	// "github.com/gorilla/mux"
	// "io"
	"../auction"
	"bytes"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var TA_SERVER = "http://localhost:6688"
var REGISTER_SERVER = TA_SERVER + "/peer"
var BC_DOWNLOAD_SERVER = TA_SERVER + "/upload"
var SELF_ADDR = "http://localhost:6686"
var FIRST_ADDR = "http://localhost:6686"

var PEERS_SIZE = int32(32)

var SBC data.SyncBlockChain
var Peers data.PeerList
var ifStarted bool

var auctioneer auction.Auctioneer
var bidder auction.Bidder
var miner auction.Miner

func initial() {
	// This function will be executed before everything else.
	// Do some initialization here.
	SBC = data.NewBlockChain()
}

// Register ID, download BlockChain, start HeartBeat
func Start(w http.ResponseWriter, r *http.Request) {
	if !ifStarted {
		id := int32(6686)
		if len(os.Args) > 1 {
			SELF_ADDR = "http://localhost:" + os.Args[1]
			id64, _ := strconv.ParseInt(string(os.Args[1]), 10, 32)
			id = int32(id64)
		}
		initial()
		Peers = data.NewPeerList(id, PEERS_SIZE)
		auctioneer = auction.Auctioneer{int(Peers.GetSelfId()), SELF_ADDR, 0}
		bidder = auction.Bidder{ID: int(Peers.GetSelfId()), Address: SELF_ADDR}
		miner = auction.Miner{ID: int(Peers.GetSelfId()), Address: SELF_ADDR, IsMiner: false}
		min, ok := r.URL.Query()["miner"]
		if ok && min[0] == "true" {
			miner.IsMiner = true
		}
		first, ok := r.URL.Query()["first"]
		if ok && first[0] == "true" {
			var gBlock p2.Block
			gBlock.FirstBlock()
			SBC.Insert(gBlock)
		} else {
			hearbeatData, err := GetPeerMap()
			if err != nil {
				fmt.Fprintf(w, "GetPeerMap error")
				return
			}
			Peers.Add(hearbeatData.Addr, hearbeatData.Id)
			Peers.InjectPeerMapJson(hearbeatData.PeerMapJson, SELF_ADDR)
			blockChainJSON, err := Download()
			if err != nil {
				fmt.Fprintf(w, "Donwload error")
				return
			}
			SBC.UpdateEntireBlockChain(blockChainJSON)
		}
		rand.Seed(time.Now().UnixNano())
		go StartHeartBeat()
		ifStarted = true
	}
}

// Display peerList and sbc
func Show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n%s", Peers.Show(), SBC.Show())
}

// Register to TA's server, get an ID
func Register() (int32, error) {
	resp, err := http.Get(REGISTER_SERVER)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(string(body), 10, 32)
	return int32(i), err
}

func GetPeerMap() (data.HeartBeatData, error) {
	resp, err := http.Get(FIRST_ADDR + "/peerMap?addr=" + SELF_ADDR + "&id=" + strconv.Itoa(int(Peers.GetSelfId())))
	if err != nil {
		return data.HeartBeatData{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data.HeartBeatData{}, err
	}
	var heartBeat data.HeartBeatData
	if err := json.Unmarshal(body, &heartBeat); err != nil {
		return data.HeartBeatData{}, err
	}
	return heartBeat, err
}

// Download blockchain from TA server
func Download() (string, error) {
	peers := Peers.Copy()
	for addr := range peers {
		resp, err := http.Get(addr + "/upload")
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		return string(body), err
	}
	return "", errors.New("error")
}

func UploadPeerMap(w http.ResponseWriter, r *http.Request) {
	addr, ok1 := r.URL.Query()["addr"]
	idString, ok2 := r.URL.Query()["id"]
	id, err := strconv.Atoi(idString[0])
	if ok1 && ok2 && err == nil {
		Peers.Add(addr[0], int32(id))
		peerMapJSON, err := Peers.PeerMapToJson()
		if err != nil {
			// handle error
		}
		heartbeat := data.HeartBeatData{Id: Peers.GetSelfId(), PeerMapJson: peerMapJSON, Addr: SELF_ADDR}
		heartbeatJSON, err := json.Marshal(heartbeat)
		if err != nil {
			// handle error
		}
		fmt.Fprint(w, string(heartbeatJSON))
	}
}

// Upload blockchain to whoever called this method, return jsonStr
func Upload(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := SBC.BlockChainToJson()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprint(w, blockChainJson)
}

// Upload a block to whoever called this method, return jsonStr
func UploadBlock(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	i, err := strconv.ParseInt(path[2], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		block, ok := SBC.GetBlock(int32(i), path[3])
		if ok == false {
			w.WriteHeader(http.StatusNoContent)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, block.Encode())
	}
}

// Received a heartbeat
func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		// handle error
		fmt.Println(err)
		return
	}
	var heartBeat data.HeartBeatData
	if err := json.Unmarshal(body, &heartBeat); err != nil {
		// handle error
		fmt.Println(err)
		return
	}
	Peers.Add(heartBeat.Addr, heartBeat.Id)
	Peers.InjectPeerMapJson(heartBeat.PeerMapJson, SELF_ADDR)
	if heartBeat.IfNewBlock {
		handleNewBlock(heartBeat)
		ForwardHeartbeat(heartBeat)
	} else if heartBeat.IfBid {
		handleNewBid(heartBeat)
		ForwardHeartbeat(heartBeat)
	}
}

func ForwardHeartbeat(heartBeat data.HeartBeatData) {
	if heartBeat.Hops = heartBeat.Hops - 1; heartBeat.Hops != 0 {
		heartBeat.Addr = SELF_ADDR
		heartBeat.Id = Peers.GetSelfId()
		SendHeartBeat(heartBeat)
	}
}

func handleNewBlock(heartBeat data.HeartBeatData) {
	var block p2.Block
	block.Decode(heartBeat.BlockJson)
	if exist := SBC.CheckParentHash(block); exist {
		ok := CheckNonce(block)
		if ok {
			SBC.Insert(block)
		}
	} else {
		AskForBlock(block.Header.Height-1, block.Header.ParentHash)
		if exist = SBC.CheckParentHash(block); exist {
			ok := CheckNonce(block)
			if ok {
				SBC.Insert(block)
			}
		}
	}
}

func handleNewBid(heartBeat data.HeartBeatData) {
	if miner.IsMiner {
		id, mpt, err := miner.Min(heartBeat.BidJson)
		if err != nil {
			log.Println(err)
			return
		}
		go StartTryingNonces(mpt)
		miner.Trans = append(miner.Trans, id)
	}
}

// Ask another server to return a block of certain height and hash
func AskForBlock(height int32, hash string) {
	peers := Peers.Copy()
	for addr := range peers {
		resp, err := http.Get(addr + "/block/" + strconv.Itoa(int(height)) + "/" + hash)
		if err != nil {
			// handle error
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// handle error
				continue
			}
			var block p2.Block
			block.Decode(string(body))
			if exist := SBC.CheckParentHash(block); !exist {
				AskForBlock(block.Header.Height-1, block.Header.ParentHash)
				if exist = SBC.CheckParentHash(block); exist {
					// ok := CheckNonce(block)
					// if ok {
					SBC.Insert(block)
					// }
				}
			} else {
				// ok := CheckNonce(block)
				// if ok {
				SBC.Insert(block)
				// }
			}
			break
		}
	}
}

func CheckNonce(block p2.Block) bool {
	str := sha3.Sum256([]byte(block.Header.ParentHash + block.Header.Nonce + block.Value.Root))
	result := hex.EncodeToString(str[:])
	return strings.HasPrefix(result, "00000")
}

func SendHeartBeat(heartBeatData data.HeartBeatData) {
	list := Peers.Copy()
	heartBeatJSON, err := json.Marshal(heartBeatData)
	if err != nil {
		// handle error
		panic(err)
	}
	for addr := range list {
		go func(addr string) {
			body := bytes.NewBuffer(heartBeatJSON)
			_, err = http.Post(addr+"/heartbeat/receive", "application/json", body)
			if err != nil {
				// handle error
				log.Println(err)
				Peers.Delete(addr)
			}
		}(addr)
	}
}

func StartHeartBeat() {
	for range time.Tick(time.Second * 10) {
		str, err := Peers.PeerMapToJson()
		if err != nil {
			log.Println(err)
			str = ""
		}
		heartBeatData := data.PrepareHeartBeatData(&SBC, Peers.GetSelfId(), str, SELF_ADDR)
		SendHeartBeat(heartBeatData)
	}
}

func PostBid(w http.ResponseWriter, r *http.Request) {
	id, bidJSON, err := bidder.PostBid(Peers.GetSelfId(), r)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	str, err := Peers.PeerMapToJson()
	if err != nil {
		log.Println(err)
		str = ""
	}
	heartBeatData := data.PrepareBidData(Peers.GetSelfId(), str, SELF_ADDR, bidJSON)
	SendHeartBeat(heartBeatData)
	bidder.BidList = append(bidder.BidList, id)
}

func PostItem(w http.ResponseWriter, r *http.Request) {
	mpt, err := auctioneer.PostItem(r)
	if err != nil {
		log.Println(err)
		return
	}
	auctioneer.ItemNum++
	go StartTryingNonces(mpt)
}

func StartTryingNonces(mpt p1.MerklePatriciaTrie) {
	success := false
	for !success {
		latestBlocks, ok := SBC.GetLatestBlocks()
		if ok {
			var x string
			x, success = TryNonces(latestBlocks, mpt.Root)
			if success {
				str, err := Peers.PeerMapToJson()
				if err != nil {
					str = ""
				}
				var block p2.Block
				block.Initial(latestBlocks[0].Header.Height+1, latestBlocks[0].Header.Hash, mpt)
				block.Header.Nonce = x
				SBC.Insert(block)
				heartBeadData := data.NewHeartBeatData(true, Peers.GetSelfId(), block.Encode(), str, SELF_ADDR)
				heartBeadData.Hops = 1
				SendHeartBeat(heartBeadData)
			}
		}
	}
}

func TryNonces(latestBlocks []p2.Block, Root string) (string, bool) {
	var result string
	var x string
	success := true
	for !strings.HasPrefix(result, "00000") {
		if SBC.GetLength() > latestBlocks[0].Header.Height {
			success = false
			break
		}
		bytes := make([]byte, 8)
		rand.Read(bytes)
		x = hex.EncodeToString(bytes)
		resultSum := sha3.Sum256([]byte(latestBlocks[0].Header.Hash + x + Root))
		result = hex.EncodeToString(resultSum[:])
	}
	return x, success
}

func ListItem(w http.ResponseWriter, r *http.Request) {
	aIDStr, ok := r.URL.Query()["auctioneerID"]
	iIDStr, ok2 := r.URL.Query()["itemID"]
	lastBlocks, ok3 := SBC.GetLatestBlocks()
	if !ok3 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	chains := canonicalData(lastBlocks)
	var rawData []byte
	var err error
	if ok && ok2 {
		aID, _ := strconv.Atoi(aIDStr[0])
		iID, _ := strconv.Atoi(iIDStr[0])
		itemData := bidder.ListItem(chains, aID, iID)
		rawData, err = json.Marshal(itemData)
	} else {
		itemsData := bidder.ListItems(chains)
		rawData, err = json.Marshal(itemsData)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(rawData))
}

func FinalizeAuction(w http.ResponseWriter, r *http.Request) {
	iIDStr, ok := r.URL.Query()["itemID"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for true {
		if SBC.CheckCanonical() {
			break
		}
	}
	lastBlocks, ok3 := SBC.GetLatestBlocks()
	if !ok3 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	chains := canonicalData(lastBlocks)
	iID, _ := strconv.Atoi(iIDStr[0])
	itemData := bidder.ListItem(chains, auctioneer.ID, iID)
	winner := auctioneer.DetermineWinner(itemData[0])
	StartTryingNonces(winner)
}

func canonicalData(lastBlocks []p2.Block) [][]auction.TrieWithTime {
	var chains [][]auction.TrieWithTime
	for _, block := range lastBlocks {
		var chain []auction.TrieWithTime
		chain = append(chain, auction.TrieWithTime{block.Value, block.Header.Timestamp})
		for block.Header.Height != 0 {
			block, _ = SBC.GetParentBlock(block)
			chain = append(chain, auction.TrieWithTime{block.Value, block.Header.Timestamp})
		}
		chains = append(chains, chain)
	}
	return chains
}

func Canonical(w http.ResponseWriter, r *http.Request) {
	lastBlocks, ok := SBC.GetLatestBlocks()
	if !ok {
		//handle error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var chains string
	for i, block := range lastBlocks {
		chain := "Chain" + strconv.Itoa(i) + "\n"
		chain += block.Info()
		for block.Header.Height != 0 {
			block, _ = SBC.GetParentBlock(block)
			chain += block.Info()
		}
		chains += chain
	}
	fmt.Fprint(w, chains)
}
