package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerkle(t *testing.T) {
	data := [][]byte{
		[]byte("hello1"),
		[]byte("hello2"),
		[]byte("hello3"),
		[]byte("hello4"),
		[]byte("hello5"),
		[]byte("hello6"),
		[]byte("hello7"),
	}
	leafNodes := MakeNodes(data)
	root := MakeMerkleTree(leafNodes)
	proof := CalculateProof(leafNodes[6])
	assert.True(t, VerifyProof(leafNodes[6].Hash, root.Hash, proof))
}
