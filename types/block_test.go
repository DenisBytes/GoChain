package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DenisBytes/GoChain/crypto"
	"github.com/DenisBytes/GoChain/proto"
	"github.com/DenisBytes/GoChain/util"
)

func TestHashBlock(t *testing.T) {
	block := util.RandomBlock()
	hash := HashBlock(block)
	assert.Equal(t, 32, len(hash))
}

func TestSignVerifyBlock(t *testing.T) {
	var (
		block   = util.RandomBlock()
		privKey = crypto.GeneratePrivateKey()
		pubKey  = privKey.Public()
	)

	sig := SignBlock(privKey, block)
	assert.Equal(t, 64, len(sig.Bytes()))
	assert.True(t, sig.Verify(pubKey, HashBlock(block)))

	assert.Equal(t, block.PublicKey, pubKey.Bytes())
	assert.Equal(t, block.Signature, sig.Bytes())

	assert.True(t, VerifyBlock(block))

	invalidPRivKey := crypto.GeneratePrivateKey()
	block.PublicKey = invalidPRivKey.Public().Bytes()

	assert.False(t, VerifyBlock(block))
}

func TestCalculateRootHash(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	block := util.RandomBlock()
	tx := &proto.Transaction{
		Version: 1,
	}
	block.Transactions = append(block.Transactions, tx)
	SignBlock(privKey, block)

	assert.True(t, VerifyRootHash(block))
	assert.Equal(t, 32, len(block.Header.RootHash))
}
