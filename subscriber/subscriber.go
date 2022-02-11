package subscriber

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/praveenkumarKajla/mocketh/client"
	"github.com/praveenkumarKajla/mocketh/models"
	"github.com/praveenkumarKajla/mocketh/token"
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

	// generate eth filterQuery to get logs between range of blocks
	query := Subscriber.ethclient.FilterQuery(startblock, endblockInt, Subscriber.ContractAddress)

	// Query logs by filter Query
	logs, err := ethclient.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	// parse the JSON abi
	contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenABI)))
	if err != nil {
		log.Fatal(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")

	// keccak256 hash of each event log function signature
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	for _, vLog := range logs {

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
			fmt.Printf("Log Index: %d\n", vLog.Index)
			fmt.Printf("Log Data: %d\n", vLog.Data)
			fmt.Printf("Log Name: Transfer\n")

			var transferEvent models.LogTransfer

			err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			// val, err := contractAbi.Unpack("Transfer", vLog.Data)

			if err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			event := &models.Erc20TransferEvent{
				From:   transferEvent.From.Hex(),
				To:     transferEvent.To.Hex(),
				Tokens: transferEvent.Tokens.String(),
			}

			result, _ := Subscriber.collection.InsertOne(ctx, event)

			fmt.Printf("From: %s\n", transferEvent.From.Hex())
			fmt.Printf("To: %s\n", transferEvent.To.Hex())
			fmt.Printf("Tokens: %s\n", transferEvent.Tokens.String())
			fmt.Println("dbID ", result)
			fmt.Printf("\n\n")

		}

	}
}
