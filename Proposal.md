# Auction Marketplace

## What?
Using Blockchain to create a marketplace for auction

## Why?
Auctioneer may lie about which bidder provide hightest price for this auction. When using blockchain to implement a marketplace for auction, every bidder have to put transaction in the blockchain, which can provide immutability, and every can see the transaction whcih can confirm the actual winner.

## How?
- Auctioneer can post the information about the item to blockchain, such as description and starting price, etc. Also Auctioneer also can based on transactions to determine the winner of this auction.
- Bidder can post their price to minner. Also they can review the recent transcation of the auction and check the auctioneer determine winner correctlly.
- Minner would create block with transaction based on the price minner posted. And Minner can get reward from the each payment and service fee from bidder.

## Functionalities
### Auctioneer
1. Post item info such as description, phote, end time and starting price, etc. to the blockchain.
    > Midpoint--before 04/22
2. Determine the winner and post the final result of auction, winner info, to the blockchain.

### Bidder
1. Find all auction available (before end time) in blockchain.
    > Midpoint--before 04/26
2. Post price for specific auction to minner. 
    > Midpoint--before 05/01
3. Check specific auction transaction
4. Check if the winner of specific auction is valid.

### Minner
1. Creat block based on provided price of bidder.

## Success
Making sure that every bidder can censor the auction and checking if the winner transaction is in canonical chain and this winner provide hightest price.
