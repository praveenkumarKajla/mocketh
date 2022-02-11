package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mileusna/crontab"
	"github.com/praveenkumarKajla/mocketh/client"
	"github.com/praveenkumarKajla/mocketh/models"
	"github.com/praveenkumarKajla/mocketh/subscriber"
)

func main() {
	ctab := crontab.New()

	collection, err := client.DBClient.GetCollection("mocketh_test_db", "blocks")
	if err != nil {
		fmt.Println("dbClient.GetCollection Error", err)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var initBlock = &models.Blocks{StartBlock: "11955000"}
	result, err := collection.InsertOne(ctx, initBlock)
	if err != nil {
		fmt.Println("collection.InsertOne", err)
		return
	}
	fmt.Println("initBlockLog", result)

	erc20Subscriber, err := subscriber.NewErc20Subscriber(client.EthClient, "TestTokenSubscriber", "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174", collection)
	if err != nil {
		fmt.Println("NewErc20Subscriber Error", err)
		return
	}
	ctab.MustAddJob("* * * * *", erc20Subscriber.DoRun)
}
