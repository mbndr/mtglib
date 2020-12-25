# MTG Library
> _Magic: The Gathering_ collection viewer / filter with card data from _Scryfall_ and collection data from _Helvault_

```bash
$ # Move Scryfall bulk data (json) and Helvault export (csv) to data/
$ git clone https://github.com/mbndr/mtglib
$ cd mtglib

$ sqlite3 data/library.db
sqlite> .read db/create-schemes.sql
sqlite> .quit

$ go build -o mtglib cmd/main.go
$ ./mtglib -import scryfall -file data/scryfall_cards.json
$ ./mtglib -import helvault -file data/helvault.csv
$ ./mtglib # Run webserver
```

## TODO
- Load other Scryfall data (sets etc)
- Web interface with filters (color, text, type, name, ...)
- Automatic loading and caching of card images
- Paginator on card page?