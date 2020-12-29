package filter

type TypeFilter struct {
	Type string
}

func (f *TypeFilter) WhereClause() string {
	return "s.type_line LIKE  '%' || ? || '%'"
}

func (f *TypeFilter) Parameters() []interface{} {
	return []interface{}{f.Type}
}
