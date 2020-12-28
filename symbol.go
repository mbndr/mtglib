package mtglib

import (
	"database/sql"
	"html/template"
	"path"
	"regexp"
	"strings"

	"github.com/mbndr/mtglib/db"
)

// Symbol represents a card symbol
type Symbol struct {
	Symbol string `json:"symbol"`
	SvgURI string `json:"svg_uri"`
	Title  string `json:"english"`
}

// SymbolCollection holds all symbols and has functions to parse them
type SymbolCollection struct {
	// key is the symbol as string (e.g. {G})
	symbols map[string]Symbol
}

var symbolRegex = regexp.MustCompile(`({.*?})`)

// LoadSymbols loads all symbols
func LoadSymbols() (SymbolCollection, error) {
	sc := SymbolCollection{
		make(map[string]Symbol),
	}

	err := db.Select("SELECT symbol, svg_uri, title FROM scryfall_symbols", func(rows *sql.Rows) error {
		var s Symbol
		err := rows.Scan(&s.Symbol, &s.SvgURI, &s.Title)
		sc.symbols[s.Symbol] = s
		return err
	})

	return sc, err
}

func (sc *SymbolCollection) symbolToImage(sym string) string {
	symbol, ok := sc.symbols[sym]
	if !ok {
		return sym
	}
	return "<img class=\"ml-symbol\" src=\"/resources/" + path.Base(symbol.SvgURI) + "\" title=\"" + symbol.Title + "\" />"
}

// HTMLImages parses a symbol string (e.g. {2}{G} to HTML images)
func (sc *SymbolCollection) HTMLImages(symbolStr string) template.HTML {
	symbolSplit := strings.SplitAfter(symbolStr, "}")
	images := ""

	for _, sym := range symbolSplit[:len(symbolSplit)-1] {
		images += sc.symbolToImage(sym)
	}

	return template.HTML(images)
}

// ParseInText pareses symbols in a text
func (sc *SymbolCollection) ParseInText(text string) template.HTML {
	text = symbolRegex.ReplaceAllStringFunc(text, sc.symbolToImage)

	str := ""
	for _, s := range strings.Split(text, "\n") {
		str += "<p>" + s + "</p>"
	}

	return template.HTML(str)
}
