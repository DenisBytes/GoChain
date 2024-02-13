package types

import (
	"crypto/sha256"

	pb "google.golang.org/protobuf/proto"

	"github.com/DenisBytes/GoChain/crypto"
	"github.com/DenisBytes/GoChain/proto"
)

// Hashblock returns a SHA-256 of the header.
func HashBlock(block *proto.Block) []byte {
	b, err := pb.Marshal(block)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)

	return hash[:]
}

func SignBlock(pk *crypto.PrivateKey, b *proto.Block) *crypto.Signature {
	return pk.Sign(HashBlock(b))
}
