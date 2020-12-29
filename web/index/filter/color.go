package filter

import "strings"

type ColorFilter struct {
	Colors        []string
	MonocolorOnly bool
}

func (f *ColorFilter) WhereClause() string {
	var whereClauses []string
	whereStr := ""

	if len(f.Colors) != 0 {
		for range f.Colors {
			whereClauses = append(whereClauses, "s.colors LIKE '%' || ? || '%'")
		}

		whereStr = "(" + strings.Join(whereClauses, " OR ") + ")"
	}

	if f.MonocolorOnly {
		if whereStr != "" {
			whereStr += " AND "
		}
		whereStr += "s.colors NOT LIKE '%|%'"
	}

	return whereStr
}

func (f *ColorFilter) Parameters() []interface{} {
	p := make([]interface{}, len(f.Colors))

	for i, c := range f.Colors {
		p[i] = c
	}

	return p
}
