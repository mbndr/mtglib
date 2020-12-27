package mtglib

import (
	"database/sql"
	"html/template"
	"path"
	"strings"

	"github.com/mbndr/mtglib/db"
	"github.com/mbndr/mtglib/scryfall"
)

func LoadSymbols() ([]scryfall.Symbol, error) {
	list := []scryfall.Symbol{}

	err := db.Select("SELECT symbol, svg_uri, title FROM scryfall_symbols", func(rows *sql.Rows) error {
		var s scryfall.Symbol
		err := rows.Scan(&s.Symbol, &s.SvgURI, &s.Title)
		list = append(list, s)
		return err
	})

	return list, err
}

func SymbolToHTMLImage(symbols string, symbolList []scryfall.Symbol) []template.HTML {
	symbolSplit := strings.SplitAfter(symbols, "}")
	imageLinks := []template.HTML{}

	// TODO: do this into symbol wrapper list object
	smap := make(map[string]scryfall.Symbol)
	for _, sym := range symbolList {
		smap[sym.Symbol] = sym
	}

	for _, sym := range symbolSplit[:len(symbolSplit)-1] { // TODO: return whole html (with title etc)
		imageLinks = append(imageLinks, template.HTML("<img width=\"20\" src=\"/resources/"+path.Base(smap[sym].SvgURI)+"\" title=\""+smap[sym].Title+"\" />"))
	}

	return imageLinks
}
