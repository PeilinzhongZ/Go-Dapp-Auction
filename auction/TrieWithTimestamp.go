package auction

import "../p1"

type TrieWithTime struct {
	Trie      p1.MerklePatriciaTrie
	Timestamp int64
}
