package database

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Pagination details
type Pagination struct {
	// The current page
	//
	// in: query
	// type: integer
	// enum: 1
	// minimum: 1
	// required: false
	Page int `db:"page" json:"page"`
	// The per page limit
	//
	// in: query
	// type: integer
	// enum: 20
	// minimum: 1
	// maximum: 100
	// required: false
	PerPage int `db:"per_page" json:"per_page"`
	// The column to sort on
	//
	// in: query
	// required: false
	// type: string
	// enum: created,updated
	// description:
	//   Sort order:
	//   - `created` - When the record was created in the database
	//   - `updated` - When the record was last touched in the database
	Sort string `db:"sort" json:"sort"`
	// The direction of the sort
	//
	// in: query
	// required: false
	// type: string
	// enum: asc,desc
	// description:
	//   Sort order:
	//   - `asc` - Ascending, from A to Z
	//   - `desc` - Descending, from Z to A
	Direction string `db:"direction" json:"direction"`
}

// PaginationResults Pagination details
type PaginationResults struct {
	// The current page
	//
	// in: query
	// type: integer
	// enum: 1
	// minimum: 1
	// required: false
	Page int `db:"page" json:"page"`
	// The per page limit
	//
	// in: query
	// type: integer
	// enum: 20
	// minimum: 1
	// maximum: 100
	// required: false
	PerPage int `db:"per_page" json:"per_page"`
	// The column to sort on
	//
	// in: query
	// required: false
	// type: string
	// enum: created
	// description:
	//   Sort order:
	//   - `created` - When the record was created in the database
	Sort string `db:"sort" json:"sort"`
	// The direction of the sort
	//
	// in: query
	// required: false
	// type: string
	// enum: asc,desc
	// description:
	//   Sort order:
	//   - `asc` - Ascending, from A to Z
	//   - `desc` - Descending, from Z to A
	Direction string `db:"direction" json:"direction"`
}

// NewPagination to initialize the pagination
func NewPagination() Pagination {
	return Pagination{
		Page:      0,
		PerPage:   20,
		Sort:      "created_on",
		Direction: "desc",
	}
}

// PaginationParams simple function to get, validate, and compute
// pagination values
func PaginationParams(r *http.Request) (Pagination, error) {
	qparams := r.URL.Query()
	singleSpacePattern := regexp.MustCompile(`\s+`)

	pagi := NewPagination()
	fmt.Println("First")
	fmt.Println(pagi)

	if val, ok := qparams["per_page"]; ok {
		perPage, err := strconv.Atoi(val[0])
		if err != nil {
			return Pagination{}, fmt.Errorf("invalid perPage format: %s", val[0])
		}
		if perPage < 0 {
			perPage = 20
		}
		if perPage > 100 {
			perPage = 100
		}
		pagi.PerPage = perPage
	}

	if val, ok := qparams["page"]; ok {
		page, err := strconv.Atoi(val[0])
		if err != nil {
			return Pagination{}, fmt.Errorf("invalid page format: %s", val[0])
		}
		page = page - 1
		if page < 0 {
			page = 0
		}
		pagi.Page = page * pagi.PerPage
	}

	if val, ok := qparams["sort"]; ok {
		val[0] = singleSpacePattern.ReplaceAllString(val[0], "")
		sort := strings.ToLower(val[0])
		// "created" set by default
		if sort == "updated" {
			pagi.Sort = "updated_on"
		}
		if sort == "id" {
			pagi.Sort = "id"
		}
		fmt.Println("Last")
		fmt.Println(pagi)
	}

	if val, ok := qparams["direction"]; ok {
		val[0] = singleSpacePattern.ReplaceAllString(val[0], "")
		sort := strings.ToLower(val[0])
		// "created" set by default
		if sort == "asc" {
			pagi.Direction = "asc"
		}
	}

	return pagi, nil
}

// PaginationQuery Placeholders ('?') can only be used to insert dynamic,
//
//	escaped values for filter parameters (e.g. in the WHERE part), where
//
// data values should appear, not for SQL keywords, identifiers etc. You
// cannot use it to dynamically specify the ORDER BY OR GROUP BY values.
// https://stackoverflow.com/questions/30867337/golang-order-by-issue-with-mysql
func PaginationQuery(pagi Pagination, q string) string {
	q = strings.ReplaceAll(q, ":per_page", fmt.Sprint(pagi.PerPage))
	q = strings.ReplaceAll(q, ":page", fmt.Sprint(pagi.Page))
	q = strings.ReplaceAll(q, ":sort", pagi.Sort)
	q = strings.ReplaceAll(q, ":direction", pagi.Direction)
	return q
}
