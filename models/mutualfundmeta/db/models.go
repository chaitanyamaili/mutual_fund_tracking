package db

import (
	"database/sql"
	"time"
)

type MutualFundMeta struct {
	ID             int64        `db:"id"`
	FundHouse      string       `db:"fund_house"`
	SchemaType     string       `db:"scheme_type"`
	SchemaCategory string       `db:"scheme_category"`
	SchemaCode     string       `db:"scheme_code"`
	SchemaName     string       `db:"scheme_name"`
	CreatedOn      time.Time    `db:"created_on"`
	UpdatedOn      time.Time    `db:"updated_on"`
	DeletedOn      sql.NullTime `db:"deleted_on"`
	// DeletedOn      time.Time `db:"deleted_on"`
}
