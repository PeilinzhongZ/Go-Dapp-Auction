package data

type HeartBeatData struct {
	IfNewBlock  bool   `json:"ifNewBlock"`
	Id          int32  `json:"id"`
	BlockJson   string `json:"blockJson"`
	PeerMapJson string `json:"peerMapJson"`
	Addr        string `json:"addr"`
	Hops        int32  `json:"hops"`
	IfBid       bool   `json:"ifBid`
	BidJson     string `json:"bidJson`
}

func NewHeartBeatData(ifNewBlock bool, id int32, blockJson string, peerMapJson string, addr string) HeartBeatData {
	return HeartBeatData{IfNewBlock: ifNewBlock, Id: id, BlockJson: blockJson, PeerMapJson: peerMapJson, Addr: addr}
}

func PrepareHeartBeatData(sbc *SyncBlockChain, selfId int32, peerMapJson string, addr string) HeartBeatData {
	heartBeatData := NewHeartBeatData(false, selfId, "", peerMapJson, addr)
	heartBeatData.Hops = 1
	return heartBeatData
}

func PrepareBidData(selfId int32, peerMapJson string, addr string, bidJson string) HeartBeatData {
	heartBeatData := HeartBeatData{Id: selfId, PeerMapJson: peerMapJson, Addr: addr, IfBid: true, BidJson: bidJson}
	heartBeatData.Hops = 1
	return heartBeatData
}
