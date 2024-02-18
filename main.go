package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/chaitanyamaili/mutual_fund_tracking/internal/handlers"
	"github.com/chaitanyamaili/mutual_fund_tracking/pkg/database"
	"github.com/chaitanyamaili/mutual_fund_tracking/pkg/logger"
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
	// meterName is the name of the meter.
	meterName = "github.com/chaitanyamaili/mutual_fund_tracking"
)

// --------------------------------------------------
// Logger
// --------------------------------------------------
// Initiate a new logger.
var log = logger.WithFormatter(os.Stdout, DefaultAddSource, DefaultLogFormat)

// var provider *metric.MeterProvider

func main() {
	if err := run(); err != nil {
		log.Error(err.Error())
		// panic(err)
	}
}

func run() error {
	log.Info("Starting server...")
	// -------------------------------------------------------------------
	// Databases
	// -------------------------------------------------------------------
	log.Info("startup.db", "status", "initializing DBs")

	db, err := database.Open(database.Config{
		Type:         "postgres",
		User:         "postgres",
		Password:     "root",
		Host:         "localhost", //"mf_db","172.25.0.2"
		Port:         6543,
		Name:         "mf",
		MaxIdleConns: 0,
		MaxOpenConns: 0,
		DisableTLS:   true,
	})
	if err != nil {
		log.Error("startup.db", "status", "error", "error", err.Error())
		panic(fmt.Errorf("connecting to db: %w", err))
	}
	defer func() {
		log.Info("shutdown", "status", "stopping db", "host", "localhost")
		_ = db.Close()
	}()

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Error("Ping.db", "status", "error", "error", err.Error())
	} else {
		log.Info("Ping.db", "status", "success")
	}

	// -------------------------------------------------------------------
	// RWMux for lock DBs in transaction mode (deadlocks = yuck)
	// -------------------------------------------------------------------
	log.Info("startup.remux", "status", "created")
	rwmux := &sync.RWMutex{}
	apiMux := handlers.APIMux(
		handlers.APIMuxConfig{
			Log:   log,
			DB:    db,
			RWMux: rwmux,
		},
	)

	log.Info("startup.api", "status", "created")
	api := http.Server{
		Addr:    "localhost:8080",
		Handler: apiMux,
	}

	log.Debug(
		fmt.Sprintf("\n%s\nAPI STARTED\nHost: %s\n%s\n",
			strings.Repeat("-", 75),
			api.Addr,
			strings.Repeat("-", 75),
		),
	)
	err = api.ListenAndServe()
	if err != nil {
		log.Error("startup.api", "status", "error", "error", err.Error())
		panic(fmt.Errorf("starting server: %w", err))
	}

	return nil
}

// func scrapeMFData(w http.ResponseWriter, r *http.Request) {
// 	log.Info(fmt.Sprintf("Scraping %s with method %s", r.URL.Path, r.Method))
// 	params := httptreemux.ContextParams(r.Context())
// 	fmt.Fprintf(w, "Scraping %s\n", params["key"])
// 	resData, err := mutualfund.NewHandler(log).GetLatestNavData(params["key"])
// 	if err != nil {
// 		log.Error(fmt.Sprintf("Failed to get nav data for %s, error: %s", params["key"], err.Error()))
// 	}
// 	fmt.Fprintf(w, "Fund: %s, Latest NAV Value: %s for date: %s\n", resData.Meta.FuncdHouse, resData.Data[0].Nav, resData.Data[0].Date)

// 	meter := provider.Meter(meterName)

// 	opt := api.WithAttributes(
// 		attribute.Key("A").String("B"),
// 		attribute.Key("C").String("D"),
// 	)

// 	// This is the equivalent of prometheus.NewHistogramVec
// 	histogram, err := meter.Float64Histogram(
// 		"nav_value",
// 		api.WithDescription(fmt.Sprintf("FuncdHouse: %s, SchemaName: %s histogram", resData.Meta.FuncdHouse, resData.Meta.SchemeName)),
// 	)
// 	if err != nil {
// 		log.Error(fmt.Sprintf("Failed to create histogram: %v", err))
// 	}
// 	navFloat64, _ := strconv.ParseFloat(resData.Data[0].Nav, 64)
// 	histogram.Record(r.Context(), navFloat64, opt)
// }
