package subscriber

import (
	"context"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/praveenkumarKajla/mocketh/client"
	"github.com/praveenkumarKajla/mocketh/models"
	contract "github.com/praveenkumarKajla/mocketh/token"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// No. of blocks to poll logs at a time
const (
	blockIntervals = 100
)

type Erc20Subscriber struct {
	LastRunAt       time.Time
	Name            string
	ContractAddress common.Address
	ethclient       *client.ETHClient
	collection      *mongo.Collection
}

func NewErc20Subscriber(
	ethClient *client.ETHClient,
	Name string,
	contractAddress common.Address,
	collection *mongo.Collection,
) (*Erc20Subscriber, error) {

	return &Erc20Subscriber{
		LastRunAt:       time.Now(),
		Name:            Name,
		ContractAddress: contractAddress,
		ethclient:       ethClient,
		collection:      collection,
	}, nil

}

func (Subscriber *Erc20Subscriber) DoRun(Erc20Token *models.ERC20) (*models.ERC20, error) {
	logrus.Info("Running Subscriber")
	ethclient := Subscriber.ethclient
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var endBlock *big.Int
	// if indexing this token for first time get the latest block
	if Erc20Token.LastIndexedBlock != 0 {
		endBlock = big.NewInt(Erc20Token.LastIndexedBlock)
	} else {
		lastBlock, err := Subscriber.ethclient.Client.BlockByNumber(ctx, nil)
		if err != nil {
			logrus.Fatal(err)
		}
		endBlock = lastBlock.Number()
	}

	startblock := big.NewInt(0).Sub(endBlock, big.NewInt(blockIntervals))
	logrus.Info(startblock, endBlock)
	// generate eth filterQuery to get logs between range of blocks

	query := ethclient.FilterQuery(startblock, endBlock, Subscriber.ContractAddress)
	// Query logs by filter Query
	logs, err := ethclient.Client.FilterLogs(ctx, query)
	if err != nil {
		logrus.Fatal(err)
	}
	// parse the JSON abi
	contractAbi, err := abi.JSON(strings.NewReader(string(contract.ContractABI)))
	if err != nil {
		logrus.Fatal(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")

	// keccak256 hash of each event log function signature
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	for _, vLog := range logs {

		switch vLog.Topics[0].Hex() {
		// if need to get approval logs too add the approval sig hash
		case logTransferSigHash.Hex():
			logrus.Infof("Log Index: %d\n", vLog.Index)
			logrus.Infof("Log Data: %d\n", vLog.Data)
			logrus.Infof("Log Name: Transfer\n")
			var transferEvent models.LogTransfer

			err = contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			logrus.Info("transferEventData : ", vLog.Data)
			if err != nil {
				logrus.Fatal(err)
			}
			for _, topic := range vLog.Topics {
				logrus.Info(topic, topic.Hex())
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			event := &models.Erc20TransferEvent{
				From:        transferEvent.From.Hex(),
				To:          transferEvent.To.Hex(),
				Tokens:      transferEvent.Value.String(),
				BlockNumber: vLog.BlockNumber,
				TxHash:      vLog.TxHash.Hex(),
			}

			result, err := Subscriber.collection.InsertOne(ctx, event)
			if err != nil {
				logrus.Fatal(err)
			}

			logrus.Info("dbID : event ", result)
			// logrus.Infof("\n\n")

		}

	}
	Erc20Token.LastIndexedBlock = endBlock.Int64()
	return Erc20Token, nil
}

// handles the new transfer event fired and store them to db
func (Subscriber *Erc20Subscriber) UpcomingEvents() (*models.ERC20, error) {
	tc, err := contract.NewContract(Subscriber.ContractAddress, Subscriber.ethclient.WssClient)
	if err != nil {
		logrus.Fatalf("error connections to contract: %s", err.Error())
	}

	// Create transfers channel.
	transfers := make(chan *contract.ContractTransfer)

	// Subscribe to ethereum contract transfers event.
	sub, err := tc.WatchTransfer(nil, transfers, nil, nil)
	if err != nil {
		logrus.Fatalf("error subscribing to event: %s", err.Error())
	}

	for {
		select {
		case err := <-sub.Err():
			logrus.Fatalf(": %s", err.Error())
		case t := <-transfers:
			event := &models.Erc20TransferEvent{
				From:        t.From.Hex(),
				To:          t.To.Hex(),
				Tokens:      t.Value.String(),
				BlockNumber: t.Raw.BlockNumber,
				TxHash:      t.Raw.TxHash.Hex(),
			}

			log.Println(event)
			result, err := Subscriber.collection.InsertOne(context.Background(), event)
			if err != nil {
				logrus.Fatal(err)
			}

			logrus.Info(result)
		}
	}

}
