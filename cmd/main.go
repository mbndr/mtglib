package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mbndr/mtglib"
)

// TODO: HERE: interface idea, import scryfall data
// für leichten import: Keine direkten Fremdschlüssel, sondern scryfall_id, dann kann geschaut werden, ob karte vorhanden.
//   zb Kartensuche: Dann zeichen ob schon in sammlung
//   gut, weil helvault csv kann ohne probleme gelöscht und wieder importiert werden (loose coupling)
// Für Bilder: bei aufruf (/image/card/{scryfallId}) schauen ob schon heruntergeladen, ansonsten runterladen und dann anzeigen (TEST!)

// Card get from db, for general use (joined data)
type Card struct {
	ScryfallID  string
	ScryfallURI string
	Languages   []string
	Colors      []rune
	Sets        []string
	// ...??
}

func main() {
	importType := flag.String("import", "", "What kind of file to import")
	importFile := flag.String("file", "", "File to import to database")

	flag.Parse()

	if *importType != "" && *importFile != "" {
		log.Printf("Imporing %s (%s)", *importFile, *importType)
		err := loadFile(*importType, *importFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Successfully imported %s (%s)", *importFile, *importType)
		return
	}

	fmt.Println("serve")
}

func loadFile(typ string, path string) error {
	r, err := os.Open(path)
	if err != nil {
		return err
	}

	if !confirmPrompt(fmt.Sprintf("This will erase all %s data. Go on?", typ)) {
		return errors.New("import canceled by user")
	}

	if typ == "scryfall" {
		return mtglib.LoadScryfall(r)
	} else if typ == "helvault" {
		return mtglib.LoadHelvault(r)
	}

	return errors.New("invalid import type")
}

func confirmPrompt(text string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s [y/n]: ", text)

	answer, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	answer = strings.ToLower(strings.TrimSpace(answer))

	if answer == "y" || answer == "yes" {
		return true
	}

	return false
}
