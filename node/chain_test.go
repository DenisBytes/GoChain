package node

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DenisBytes/GoChain/crypto"
	"github.com/DenisBytes/GoChain/proto"
	"github.com/DenisBytes/GoChain/types"
	"github.com/DenisBytes/GoChain/util"
)

func randomBlock(t *testing.T, chain *Chain) *proto.Block {
	privKey := crypto.GeneratePrivateKey()
	block := util.RandomBlock()
	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	require.Nil(t, err)
	block.Header.PrevHash = types.HashBlock(prevBlock)
	types.SignBlock(privKey, block)

	return block
}

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())
	require.Equal(t, 0, chain.Height())

	_, err := chain.GetBlockByHeight(0)
	require.Nil(t, err)
}

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())
	for i := 1; i < 100; i++ {
		block := randomBlock(t, chain)
		blockHash := types.HashBlock(block)

		require.Nil(t, chain.AddBlock(block))

		fetchedBlock, err := chain.GetBlockByHash(blockHash)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlockByHeight)
	}
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())
	for i := 1; i < 100; i++ {
		b := randomBlock(t, chain)

		require.Nil(t, chain.AddBlock(b))
		require.Equal(t, chain.Height(), i)
	}
}
