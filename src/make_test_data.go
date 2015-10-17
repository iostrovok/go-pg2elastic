package main

import (
	//"database/sql"
	"fmt"
	"math/rand"

	"gopkg.in/gorp.v1"

	"dbihelper"
)

const (
	countTestProducts = 1000000

	productsSql = "INSERT INTO products(title) VALUES ($1) RETURNING id"
	offersSql   = "INSERT INTO offers(title, products_id) VALUES ($1, $2)"
	dropTables  = `
DROP TABLE IF EXISTS offers;
DROP TABLE IF EXISTS products;
`
	productsSqlCreate = ` 
CREATE TABLE products (id serial NOT NULL, title text, CONSTRAINT products_pk_id PRIMARY KEY (id)) WITH (OIDS=FALSE);
`

	offersSqlCreate = `
CREATE TABLE offers
(
  id serial NOT NULL,
  title text,
  products_id integer,
  CONSTRAINT offers_pkey PRIMARY KEY (id),
  CONSTRAINT products_id_offers_id_fk FOREIGN KEY (products_id)
      REFERENCES products (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
)
WITH (
  OIDS=FALSE
);
`
)

func main() {
	dbmap := dbihelper.GetDBI()
	create_tables(dbmap)
	load_test_data(dbmap, countTestProducts)
}

func create_tables(dbmap *gorp.DbMap) {

	if _, err := dbmap.Exec(dropTables); err != nil {
		panic(err)
	}

	if _, err := dbmap.Exec(productsSqlCreate); err != nil {
		panic(err)
	}

	if _, err := dbmap.Exec(offersSqlCreate); err != nil {
		panic(err)
	}
}

func load_test_data(dbmap *gorp.DbMap, count int) {
	for count > 0 {
		count--

		var lastProductsId int
		title := randString()
		err := dbmap.Db.QueryRow(productsSql, title).Scan(&lastProductsId)
		if err != nil {
			panic(err)
		}

		line := fmt.Sprintf("%d. lastInsertId: %d\n", count, lastProductsId)
		countOffers := rand.Intn(100)
		for countOffers > 0 {
			countOffers--
			title := randString()
			_, err := dbmap.Exec(offersSql, title, lastProductsId)
			if err != nil {
				panic(err)
			}
			line += fmt.Sprintf("+")

		}

		if count%100 == 0 {
			fmt.Println(line)
		}
	}
}

func randString() string {
	line := []string{"q", "a", "z", "w", "s", "x", "e", "d", "c", "r", "f", "v", "t",
		"g", "b", "y", "h", "n", "u", "j", "m", "i", "k", "o", "l", "p", "Q", "A", "Z",
		"W", "S", "X", "E", "D", "C", "R", "F", "V", "T", "G", "B", "Y", "H", "N", "U",
		"J", "M", "I", "K", "O", "L", "P", " ", " ", " ", " ", " ", " ", " ", " ", " ",
	}
	l := 3 + rand.Intn(16)
	c := len(line)
	out := ""
	for l > 0 {
		l--
		x := rand.Intn(c)
		out += line[x]
	}

	return out
}
