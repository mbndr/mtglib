package web

import (
	"html/template"
	"net/http"

	"github.com/mbndr/mtglib"
)

// Serve starts the http server
// TODO: cleanup
func Serve(port string) error {
	// load all cards
	cards, oracleIDs, err := mtglib.LoadCards()
	if err != nil {
		return err
	}

	symbolList, err := mtglib.LoadSymbols()
	if err != nil {
		return err
	}

	// index page
	tpl, err := template.ParseFiles("html/index.html")
	if err != nil {
		return err
	}

	index := &indexHandler{
		DistinctCardCount: len(cards),
		TotalCardCount:    mtglib.TotalLibraryCardCount(),
		cards:             cards,
		oracleIDs:         oracleIDs,
		tpl:               tpl,
	}

	// detail page
	tpl, err = template.ParseFiles("html/detail.html")
	if err != nil {
		return err
	}

	detail := &detailHandler{
		cards:   cards,
		symbols: symbolList, // TODO: use custom symbolcollection
		tpl:     tpl,
	}

	// setting up webserver
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/detail/", detail)
	http.Handle("/", index)

	return http.ListenAndServe(port, nil)
}

func httpError(w http.ResponseWriter, err error) {
	http.Error(w, "500 internal server error", http.StatusInternalServerError)
}
