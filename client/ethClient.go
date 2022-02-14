package client

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/praveenkumarKajla/mocketh/config"
	"github.com/praveenkumarKajla/mocketh/models"
	"github.com/sirupsen/logrus"
)

type ETHClient struct {
	connString string
	Client     *ethclient.Client
	WssClient  *ethclient.Client
}

var EthClient *ETHClient

func init() {
	infuraConnString := config.Config.GetString("infuraConnString")
	infuraConnStringWss := config.Config.GetString("infuraConnStringWss")

	client, err := NewETHClient(infuraConnString, infuraConnStringWss)
	if err != nil {
		logrus.Info("client.NewETHClient Error", err)
		return
	}
	EthClient = client
}

func NewETHClient(connString string, wssConnString string) (*ETHClient, error) {
	client, err := ethclient.Dial(connString)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	// wss client is for notifying through channel whenever new Transer event is fired
	wssClient, err := ethclient.Dial(wssConnString)
	if err != nil {
		return nil, err
	}

	logrus.Info("we have a connection")
	return &ETHClient{connString: connString, Client: client, WssClient: wssClient}, nil
}

// FilterQuery between desired block range
func (ethClient *ETHClient) FilterQuery(
	startblock *big.Int,
	endblock *big.Int,
	contractAddress common.Address) ethereum.FilterQuery {

	return ethereum.FilterQuery{
		FromBlock: startblock,
		ToBlock:   endblock,
		Addresses: []common.Address{
			contractAddress,
		},
	}
}

// returns ERC20 struct to store in db
func GetERC20(ctx context.Context, addr common.Address, name string, fromBlock int64) (*models.ERC20, error) {
	erc20 := &models.ERC20{
		Address:          addr.Bytes(),
		Name:             name,
		LastIndexedBlock: fromBlock,
	}
	return erc20, nil
}
