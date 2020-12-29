package sorting

// Sorting could be done like the filters, but they are more complex
type Sorting struct {
	SortBy    string
	SortOrder string
}

func (s *Sorting) Validate() bool {
	if s.SortOrder != "asc" && s.SortOrder != "desc" {
		return false
	}

	// TODO: rarity?
	if s.SortBy != "cmc" && s.SortBy != "type" && s.SortBy != "name" {
		return false
	}

	return true
}

// GetOrderBy returns the string relevant for the SQL query
func (s *Sorting) GetOrderBy() string {
	sortBy := ""

	if s.SortBy == "cmc" {
		sortBy = "s.cmc"
	} else if s.SortBy == "type" {
		sortBy = "s.type_line"
	} else if s.SortBy == "name" {
		sortBy = "s.name"
	}

	return sortBy + " " + s.SortOrder
}
