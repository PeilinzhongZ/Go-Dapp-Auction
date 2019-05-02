# Auction Marketplace

## What?
Using Blockchain to create a marketplace for auction

## Why?
Auctioneer may lie about which bidder provide hightest price for this auction. When using blockchain to implement a marketplace for auction, every bidder have to put transaction in the blockchain, which can provide immutability, and every can see the transaction whcih can confirm the actual winner.

## How?
- Auctioneer can post the information about the item to blockchain, such as description and starting price, etc. Also Auctioneer also can based on transactions to determine the winner of this auction.
- Bidder can post their price to miner. Also they can review the recent transcation of the auction and check the auctioneer determine winner correctlly.
- Miner would create block with transaction based on the price miner posted. And Miner can get reward from the each payment and service fee from bidder.

## Functionalities
### Auctioneer
1. Post item info such as description, phote, end time and starting price, etc. to the blockchain.
    > Midpoint--before 04/22
2. Determine the winner and post the final result of auction, winner info, to the blockchain. (Determine when 6 blocks are created after end time)

### Bidder
1. Find all auction available (before end time) in blockchain.
    > Midpoint--before 04/26
2. Post price for specific auction to miner.
    > Midpoint--before 05/01
3. Check specific auction transaction
4. Check if the winner of specific auction is valid.

### Miner
1. Creat block based on provided price of bidder.

## Success
- Making sure that every bidder can censor the auction and checking if the winner transaction is in canonical chain and this winner provide hightest price.

## Data structure:
```
type Auctioneer struct {  // represent the Info of Auctioneer
	ID      int       //ID of Auctioneer (Port Number)
	Address string    //IP address and Port Number
	ItemNum int       // Number of Item posted
}
```
- ItemNum would automatically increase by 1 every time post a Item in blockchain
```
type Bidder struct {    // represent the Info of Bidder
	ID      int     // ID of Bidder (Port Number)
	Address string  // IP address and Port Number
}
```

```
type ItemInfo struct {
	Name        string  // Name of Item
	Description string  // Description of Item
	Price       int     // Starting price of Item
	End         int64   // End Time for auction (Unix Timestamp)
}
```
- ItemInfo used for data send by auctioneer to API
```
type Item struct {
	AuctioneerID int       // Auctioneer ID
	ItemID       int       // Item ID
	Info         ItemInfo  //ItemInfo
}
```
- Item is used for parsing Item in mpt, AuctioneerID and ItemID response for the unique Item ID in blockchain. (ItemID are the number of Item this auctioneer already post)

```
type BidInfo struct {
	AuctioneerID int  // represent the Item belong to which auctioneer
	ItemID       int  // represent the Item ID of specific auctioneer
	Bid          int  // the price provide for this bid
}
```
- BidInfo used for parsing data send by Bidder to API
- Combining AutioneerID and ItemID would represent the ID of this Item
```
type BidDetail struct {
	BidderID int      // Bidder ID bid on this Item
	BidInfo  BidInfo  // Info of this bid
}
```
- BidDetail used for sending Bid to Miner
```
type Transaction struct {
	MinerID int        // ID of miner
	Detail  BidDetail  // Detail of this bid
}
```
- Transaction is used for parsing Bid in mpt
```
type ItemData struct {
	Iteminfo Item  // Information of this Item
	Trans    []Transaction  // all bids for this specific Item
}
```
- ItemData used to store all bids for specific Item
- ItemData would be store in map access with this ItemID (string of combination of AuctionID and ItemID)

## Implementation: (Accomplished by Midpoint)
### Auction:
1. Post item info to blockchain.
    - Parsing the body of POST request into ItemInfo struct then insert each element in ItemInfo into MPT
    - Using this MPT and start trying Nonce
### Bidder
1. Find all auction available in blockchain
    - Get and store MPT in canonical blocks
    - Parse MPT form the mpt of first block, if transaction Type is "ItemInfo", parse this mpt into ItemInfo, then create ItemData for this ItemInfo and store in a map. Otherwise, parse this mpt into Transaction, then add to list of transaction for specific ItemInfo in the map.
    - Marshal the map into json data, reply to user

2. Post price for specific auction to miner
    - Parse the post request send by Bidder into BidInfo.
    - Create BidDetail by add parameter BidderID (Port number), then send BidDetail to Node on PeerList.

### Each Node
- Receive the bid send by Bidder, and Forward to others in its Peerlist

### Changes in previous work:
1. After the server start, it would create corresponding auctioneer and bidder instance based on Port Number and IP address
2. Alter HeartBeatData to support BidDetail.
3. Alter ForwardHeartBeat and ReceiveHeartBeat to support BidDetail 