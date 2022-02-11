package client

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ETHClient struct {
	connString string
	Client     *ethclient.Client
}

func NewETHClient(connString string) (*ETHClient, error) {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/3e411f27aa87416885b43cdc9c7456b1")
	if err != nil {
		return nil, err
	}

	fmt.Println("we have a connection")
	return &ETHClient{connString: connString, Client: client}, nil
}

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
