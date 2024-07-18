package main

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

func makeNodes(data [][]byte) []*Node {
	nodeList := []*Node{}
	for i, v := range data {
		nodeList = append(nodeList, &Node{
			Hash:  Hasher(v),
			Index: i,
		})
	}
	return nodeList
}

func printHashes(list []*Node) {
	for _, v := range list {
		fmt.Printf(hexutil.Encode(v.Hash) + " ")
	}
	fmt.Println("\n")
}

func makeMerkleTree(leafNodes []*Node) *Node {
	printHashes(leafNodes)

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
	return makeMerkleTree(newLeaf)
}

type Proof struct {
	Sibling []byte
	Position int
}

func calculateProof(leafNode *Node) []*Proof {
	proof := []*Proof{}
	for leafNode.Parent != nil {
		sibling := leafNode.SiblingHash
		var position int
		if leafNode.Index % 2 == 0 {
			position = 1
		} 
		proof = append(proof, &Proof{
			Sibling: sibling,
			Position: position,
		})
		leafNode = leafNode.Parent
	}
	return proof
}

func verifyProof(leaf, root []byte, proof []*Proof) bool {
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

func main() {
	data := [][]byte{
		[]byte("hello1"),
		[]byte("hello2"),
		[]byte("hello3"),
		[]byte("hello4"),
		[]byte("hello5"),
		[]byte("hello6"),
		[]byte("hello7"),
	}
	leafNodes := makeNodes(data)
	root := makeMerkleTree(leafNodes)

	fmt.Println(hexutil.Encode(root.Hash))

	proof := calculateProof(leafNodes[6])

	for _, v := range proof {
		fmt.Println(hexutil.Encode(v.Sibling))
	}

	fmt.Println(verifyProof(leafNodes[6].Hash, root.Hash, proof))
}


