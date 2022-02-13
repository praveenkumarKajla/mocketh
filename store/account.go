package store

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/praveenkumarKajla/mocketh/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Account struct {
	Collection *mongo.Collection
}

func NewWithCollection(collection *mongo.Collection) *Account {
	return &Account{
		Collection: collection,
	}
}

//  Get all the ERC20 tokens
func (acc *Account) ListOldERC20(ctx context.Context) ([]*models.ERC20, error) {
	var erc20s []*models.ERC20
	cursor, err := acc.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &erc20s); err != nil {
		return nil, err
	}
	logrus.Info(erc20s)
	return erc20s, nil
}

//  Add new ERC20
func (acc *Account) InsertERC20(ctx context.Context, token *models.ERC20) (id *mongo.InsertOneResult, err error) {
	result, err := acc.Collection.InsertOne(ctx, token)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (acc *Account) FindERC20(ctx context.Context, address common.Address) (*models.ERC20, error) {
	var erc20 models.ERC20
	addressBytes := address.Bytes()
	if err := acc.Collection.FindOne(ctx, bson.M{"address": addressBytes}).Decode(&erc20); err != nil {
		return nil, err
	}
	logrus.Info(erc20.LastIndexedBlock)
	return &erc20, nil

}

// Update the last block indexed
func (acc *Account) UpdateERC20Block(ctx context.Context, erc20 *models.ERC20) error {
	filter := bson.M{"_id": bson.M{"$eq": erc20.ID}}
	update := bson.M{"$set": bson.M{"block_number": erc20.LastIndexedBlock}}
	if _, err := acc.Collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}
