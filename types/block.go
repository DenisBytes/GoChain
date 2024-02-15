package types

import (
	"bytes"
	"crypto/sha256"

	pb "google.golang.org/protobuf/proto"

	"github.com/DenisBytes/GoChain/crypto"
	"github.com/DenisBytes/GoChain/proto"
	"github.com/cbergoon/merkletree"
)

type TxHash struct {
	hash []byte
}

func NewTxHash(hash []byte) TxHash {
	return TxHash{
		hash: hash,
	}
}

func (h TxHash) CalculateHash() ([]byte, error) {
	return h.hash, nil
}

func (h TxHash) Equals(other merkletree.Content) (bool, error) {
	equals := bytes.Equal(h.hash, other.(TxHash).hash)
	return equals, nil
}

// Hashblock returns a SHA-256 of the header.
func HashBlock(block *proto.Block) []byte {
	return HashHeader(block.Header)
}

func HashHeader(header *proto.Header) []byte {
	b, err := pb.Marshal(header)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)

	return hash[:]
}

func SignBlock(pk *crypto.PrivateKey, b *proto.Block) *crypto.Signature {

	if len(b.Transactions) > 0 {
		tree, err := GetMerleTree(b)
		if err != nil {
			panic(err)
		}
		b.Header.RootHash = tree.MerkleRoot()
	}

	hash := HashBlock(b)
	sig := pk.Sign(hash)

	b.PublicKey = pk.Public().Bytes()
	b.Signature = sig.Bytes()

	return sig
}

func VerifyBlock(b *proto.Block) bool {
	if len(b.Transactions) > 0 {
		if !VerifyRootHash(b) {
			return false
		}
	}
	if len(b.PublicKey) != crypto.PubKeyLen {
		return false
	}
	if len(b.Signature) != crypto.SignatureLen {
		return false
	}
	sig := crypto.SignatureFromBytes(b.Signature)
	pubKey := crypto.PublicKeyFromBytes(b.PublicKey)
	hash := HashBlock(b)

	return sig.Verify(pubKey, hash)
}

func VerifyRootHash(b *proto.Block) bool {
	tree, err := GetMerleTree(b)
	if err != nil {
		return false
	}
	valid, err := tree.VerifyTree()
	if err != nil {
		return false
	}

	if !valid {
		return false
	}

	return bytes.Equal(b.Header.RootHash, tree.MerkleRoot())
}

func GetMerleTree(b *proto.Block) (*merkletree.MerkleTree, error) {

	list := make([]merkletree.Content, len(b.Transactions))
	for i := 0; i < len(b.Transactions); i++ {
		list[i] = NewTxHash(HashTransaction(b.Transactions[i]))
	}

	t, err := merkletree.NewTree(list)
	if err != nil {
		return nil, err
	}

	return t, nil
}
