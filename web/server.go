package web

import (
	"net/http"

	"github.com/mbndr/mtglib"
	"github.com/mbndr/mtglib/web/detail"
	"github.com/mbndr/mtglib/web/index"
)

// Serve starts the http server
func Serve(port string) error {
	// load all cards
	cards, err := mtglib.LoadCards()
	if err != nil {
		return err
	}

	symbols, err := mtglib.LoadSymbols()
	if err != nil {
		return err
	}

	// setting up webserver
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/detail/", detail.NewHandler(cards, symbols))
	http.Handle("/", index.NewHandler(cards, symbols, mtglib.TotalLibraryCardCount()))

	return http.ListenAndServe(port, nil)
}

func httpError(w http.ResponseWriter, err error) {
	http.Error(w, "500 internal server error", http.StatusInternalServerError)
}
