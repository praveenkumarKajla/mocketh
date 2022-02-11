package subscriber

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/praveenkumarKajla/mocketh/client"
	"github.com/praveenkumarKajla/mocketh/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Erc20Subscriber struct {
	LastRunAt       time.Time
	Name            string
	ContractAddress common.Address
	ethclient       *client.ETHClient
	collection      *mongo.Collection
}

func NewErc20Subscriber(ethClient *client.ETHClient, Name string, HexAddress string, collection *mongo.Collection) (*Erc20Subscriber, error) {
	contractAddress := common.HexToAddress("0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174")

	return &Erc20Subscriber{
		LastRunAt:       time.Now(),
		Name:            Name,
		ContractAddress: contractAddress,
		ethclient:       ethClient,
		collection:      collection,
	}, nil

}

func (Subscriber *Erc20Subscriber) DoRun() {
	ethclient := Subscriber.ethclient.Client
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var block models.Blocks
	if err := Subscriber.collection.FindOne(ctx, bson.M{}).Decode(&block); err != nil {
		log.Fatal(err)
	}
	startblock, _ := strconv.Atoi(block.StartBlock)
	endblock, err := ethclient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	endblockInt := endblock.Number.Int64()
	fmt.Println(startblock, endblockInt)

}
