package mtglib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ScryfallSymbol represents a mtg symbol
type ScryfallSymbol struct {
	Symbol string `json:"symbol"`
	SvgURI string `json:"svg_uri"`
	Title  string `json:"english"`
}

type scryfallSymbolList struct {
	Data []*ScryfallSymbol `json:"data"`
}

// ImportMeta imports metadata such as symbols
func ImportMeta() error {
	// TODO: also set data
	return importSymbols()
}

func importSymbols() error {
	err := truncateTable("scryfall_symbols")
	if err != nil {
		return err
	}

	res, err := http.Get("https://api.scryfall.com/symbology")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	symbolList := &scryfallSymbolList{}

	err = json.Unmarshal(bytes, &symbolList)
	if err != nil {
		return err
	}

	// TODO: HERE: REFACTOR: simpler DB insert, select etc (with callback?)
	return nil
}
