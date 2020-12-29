package filter

import (
	"bytes"
	"database/sql"
	"strings"

	"github.com/mbndr/mtglib/db"
)

// TODO: add limit
const baseSelect = `SELECT s.oracle_id FROM helvault_library h INNER JOIN scryfall_cards s ON s.scryfall_id = h.scryfall_id`

// TODO: implement filter interface which returns a where clause for sql?
type Filter interface {
	WhereClause() string
	Parameters() []interface{}
}

// GetFilterResult gets all oracleIDs from the database with applied filters (AND CONCATINATED)
func GetFilterResult(filters []Filter, offset, limit int) ([]string, error) {
	sqlStmt, parameters := buildSQLStatement(filters, offset, limit)

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
func buildSQLStatement(filters []Filter, offset, limit int) (string, []interface{}) {
	whereClauses := make([]string, len(filters))
	var parameters []interface{} // cannot tell total length

	for i, f := range filters {
		whereClauses[i] = f.WhereClause()
		parameters = append(parameters, f.Parameters()...)
	}

	buf := bytes.NewBufferString(baseSelect)
	if len(filters) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(whereClauses, " AND "))
	}
	buf.WriteString(" GROUP BY s.oracle_id ")
	buf.WriteString(" LIMIT ? OFFSET ?")

	parameters = append(parameters, limit, offset)

	return buf.String(), parameters
}
