package main

import (
	"context"
	"fmt"
	"time"

	"github.com/praveenkumarKajla/mocketh/client"
	"github.com/praveenkumarKajla/mocketh/models"
	"github.com/praveenkumarKajla/mocketh/subscriber"
)

func main() {
	// ctab := crontab.New()
	dbClient, err := client.NewMongoClient("mongodb://mongoadmin:secret@localhost:27888/?authSource=admin")
	if err != nil {
		fmt.Println("client.NewMongoClient Error", err)
		return
	}
	collection, err := dbClient.GetCollection("mocketh_test_db", "blocks")
	if err != nil {
		fmt.Println("dbClient.GetCollection Error", err)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var initBlock = &models.Blocks{StartBlock: "11955000"}

	result, _ := collection.InsertOne(ctx, initBlock)
	fmt.Println("initBlockLog", result)

	ethClient, err := client.NewETHClient("https://mainnet.infura.io/v3/3e411f27aa87416885b43cdc9c7456b1")
	if err != nil {
		fmt.Println("client.NewETHClient Error", err)
		return
	}
	erc20Subscriber, err := subscriber.NewErc20Subscriber(ethClient, "TestTokenSubscriber", "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174", collection)
	if err != nil {
		fmt.Println("NewErc20Subscriber Error", err)
		return
	}
	erc20Subscriber.DoRun()
}
