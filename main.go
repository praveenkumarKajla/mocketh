package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/praveenkumarKajla/mocketh/indexer"
	"github.com/praveenkumarKajla/mocketh/models"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// Add Logrus formatting
func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	ll, err := logrus.ParseLevel("debug")
	if err != nil {
		ll = logrus.DebugLevel
	}
	logrus.SetLevel(ll)
}

func main() {
	sigs := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
		defer signal.Stop(sigs)

		log.Debug("Shutting down", "signal", <-sigs)
		cancel()
	}()
	indexService, err := indexer.New()
	if err != nil {
		return
	}
	indexerInstance = indexService
	// load bootstrap tokens for indexing
	erc20Addresses := indexer.LoadTokensFromConfig()

	// Initialize ERC20 tokens for given addresses
	if err := indexService.Erc20TokensSubscriber(ctx, erc20Addresses); err != nil {
		logrus.Error("Fail to subscribe ERC20Tokens and write to database", "err", err)
		return
	}
	//  start the periodic task
	indexService.Listen()

	// get transfer events through Api
	router := mux.NewRouter()
	fmt.Println("HI")

	router.HandleFunc("/add", addAddress).Methods("POST")

	http.ListenAndServe(":8080", router)
}

func addAddress(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var erc20 models.AccountPayload
	logrus.Info(json.NewDecoder(request.Body))
	err := json.NewDecoder(request.Body).Decode(&erc20)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{"message":"failure"}`))
		return
	}
	logrus.Info(erc20)

	err = indexerInstance.GetOrAddErc20(context.Background(), erc20.Address, erc20.Name)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	response.WriteHeader(http.StatusOK)
	response.Write([]byte(`{"message":"success"}`))
}

var indexerInstance *indexer.Indexer
