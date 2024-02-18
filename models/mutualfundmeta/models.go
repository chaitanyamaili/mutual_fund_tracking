package mutualfundmeta

import (
	"time"
	"unsafe"

	"github.com/chaitanyamaili/mutual_fund_tracking/models/mutualfundmeta/db"
)

// MutualFundMeta is the mutual fund meta data.
type MutualFundMeta struct {
	ID             int       `json:"id"`
	FundHouse      string    `json:"fund_house"`
	SchemaType     string    `json:"scheme_type"`
	SchemaCategory string    `json:"scheme_category"`
	SchemaCode     string    `json:"scheme_code"`
	SchemaName     string    `json:"scheme_name"`
	CreatedOn      time.Time `json:"created_on"`
	UpdatedOn      time.Time `json:"updated_on"`
	DeleteOn       time.Time `json:"delete_on"`
}

// NewMutualFundMeta is the new mutual fund meta data.
type NewMutualFundMeta struct {
	FundHouse      string `json:"fund_house"`
	SchemaType     string `json:"scheme_type"`
	SchemaCategory string `json:"scheme_category"`
	SchemaCode     string `json:"scheme_code"`
	SchemaName     string `json:"scheme_name"`
}

// UpdateMutualFundMeta is the update mutual fund meta data.
type UpdateMutualFundMeta struct {
	FundHouse      string `json:"fund_house"`
	SchemaType     string `json:"scheme_type"`
	SchemaCategory string `json:"scheme_category"`
	SchemaCode     string `json:"scheme_code"`
	SchemaName     string `json:"scheme_name"`
}

func toMutualFundMeta(dbRS db.MutualFundMeta) MutualFundMeta {
	p := (*MutualFundMeta)(unsafe.Pointer(&dbRS))
	return *p
}

func toMutualFundMetaSlice(dbSRs []db.MutualFundMeta) []MutualFundMeta {
	rs := make([]MutualFundMeta, len(dbSRs))
	for i, dbSR := range dbSRs {
		rs[i] = toMutualFundMeta(dbSR)
	}
	return rs
}
