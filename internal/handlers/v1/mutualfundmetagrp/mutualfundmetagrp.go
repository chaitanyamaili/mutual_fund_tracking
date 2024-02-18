package mutualfundmetagrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/chaitanyamaili/mutual_fund_tracking/models/mutualfundmeta"
	api "github.com/chaitanyamaili/mutual_fund_tracking/pkg/api"
	"github.com/chaitanyamaili/mutual_fund_tracking/pkg/database"
)

// Handlers is the set of all handlers for mutualfundmetagrp.
type Handlers struct {
	MutualFundMeta mutualfundmeta.Core
}

// Create a new mutual_fund_meta record.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := api.GetContextValues(ctx)
	if err != nil {
		fmt.Printf("unable to get context: %v", err)
		return err
	}

	var newRecord mutualfundmeta.NewMutualFundMeta
	if err := api.Decode(r, &newRecord); err != nil {
		fmt.Printf("unable to decode request: %v", err)
		return api.NewRequestError(err, http.StatusBadRequest)
	}

	rs, err := h.MutualFundMeta.Create(ctx, newRecord, v.Now)
	if err != nil {
		fmt.Printf("unable to create mutial_fund_meta: %v", err)
		return err
	}

	return api.Respond(ctx, w, rs, http.StatusCreated)

}

// Query all the mutual_fund_meta records
func (h Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Get the pagination parameters from the request.
	fmt.Println("Querying for MutualFundMeta")
	pagi, err := database.PaginationParams(r)
	if err != nil {
		fmt.Printf("unable to get pagination parameters: %v", err)
		return err
	}

	// Query for the MutualFundMeta records.
	fmt.Println("Querying for MutualFundMeta records")
	rs, err := h.MutualFundMeta.Query(ctx, pagi)
	if err != nil {
		fmt.Printf("unable to query for MutualFundMeta: %v\n", err)
		switch {
		case errors.Is(err, mutualfundmeta.ErrNotFound):
			return api.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("unable to query for Requesting Tech: %w", err)
		}
	}
	fmt.Printf("MutualFundMeta query successful: %v\n", rs)

	return api.Respond(ctx, w, rs, http.StatusOK)
}
