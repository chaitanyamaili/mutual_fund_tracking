package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/chaitanyamaili/mutual_fund_tracking/mutualfund"
	"github.com/chaitanyamaili/mutual_fund_tracking/pkg/logger"
	"github.com/dimfeld/httptreemux"
)

const (
	// DefaultLogFormat is 'json'.
	DefaultLogFormat = "json"
	// DefaultHttpPort is 8080.
	DefaultHttpPort = ":8080"
	// DefaultAddSource is false.
	DefaultAddSource = false
	// MutualFundLatestBaseURL is https://api.mfapi.in/mf/%s/latest
	MutualFundLatestBaseURL = "https://api.mfapi.in/mf/%s/latest"
	// MutualFundHistoryBaseURL is https://api.mfapi.in/mf/%s
	MutualFundHistoryBaseURL = "https://api.mfapi.in/mf/%s"
)

// --------------------------------------------------
// Logger
// --------------------------------------------------
// Initiate a new logger.
var log = logger.WithFormatter(os.Stdout, DefaultAddSource, DefaultLogFormat)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	log.Info("Starting server...")
	router := httptreemux.NewContextMux()
	router.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// implement the routes here
	// v1/mutualfund/scrape/[key]
	router.GET("/v1/mutualfund/scrape/:key", scrapeMFData)

	http.ListenAndServe(DefaultHttpPort, router)
	log.Info(fmt.Sprintf("Server started on port %s", DefaultHttpPort))
	return nil
}

func scrapeMFData(w http.ResponseWriter, r *http.Request) {
	log.Info(fmt.Sprintf("Scraping %s with method %s", r.URL.Path, r.Method))
	params := httptreemux.ContextParams(r.Context())
	fmt.Fprintf(w, "Scraping %s\n", params["key"])
	resData, err := mutualfund.NewHandler(log).GetLatestNavData(params["key"])
	if err != nil {
		log.Error(fmt.Sprintf("Failed to get nav data for %s, error: %s", params["key"], err.Error()))
	}
	fmt.Fprintf(w, "Fund: %s, Latest NAV Value: %s for date: %s\n", resData.Meta.FuncdHouse, resData.Data[0].Nav, resData.Data[0].Date)
}
