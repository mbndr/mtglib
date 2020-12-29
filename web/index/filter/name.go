package filter

type NameFilter struct {
	Name string
}

func (f *NameFilter) WhereClause() string {
	return "s.name LIKE  '%' || ? || '%'"
}

func (f *NameFilter) Parameters() []interface{} {
	return []interface{}{f.Name}
}
