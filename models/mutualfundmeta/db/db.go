package db

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chaitanyamaili/mutual_fund_tracking/pkg/database"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	log          *slog.Logger
	tr           database.Transactor
	db           sqlx.ExtContext
	rwmux        *sync.RWMutex
	isWithinTran bool
}

// NewStore constructs a data for api access.
func NewStore(log *slog.Logger, db *sqlx.DB, rwmux *sync.RWMutex) Store {
	return Store{
		log:   log,
		tr:    db,
		db:    db,
		rwmux: rwmux,
	}
}

// WithinTran runs passes function and do commit/rollback at the end.
func (s Store) WithinTran(ctx context.Context, fn func(sqlx.ExtContext) error) error {
	if s.isWithinTran {
		return fn(s.db)
	}
	s.rwmux.Lock()
	err := database.WithinTran(ctx, s.log, s.tr, fn)
	s.rwmux.Unlock()

	return err
}

// Tran return new Store with transaction in it.
func (s Store) Tran(tx sqlx.ExtContext) Store {
	return Store{
		log:          s.log,
		tr:           s.tr,
		db:           tx,
		isWithinTran: true,
	}
}

// -----------------------------------------------------------------------
// Database Query MutualFundMeta
// -----------------------------------------------------------------------

// Create inserts a new mutual_fund_mete into the database.
func (s Store) Create(ctx context.Context, rs MutualFundMeta) (database.DBResults, error) {
	const q = `
	INSERT INTO public.mutual_fund_meta
		(fund_house,
		scheme_type,
		scheme_category,
		scheme_code,
		scheme_name,
		created_on,
		updated_on)
	VALUES
		(:fund_house, 
		:scheme_type,
		:scheme_category,
		:scheme_code,
		:scheme_name,
		:created_on,
		:updated_on)`

	fmt.Println(q)
	res, err := database.NamedExecContext(ctx, s.log, s.db, q, rs)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return database.DBResults{}, database.NewError(database.ErrDBDuplicatedEntry, http.StatusConflict)
		}

		return database.DBResults{}, fmt.Errorf("inserting public.mutual_fund_meta: %w", err)
	}

	return res, nil
}

// Update replaces a pullrequest record in the database.
func (s Store) Update(ctx context.Context, rs MutualFundMeta) (database.DBResults, error) {
	const q = `
	UPDATE
		public.mutual_fund_meta
	SET
		fund_house = :fund_house,
		scheme_type = :scheme_type,
		scheme_category = :scheme_category,
		scheme_code = :scheme_code,
		scheme_name = :scheme_name,
		updated_on = :updated_on
	WHERE
		id = :id`

	res, err := database.NamedExecContext(ctx, s.log, s.db, q, rs)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return database.DBResults{}, database.NewError(database.ErrDBDuplicatedEntry, http.StatusConflict)
		}
		return database.DBResults{}, fmt.Errorf("updating public.mutual_fund_meta ID[%d]: %w", rs.ID, err)
	}

	return res, nil
}

// Delete removes a pullrequest from the database.
func (s Store) Delete(ctx context.Context, id string, now time.Time) (database.DBResults, error) {
	data := struct {
		ID        string    `db:"id"`
		DeletedOn time.Time `db:"deleted_on"`
	}{
		ID:        id,
		DeletedOn: now,
	}

	const q = `
	UPDATE
		public.mutual_fund_meta
	SET
		deleted_on = :deleted_on
	WHERE
		id = :id`

	res, err := database.NamedExecContext(ctx, s.log, s.db, q, data)
	if err != nil {
		return database.DBResults{}, fmt.Errorf("deleting public.mutual_fund_meta id[%s]: %w", id, err)
	}

	return res, nil
}

// Query retrieves a list of existing pullrequest from the database.
func (s Store) Query(ctx context.Context, pagi database.Pagination) ([]MutualFundMeta, error) {
	fmt.Println("In query")
	fmt.Println(pagi)
	q := database.PaginationQuery(pagi, `
	SELECT
		id,
		fund_house,
		scheme_type,
		scheme_category,
		scheme_code,
		scheme_name,
		created_on,
		updated_on,
		deleted_on
	FROM
		public.mutual_fund_meta
	WHERE
		deleted_on is null
	ORDER BY
		:sort :direction,
		id :direction`)

	// Slice to hold results
	fmt.Println(q)
	var res []MutualFundMeta
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, pagi, &res); err != nil {
		return nil, fmt.Errorf("selecting requesting tech: %w", err)
	}

	return res, nil
}

// QueryByID retrieves a list of existing pullrequests from the database.
func (s Store) QueryByID(ctx context.Context, id string) (MutualFundMeta, error) {
	data := struct {
		ID string `db:"id"`
	}{ID: id}
	const q = `
	SELECT
		id,
		fund_house,
		scheme_type,
		scheme_category,
		scheme_code,
		scheme_name,
		created_on,
		updated_on,
		deleted_on
	FROM
		public.mutual_fund_meta
	WHERE
		id = :id
		and deleted_on is null`

	// Slice to hold results
	var res MutualFundMeta
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &res); err != nil {
		return MutualFundMeta{}, err
	}

	return res, nil
}
