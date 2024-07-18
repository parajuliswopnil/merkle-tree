package merkle

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Node struct {
	Hash        []byte
	SiblingHash []byte
	Parent      *Node
	Index       int
}

func Hasher(data []byte) []byte {
	return crypto.Keccak256(data)
}

func MakeNodes(data [][]byte) []*Node {
	nodeList := []*Node{}
	for i, v := range data {
		nodeList = append(nodeList, &Node{
			Hash:  Hasher(v),
			Index: i,
		})
	}
	return nodeList
}

func PrintHashes(list []*Node) {
	for _, v := range list {
		fmt.Printf(hexutil.Encode(v.Hash) + " ")
	}
	fmt.Println()
}

func MakeMerkleTree(leafNodes []*Node) *Node {
	if len(leafNodes) == 1 {
		return leafNodes[0]
	}

	if len(leafNodes)%2 != 0 {
		leafNodes = append(leafNodes, leafNodes[len(leafNodes)-1])
	}
	newLeaf := []*Node{}
	for i := 0; i < len(leafNodes); i += 2 {
		concatHash := make([]byte, len(leafNodes[i].Hash))
		copy(concatHash, leafNodes[i].Hash)
		combinedHash := append(concatHash, leafNodes[i+1].Hash...)
		leafNodes[i].SiblingHash = leafNodes[i+1].Hash
		leafNodes[i+1].SiblingHash = leafNodes[i].Hash
		parentOfIandI1 := &Node{
			Hash:  Hasher(combinedHash),
			Index: len(newLeaf),
		}
		leafNodes[i].Parent = parentOfIandI1
		leafNodes[i+1].Parent = parentOfIandI1
		newLeaf = append(newLeaf, parentOfIandI1)
	}
	return MakeMerkleTree(newLeaf)
}

type Proof struct {
	Sibling  []byte
	Position int
}

func CalculateProof(leafNode *Node) []*Proof {
	proof := []*Proof{}
	for leafNode.Parent != nil {
		sibling := leafNode.SiblingHash
		var position int
		if leafNode.Index%2 == 0 {
			position = 1
		}
		proof = append(proof, &Proof{
			Sibling:  sibling,
			Position: position,
		})
		leafNode = leafNode.Parent
	}
	return proof
}

func VerifyProof(leaf, root []byte, proof []*Proof) bool {
	combinedHash := make([]byte, len(leaf))
	copy(combinedHash, leaf)

	for i := 0; i < len(proof); i++ {
		if proof[i].Position == 0 {
			combinedHash = Hasher(append(proof[i].Sibling, combinedHash...))
		} else {
			combinedHash = Hasher(append(combinedHash, proof[i].Sibling...))
		}
	}
	return bytes.Equal(combinedHash, root)
}

