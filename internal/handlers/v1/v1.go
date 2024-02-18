package v1

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/chaitanyamaili/mutual_fund_tracking/internal/handlers/v1/mutualfundmetagrp"
	"github.com/chaitanyamaili/mutual_fund_tracking/models/mutualfundmeta"
	"github.com/chaitanyamaili/mutual_fund_tracking/pkg/api"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Log   *slog.Logger
	DB    *sqlx.DB
	RWMux *sync.RWMutex
}

func Routes(api *api.API, cfg Config) {
	// -------------------------------------------------------------------
	// Mutual Fund Meta Group
	// -------------------------------------------------------------------
	mfmd := mutualfundmetagrp.Handlers{
		MutualFundMeta: mutualfundmeta.NewCore(cfg.Log, cfg.DB, cfg.RWMux),
	}
	api.Handle(http.MethodGet, "/v1/mutualfundmeta", mfmd.Query)
	api.Handle(http.MethodPost, "/v1/mutualfundmeta", mfmd.Create)

	// -------------------------------------------------------------------
	// Add in the Teapot
	// -------------------------------------------------------------------
	api.Handle(http.MethodGet, "/v1/teapot", Teapot)
}

func Teapot(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	lyrics := `I'm a little teapot, Short and stout,
Here is my handle. Here is my spout.
When I get all steamed up, Hear me shout,
Tip me over and pour me out!`

	type o struct {
		Lyrics string `json:"lyrics"`
	}
	output := []o{{Lyrics: lyrics}}

	return api.Respond(ctx, w, output, http.StatusTeapot)
}
