package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mbndr/logo"
	"github.com/mbndr/mtglib/db"
	"github.com/mbndr/mtglib/helvault"
	"github.com/mbndr/mtglib/scryfall"
	"github.com/mbndr/mtglib/web"
)

var (
	logger = logo.NewSimpleLogger(os.Stderr, logo.INFO, "MTGLIB ", true)

	actionServe  bool
	actionImport bool
	helvaultFile string

	serveAddr string
)

// TODO: panic when no helvault library imported
func main() {
	flag.BoolVar(&actionServe, "serve", false, "Start the webserver")
	flag.StringVar(&serveAddr, "addr", ":8080", "Adress of the webserver")
	flag.BoolVar(&actionImport, "import", false, "(Re)import scryfall data")
	flag.StringVar(&helvaultFile, "helvault", "", "(Re)import helvault library data from csv")
	flag.Parse()

	var err error

	err = db.Open("data/library.db")
	if err != nil {
		logger.Fatal("Cannot open database: ", err)
	}

	// Do the action
	if actionImport {
		confirmPrompt("Scryfall import can take a while and will erase all already imported 'Magic: The Gathering' data.\nProceed?")
		err = importScryfall()
	} else if helvaultFile != "" {
		confirmPrompt("Helvault import will erase all already imported library data.\nProceed?")
		err = helvault.Import(helvaultFile)
	} else if actionServe {
		logger.Info("Starting server listening on ", serveAddr)
		err = web.Serve(serveAddr)
	} else {
		logger.Warn("Invalid call")
		flag.Usage()
		os.Exit(1)
	}

	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Action ended successfully")
}

// run all scryfall imports and log partly progress
func importScryfall() error {
	var err error

	logger.Info("Importing cards")
	if err = scryfall.ImportCards(); err != nil {
		return err
	}
	logger.Info("Success")

	logger.Info("Importing symbols")
	if err = scryfall.ImportSymbols(); err != nil {
		return err
	}
	logger.Info("Success")

	return nil
}

func confirmPrompt(text string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s [y/n]: ", text)

	answer, err := reader.ReadString('\n')
	if err != nil {
		logger.Info("Cannot read user input")
		os.Exit(1)
	}

	answer = strings.ToLower(strings.TrimSpace(answer))

	if !(answer == "y" || answer == "yes") {
		logger.Info("Canceled by user")
		os.Exit(1)
	}
}
