package mutualfund

// Meta is the meta data of the mutual fund.
type Meta struct {
	FuncdHouse     string `json:"fund_house"`
	SchemeType     string `json:"scheme_type"`
	SchemaCategory string `json:"scheme_category"`
	SchemeCode     int    `json:"scheme_code"`
	SchemeName     string `json:"scheme_name"`
}

// Data is the data of the mutual fund.
type Data struct {
	Date string `json:"date"`
	Nav  string `json:"nav"`
}

// MutualFund is the mutual fund data.
type MutualFund struct {
	Meta   Meta   `json:"meta"`
	Data   []Data `json:"data"`
	Status string `json:"status"`
}
