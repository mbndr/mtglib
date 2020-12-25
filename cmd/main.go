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

func main() {
	importType := flag.String("import", "", "What kind of file to import")
	importFile := flag.String("file", "", "File to import to database")

	serverPort := flag.String("port", ":8080", "Port the server listens on")

	flag.Parse()

	if *importType != "" {
		log.Printf("Imporing %s (%s)", *importFile, *importType)
		err := loadFile(*importType, *importFile)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Successfully imported %s (%s)", *importFile, *importType)
		return
	}

	log.Println("Server listening on " + *serverPort)
	err := mtglib.StartServer(*serverPort)
	if err != nil {
		log.Fatal(err)
	}
}

func loadFile(typ string, path string) error {
	if typ == "meta" {
		return mtglib.ImportMeta()
	}

	if path == "" {
		return errors.New("no import path given")
	}

	r, err := os.Open(path)
	if err != nil {
		return err
	}

	if !confirmPrompt(fmt.Sprintf("This will erase all %s data. Go on?", typ)) {
		return errors.New("import canceled by user")
	}

	if typ == "scryfall" && path != "" {
		return mtglib.ImportScryfall(r)
	} else if typ == "helvault" && path != "" {
		return mtglib.ImportHelvault(r)
	}

	return errors.New("invalid/missing import arguments")
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
