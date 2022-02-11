package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Database level modelling of ERC20 token event
// Mongo DB entries
type Erc20Event struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	From       string             `json:"from,omitempty" bson:"from,omitempty"`
	To         string             `json:"to,omitempty" bson:"to,omitempty"`
	Tokens     string             `json:"tokens,omitempty" bson:"tokens,omitempty"`
	TokenOwner string             `json:"tokenOwner,omitempty" bson:"tokenOwner,omitempty"`
	Spender    string             `json:"spender,omitempty" bson:"spender,omitempty"`
}

type Blocks struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	StartBlock string             `json:"startblock,omitempty" bson:"from,omitempty"`
}

type TransferEvent struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
