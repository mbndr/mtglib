package filter

type RuleFilter struct {
	Text string
}

func (f *RuleFilter) WhereClause() string {
	return "s.oracle_text LIKE  '%' || ? || '%'"
}

func (f *RuleFilter) Parameters() []interface{} {
	return []interface{}{f.Text}
}
