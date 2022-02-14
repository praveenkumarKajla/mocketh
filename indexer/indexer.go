package indexer

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/praveenkumarKajla/mocketh/client"
	"github.com/praveenkumarKajla/mocketh/config"
	"github.com/praveenkumarKajla/mocketh/models"
	"github.com/praveenkumarKajla/mocketh/store"
	"github.com/praveenkumarKajla/mocketh/subscriber"
	"github.com/sirupsen/logrus"
)

const (
	IntervalTime = 10
)

var (
	//ErrInvalidAddress returns if invalid ERC20 address is detected
	ErrInvalidAddress = errors.New("invalid address")
	onlyOnce          sync.Once
)

type Indexer struct {
	Account   *store.Account
	ethclient *client.ETHClient
	tokenList map[common.Address]*models.ERC20
}

func New() (*Indexer, error) {
	accountcollection, err := client.DBClient.GetCollection("erc20")
	if err != nil {
		return nil, err
	}
	account := store.NewWithCollection(accountcollection)
	ethclient := client.EthClient
	return &Indexer{
		Account:   account,
		ethclient: ethclient,
	}, nil
}

func (idx *Indexer) Erc20TokensSubscriber(ctx context.Context, addresses map[string]string) error {
	logrus.Info("Address ", addresses)
	var fromBlock int64
	onlyOnce.Do(func() {
		// get the starting block from config
		// or else the indexing will start from latest block backwards
		fromBlockStr := config.Config.GetString("fromBlock")
		if fromBlockStr == "" {
			logrus.Fatal("invalid fromBlock config value")
		}
		val, err := strconv.ParseInt(fromBlockStr, 10, 64)
		if err != nil {
			// handle error
			logrus.Info(err)
		}
		fromBlock = val
	})
	// check if contract address in config exist in db
	for name, addr := range addresses {
		err := idx.GetOrAddErc20(ctx, addr, name)
		if err != nil {
			logrus.Error(err)
		}

	}
	idx.Init(ctx, fromBlock)
	return nil
}

func (idx *Indexer) Listen() error {
	err := idx.DoRun()
	if err != nil {
		logrus.Fatal("Error Running Periodic Indexer")
		return err
	}
	// to run task periodically every IntervalTime seconds
	ticker := time.NewTicker(IntervalTime * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err := idx.DoRun()
				if err != nil {
					logrus.Fatal("Error Running Periodic Indexer")
					ticker.Stop()
					return
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return nil
}

func (idx *Indexer) DoRun() error {
	sigs := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
		defer signal.Stop(sigs)

		logrus.Debug("Shutting down", "signal", <-sigs)
		cancel()
	}()

	// iterate over all ERC20 addresses to get transfer logs
	list, err := idx.Account.ListOldERC20(ctx)
	if err != nil {
		logrus.Info("Init err", err)
		return err
	}
	for _, Erc20Token := range list {
		collection, err := client.DBClient.GetCollection("transfers")
		if err != nil {
			logrus.Info("dbClient.GetCollection Error", err)
			return err
		}
		//  Subscriber at ERC20 level, add tasks related to ERC20
		erc20Subscriber, err := subscriber.NewErc20Subscriber(client.EthClient, Erc20Token.Name, common.BytesToAddress(Erc20Token.Address), collection)
		if err != nil {
			logrus.Info("NewErc20Subscriber Error", err)
			return err
		}
		// to add the upcoming transfer events in db
		go erc20Subscriber.UpcomingEvents()
		//   <-----erc20Subscriber.DoRun(Erc20Token) *****NOW*****   erc20Subscriber.UpcomingEvents() --------->
		// task to store past event logs to DB
		updatedToken, err := erc20Subscriber.DoRun(Erc20Token)
		if err != nil {
			return err
		}
		err = idx.Account.UpdateERC20Block(ctx, updatedToken)
		if err != nil {
			return err
		}

	}
	return nil
}

func (idx *Indexer) GetOrAddErc20(ctx context.Context, addr string, name string) error {
	if !common.IsHexAddress(addr) {
		return ErrInvalidAddress
	}
	address := common.HexToAddress(addr)
	val, err := idx.Account.FindERC20(ctx, address)
	// The ERC20 exists, no need to insert again
	if err == nil {
		logrus.Info("Erc20 Exist ", val)
		return errors.New("erc20 Already Exist")
	}

	erc20, err := client.GetERC20(ctx, address, name, 0)
	if err != nil {
		logrus.Error("Failed to get ERC20", "addr", addr, "err", err)
		return err
	}
	//  add the token to DB
	_, err = idx.Account.InsertERC20(ctx, erc20)
	if err != nil {
		logrus.Error("Failed to insert ERC20", "addr", addr, "err", err)
		return err
	}
	return nil
}

// get updated token list & store in memory
// not required for now
func (idx *Indexer) Init(ctx context.Context, fromBlock int64) error {
	list, err := idx.Account.ListOldERC20(ctx)
	if err != nil {
		logrus.Info("Init err", err)
		return err
	}

	tokenList := make(map[common.Address]*models.ERC20, len(list))
	tokenList[models.ETHAddress] = &models.ERC20{
		Address:          models.ETHBytes,
		LastIndexedBlock: fromBlock,
	}
	for _, e := range list {
		logrus.Info(e)
		tokenList[common.BytesToAddress(e.Address)] = e
	}
	idx.tokenList = tokenList
	return nil
}
