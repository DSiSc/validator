package tools

import (
	"crypto/sha256"
	"github.com/DSiSc/craft/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ComputeMerkleRoot(t *testing.T) {
	var data []types.Hash
	a1 := types.Hash(sha256.Sum256([]byte("a")))
	a2 := types.Hash(sha256.Sum256([]byte("b")))
	a3 := types.Hash(sha256.Sum256([]byte("c")))
	a4 := types.Hash(sha256.Sum256([]byte("d")))
	a5 := types.Hash(sha256.Sum256([]byte("e")))
	data = append(data, a1)
	data = append(data, a2)
	data = append(data, a3)
	data = append(data, a4)
	data = append(data, a5)
	hash := ComputeMerkleRoot(data)
	assert.NotEqual(t, hash, types.Hash{})
}
