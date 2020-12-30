# MTG Library
> _Magic: The Gathering_ collection viewer / filter with card data from _Scryfall_ and collection data from _Helvault_

```bash
$ git clone https://github.com/mbndr/mtglib
$ cd mtglib
$ # move Helvault export (csv) to data/

$ sqlite3 data/library.db
sqlite> .read db/create-schemes.sql
sqlite> .quit

$ go build -o mtglib cmd/main.go
$ ./mtglib -import # import scryfall data
$ ./mtglib -helvault -file data/helvault.csv # import your helvault collection
$ ./mtglib -serve # run webserver
```

## TODO
- Add detail view for multiface cards
- Color filter: Multicolor only
- Code (and comment) refactoring
- Adding cards to collections (and export them to csv / json)
- Use color-identity for filtering also?
- When finished, provide release?