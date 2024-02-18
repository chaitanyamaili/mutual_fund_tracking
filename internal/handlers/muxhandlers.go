package handlers

import (
	"log/slog"
	"net/http"
	"sync"

	v1 "github.com/chaitanyamaili/mutual_fund_tracking/internal/handlers/v1"
	"github.com/chaitanyamaili/mutual_fund_tracking/pkg/api"
	"github.com/chaitanyamaili/mutual_fund_tracking/pkg/api/middleware"
	"github.com/jmoiron/sqlx"
)

type Options struct {
	corsOrigin string
}

type APIMuxConfig struct {
	Log     *slog.Logger
	DB      *sqlx.DB
	RWMux   *sync.RWMutex
	Headers bool
}

func APIMux(cfg APIMuxConfig, options ...func(opts *Options)) http.Handler {
	var opts Options
	for _, option := range options {
		option(&opts)
	}

	// Construct the web.App which holds all routes as well as common Middleware.
	mw := make([]api.Middleware, 0, 6)
	mw = append(mw, middleware.Logger(cfg.Log))
	a := api.NewAPI(
		mw...,
	)

	// Register the 404 path not found so that we can log it
	a.NotFound(cfg.Log)

	// Accept CORS 'OPTIONS' preflight requests if config has been provided.
	// Don't forget to apply the CORS middleware to the routes that need it.
	// Example Config: `conf:"default:https://MY_DOMAIN.COM"`
	// if opts.corsOrigin != "" {
	// 	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// 		return nil
	// 	}
	// 	a.Handle(http.MethodOptions, "", h, middleware.Cors(opts.corsOrigin))
	// }

	// Load the v1 routes.
	v1.Routes(a, v1.Config{
		Log:   cfg.Log,
		DB:    cfg.DB,
		RWMux: cfg.RWMux,
	})

	return a
}
