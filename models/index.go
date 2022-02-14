package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// ETHAddress represents ether type in address type
	ETHAddress = common.BytesToAddress([]byte("ETH"))
	// ETHBytes represents ether type in bytes array type
	ETHBytes = ETHAddress.Bytes()
)

// ERC20 represents the ERC20 contract
type ERC20 struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name             string             `json:"name,omitempty" bson:"name,omitempty"`
	Address          []byte             `json:"address,omitempty" bson:"address,omitempty"`
	LastIndexedBlock int64              `json:"block_number,omitempty" bson:"block_number,omitempty"`
}

// Database level modelling of ERC20 token event
// Mongo DB entries
type Erc20TransferEvent struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	From          string             `json:"from,omitempty" bson:"from,omitempty"`
	To            string             `json:"to,omitempty" bson:"to,omitempty"`
	Tokens        string             `json:"tokens,omitempty" bson:"tokens,omitempty"`
	BlockNumber   uint64             `json:"block_number,omitempty" bson:"block_number,omitempty"`
	TxHash        string             `json:"tx_hash,omitempty" bson:"tx_hash,omitempty"`
	Address       string             `json:"address,omitempty" bson:"address,omitempty"`
	BlockLogIndex uint64             `json:"block_index,omitempty" bson:"block_index,omitempty"`
}

//struct types matching the types of the ERC-20 event Transfer
type LogTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}
