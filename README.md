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
- Add rarity (card, filter, sorting)
- Add config for sort order
- Add detail view for multiface cards