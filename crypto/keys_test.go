package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	assert.Equal(t, len(privKey.Bytes()), privKeyLen)

	pubKey := privKey.Public()
	assert.Equal(t, len(pubKey.Bytes()), pubKeyLen)
}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		seed       = "4c179dfc5d3b9f3c78e6b0a3a64e280a4c09947958959a150b59dd3023c02bae"
		privKey    = NewPrivateKeyFromString(seed)
		addressStr = "75b7d7c7776c5bcbaa30ac2fac0272806e5185b7"
	)

	assert.Equal(t, privKeyLen, len(privKey.Bytes()))

	address := privKey.Public().Address()
	assert.Equal(t, addressStr, address.String())
	fmt.Println(address)
}

func TestPrivateKeySign(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	msg := []byte("foo bar baz")

	sig := privKey.Sign(msg)
	assert.True(t, sig.Verify(pubKey, msg))

	// Test with invalid message
	assert.False(t, sig.Verify(pubKey, []byte("foo")))

	// Test with other pubkey
	otherPrivKey := GeneratePrivateKey()
	otherPubKey := otherPrivKey.Public()
	assert.False(t, sig.Verify(otherPubKey, msg))
}

func TestPublicKetToAddress(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()
	assert.Equal(t, addressLen, len(address.Bytes()))
}
