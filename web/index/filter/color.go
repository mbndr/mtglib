package filter

import "strings"

// ColorFilter filters for color identity of a card
type ColorFilter struct {
	Colors        []string
	MonocolorOnly bool
}

func (f *ColorFilter) WhereClause() string {
	var whereClauses []string
	whereStr := ""

	if len(f.Colors) != 0 {
		for range f.Colors {
			whereClauses = append(whereClauses, "s.color_identity LIKE '%' || ? || '%'")
		}

		whereStr = "(" + strings.Join(whereClauses, " OR ") + ")"
	}

	if f.MonocolorOnly {
		if whereStr != "" {
			whereStr += " AND "
		}
		whereStr += "s.color_identity NOT LIKE '%|%'"
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
