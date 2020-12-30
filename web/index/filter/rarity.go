package filter

type RarityFilter struct {
	Rarity string
}

func (f *RarityFilter) WhereClause() string {
	return "s.rarity = ?"
}

func (f *RarityFilter) Parameters() []interface{} {
	return []interface{}{f.Rarity}
}
