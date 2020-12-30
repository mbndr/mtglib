package filter

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"

	"github.com/mbndr/mtglib/db"
	"github.com/mbndr/mtglib/web/index/sorting"
)

var LogSQL = false

const baseSelect = `SELECT s.oracle_id FROM helvault_library h INNER JOIN scryfall_cards s ON s.scryfall_id = h.scryfall_id`

type Filter interface {
	WhereClause() string
	Parameters() []interface{}
}

// GetFilterResult gets all oracleIDs from the database with applied filters (AND CONCATINATED)
func GetFilterResult(filters []Filter, sorter *sorting.Sorting) ([]string, error) {
	sqlStmt, parameters := buildSQLStatement(filters, sorter)

	var oracleIDs []string

	err := db.Select(sqlStmt, func(rows *sql.Rows) error {
		var oid string
		err := rows.Scan(&oid)
		oracleIDs = append(oracleIDs, oid)
		return err
	}, parameters...)

	return oracleIDs, err
}

// build sql statement and return it with a list of parameters
func buildSQLStatement(filters []Filter, sorter *sorting.Sorting) (string, []interface{}) {
	// prepare
	whereClauses := make([]string, len(filters))
	var parameters []interface{} // cannot tell total length

	for i, f := range filters {
		whereClauses[i] = f.WhereClause()
		parameters = append(parameters, f.Parameters()...)
	}

	// build
	buf := bytes.NewBufferString(baseSelect)
	if len(filters) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(whereClauses, " AND "))
	}
	buf.WriteString(" GROUP BY s.oracle_id")

	if sorter != nil {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(sorter.GetOrderBy())
	}

	parameters = append(parameters)

	if LogSQL {
		fmt.Printf("SQL: %s\nParameter: %+v\n", buf.String(), parameters)
	}
	return buf.String(), parameters
}
