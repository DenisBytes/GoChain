package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/DenisBytes/GoChain/crypto"
	"github.com/DenisBytes/GoChain/node"
	"github.com/DenisBytes/GoChain/proto"
	"github.com/DenisBytes/GoChain/util"
)

// //Build list of Content to build tree
// var list []merkletree.Content
// list = append(list, TestContent{x: "Hello"})
// list = append(list, TestContent{x: "Hi"})
// list = append(list, TestContent{x: "Hey"})
// list = append(list, TestContent{x: "Hola"})

// //Create a new Merkle Tree from the list of Content
// t, err := merkletree.NewTree(list)
// if err != nil {
//   log.Fatal(err)
// }

// //Get the Merkle Root of the tree
// mr := t.MerkleRoot()
// log.Println(mr)

// //Verify the entire tree (hashes for each node) is valid
// vt, err := t.VerifyTree()
// if err != nil {
//   log.Fatal(err)
// }
// log.Println("Verify Tree: ", vt)

// //Verify a specific content in in the tree
// vc, err := t.VerifyContent(list[0])
// if err != nil {
//   log.Fatal(err)
// }

// log.Println("Verify Content: ", vc)

// //String representation
// log.Println(t)

func main() {

	makeNode(":3000", []string{}, true)
	time.Sleep(time.Second)
	makeNode(":4000", []string{":3000"}, false)
	time.Sleep(time.Second)
	makeNode(":5000", []string{":4000"}, false)
	for {
		time.Sleep(time.Millisecond * 500)
		makeTransaction()
	}
}

func makeNode(listenAddr string, bootstrapNodes []string, isValidator bool) *node.Node {
	cfg := node.ServerConfig{
		Version:    "gochain-0.1",
		ListenAddr: listenAddr,
	}
	if isValidator {
		cfg.PrivateKy = crypto.GeneratePrivateKey()
	}
	n := node.NewNode(cfg)
	go n.Start(listenAddr, bootstrapNodes)

	return n
}

func makeTransaction() {
	client, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)

	privKey := crypto.GeneratePrivateKey()
	tx := &proto.Transaction{
		Version: 1,
		Inputs: []*proto.TxInput{
			{
				PrevTxHash:   util.RandomHash(),
				PrevOutIndex: 0,
				PublicKey:    privKey.Public().Bytes(),
			},
		},
		Outputs: []*proto.TxOutput{
			{
				Amount:  99,
				Address: privKey.Public().Address().Bytes(),
			},
		},
	}

	_, err = c.HandleTransaction(context.TODO(), tx)
	if err != nil {
		log.Fatal(err)
	}
}
