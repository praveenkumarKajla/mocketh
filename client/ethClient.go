package client

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/praveenkumarKajla/mocketh/config"
)

type ETHClient struct {
	connString string
	Client     *ethclient.Client
}

var EthClient *ETHClient

func init() {
	infuraConnString := config.Config.GetString("infuraConnString")
	client, err := NewETHClient(infuraConnString)
	if err != nil {
		fmt.Println("client.NewETHClient Error", err)
		return
	}
	EthClient = client
}

func NewETHClient(connString string) (*ETHClient, error) {
	client, err := ethclient.Dial(connString)
	if err != nil {
		return nil, err
	}

	fmt.Println("we have a connection")
	return &ETHClient{connString: connString, Client: client}, nil
}

// FilterQuery between desired block range
func (ethClient *ETHClient) FilterQuery(
	startblock int,
	endblock int64,
	contractAddress common.Address) ethereum.FilterQuery {

	return ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(startblock)),
		ToBlock:   big.NewInt(int64(endblock)),
		Addresses: []common.Address{
			contractAddress,
		},
	}
}
