package mutualfundmeta

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/chaitanyamaili/mutual_fund_tracking/models/mutualfundmeta/db"
	"github.com/chaitanyamaili/mutual_fund_tracking/pkg/database"
	"github.com/chaitanyamaili/mutual_fund_tracking/pkg/validate"
	"github.com/jmoiron/sqlx"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound     = errors.New("mutual_fund_meta not found")
	ErrInvalidID    = errors.New("ID is not in its proper form")
	ErrInvalidAlias = errors.New("alias is not in its proper form")
	//ErrAuthenticationFailure = errors.New("authentication failed")
)

// Core manages the set of APIs for pullrequest access
type Core struct {
	store db.Store
}

// NewCore constructs a core for pullrequest api access.
func NewCore(log *slog.Logger, sqlxDB *sqlx.DB, rwmux *sync.RWMutex) Core {
	return Core{
		store: db.NewStore(log, sqlxDB, rwmux),
	}
}

// -----------------------------------------------------------------------
// CRUD Methods
// -----------------------------------------------------------------------

// Create inserts a new pullrequest into the database
func (c Core) Create(ctx context.Context, rs NewMutualFundMeta, now time.Time) (MutualFundMeta, error) {
	if err := validate.Check(rs); err != nil {
		return MutualFundMeta{}, err
	}

	dbRS := db.MutualFundMeta{
		FundHouse:      rs.FundHouse,
		SchemaType:     rs.SchemaType,
		SchemaCategory: rs.SchemaCategory,
		SchemaCode:     rs.SchemaCode,
		SchemaName:     rs.SchemaName,
		CreatedOn:      now,
		UpdatedOn:      now,
	}

	// This provides an example of how to execute a transaction if required.
	tran := func(tx sqlx.ExtContext) error {
		res, err := c.store.Tran(tx).Create(ctx, dbRS)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}
		dbRS.ID = res.LastInsertID
		return nil
	}

	if err := c.store.WithinTran(ctx, tran); err != nil {
		return MutualFundMeta{}, fmt.Errorf("tran: %w", err)
	}

	return toMutualFundMeta(dbRS), nil
}

// Update replaces a pullrequest document in the database.
func (c Core) Update(ctx context.Context, id string, urs UpdateMutualFundMeta, now time.Time) error {
	if err := validate.Check(urs); err != nil {
		return err
	}
	if err := validate.CheckID(id); err != nil {
		return ErrInvalidID
	}

	dbRS, err := c.store.QueryByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating requesting tech id[%s]: %w", id, err)
	}

	isEmpty := true
	if urs.FundHouse != "" {
		dbRS.FundHouse = strings.TrimSpace(*&urs.FundHouse)
		isEmpty = false
	}
	if urs.SchemaCategory != "" {
		dbRS.SchemaCategory = strings.TrimSpace(*&urs.SchemaCategory)
		isEmpty = false
	}
	if urs.SchemaCode != "" {
		dbRS.SchemaCode = strings.TrimSpace(*&urs.SchemaCode)
		isEmpty = false
	}
	if urs.SchemaName != "" {
		dbRS.SchemaName = strings.TrimSpace(*&urs.SchemaName)
		isEmpty = false
	}
	if urs.SchemaType != "" {
		dbRS.SchemaType = strings.TrimSpace(*&urs.SchemaType)
		isEmpty = false
	}

	// No changes were made - don't touch the DB
	if isEmpty {
		return nil
	}
	dbRS.UpdatedOn = now

	_, err = c.store.Update(ctx, dbRS)
	if err != nil {
		return fmt.Errorf("update id[%s]: %w", id, err)
	}

	return nil
}

// Delete removes a pullrequest from the database.
func (c Core) Delete(ctx context.Context, id string, now time.Time) error {
	if err := validate.CheckID(id); err != nil {
		return ErrInvalidID
	}

	_, err := c.store.QueryByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("undeleting requesting tech id[%s]: %w", id, err)
	}

	_, err = c.store.Delete(ctx, id, now)
	if err != nil {
		return fmt.Errorf("delete id[%s]: %w", id, err)
	}

	return nil
}

// Query retrieves a list of existing records from the database
func (c Core) Query(ctx context.Context, pagi database.Pagination) ([]MutualFundMeta, error) {
	res, err := c.store.Query(ctx, pagi)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return []MutualFundMeta{}, ErrNotFound
		}
		return []MutualFundMeta{}, fmt.Errorf("query: %w", err)
	}

	return toMutualFundMetaSlice(res), nil
}

// QueryByID retrieves a single records from the database by id
func (c Core) QueryByID(ctx context.Context, id string) (MutualFundMeta, error) {
	if err := validate.CheckID(id); err != nil {
		return MutualFundMeta{}, ErrInvalidID
	}

	res, err := c.store.QueryByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return MutualFundMeta{}, ErrNotFound
		}
		return MutualFundMeta{}, fmt.Errorf("query: %w", err)
	}

	return toMutualFundMeta(res), nil
}
