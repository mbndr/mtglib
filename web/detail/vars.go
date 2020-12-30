package detail

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/mbndr/mtglib"
)

type detailVars struct {
	Handler *Handler
	Card    mtglib.Card
}

func (v *detailVars) CardName(c mtglib.Card) string {
	if len(c.CardFaces) == 0 {
		return c.Name
	}
	return strings.Replace(c.Name, "//", "|", -1)
}

// CardURL returns the link to the cards image.
func (v *detailVars) CardImageURL(oracleID string) string {
	return v.Handler.cards.CardImageURL(oracleID)
}

// CardURL returns the link to the cards image.
func (v *detailVars) FaceImageURLs(oracleID string) []string {
	return v.Handler.cards.FaceImageURLs(oracleID)
}

// ManaSymbols returns mana symbols from a string like "{R}{G}"
func (v *detailVars) ManaSymbols(symbols string) template.HTML {
	if symbols == "" {
		return "-"
	}
	return v.Handler.symbols.HTMLImages(symbols)
}

// ManaSymbolsArr returns mana symbols from a array like []string{"R", "G"}
func (v *detailVars) ManaSymbolsArr(symbols []string) template.HTML {
	var str string
	for _, s := range symbols {
		str += fmt.Sprintf("{%s}", s)
	}
	return v.ManaSymbols(str)
}

// Parse symbols in text and structure it with paragraphs
func (v *detailVars) ParseOracleText(oracleText string) template.HTML {
	return v.Handler.symbols.ParseInText(oracleText)
}
