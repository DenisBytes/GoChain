package types

import (
	"crypto/sha256"

	"github.com/DenisBytes/GoChain/crypto"
	"github.com/DenisBytes/GoChain/proto"
	pb "google.golang.org/protobuf/proto"
)

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}

func SignTransaction(pk *crypto.PrivateKey, tx *proto.Transaction) *crypto.Signature {
	return pk.Sign(HashTransaction(tx))
}

func VerifyTransaction(tx *proto.Transaction) bool {
	for _, input := range tx.Inputs {
		if len(input.Signature) == 0 {
			panic("the transaction has no signature")
		}
		sig := crypto.SignatureFromBytes(input.Signature)
		pubKey := crypto.PublicKeyFromBytes(input.PublicKey)
		input.Signature = nil
		if !sig.Verify(pubKey, HashTransaction(tx)) {
			return false
		}
	}
	return true
}
