CREATE TABLE IF NOT EXISTS helvault_library (
    scryfall_id VARCHAR(255) PRIMARY KEY,
    quantity INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS scryfall_cards (
	scryfall_id VARCHAR(255) PRIMARY KEY,
	oracle_id VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
	image_uri VARCHAR(512),
	mana_cost VARCHAR(255),
	cmc FLOAT NOT NULL,
	type_line VARCHAR(255) NOT NULL,
	oracle_text VARCHAR(512),
	colors VARCHAR(255),
	color_identity VARCHAR(255) NOT NULL,
	set_code VARCHAR(10) NOT NULL,
	set_name VARCHAR(255) NOT NULL,
	rarity VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS scryfall_card_faces (
	card_id VARCHAR(255) NOT NULL,
	colors VARCHAR(255),
	image_uri VARCHAR(512),
	mana_cost VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
	type_line VARCHAR(255) NOT NULL,
	FOREIGN KEY(card_id) REFERENCES scryfall_cards(scryfall_id)
);

CREATE TABLE IF NOT EXISTS scryfall_symbols (
	symbol VARCHAR(64) PRIMARY KEY,
	svg_uri VARCHAR(255) NOT NULL,
	title VARCHAR(255)
);