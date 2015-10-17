package iterator

import (
	"log"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"

	"dbstruct"
)

const (
	SqlFull = `
SELECT 
	products.id as products_id, products.title as products_title,
	COALESCE(offers.id, 0) as offers_id, COALESCE(offers.title, '') as offers_title 
FROM products LEFT JOIN offers ON offers.products_id = products.id
WHERE $1 < products.id AND products.id <= $2
ORDER by products
`

	SqlProductId = `
SELECT products.id FROM products 
WHERE products.id > $1 ORDER BY products.id LIMIT 1 OFFSET $2
`
)

func ReadAll(dbmap *gorp.DbMap, result chan []dbstruct.Row, limit int64) {
	go _readAll(dbmap, result, limit)
}

func _readAll(dbmap *gorp.DbMap, result chan []dbstruct.Row, limit int64) {
	var lastPoductsID int64 = 0
	var lastIter bool = false

	defer close(result)

	for {
		var rows []dbstruct.Row

		// Get
		nextPoductsID, err := dbmap.SelectInt(SqlProductId, lastPoductsID, limit)
		if err != nil {
			log.Println(err)
			return
		}

		// we have to read last rows
		if nextPoductsID == 0 {
			nextPoductsID = lastPoductsID + limit
		}

		_, err = dbmap.Select(&rows, SqlFull, lastPoductsID, nextPoductsID)
		if err != nil {
			log.Println(err)
			return
		}

		if len(rows) == 0 || lastIter {
			log.Println(err)
			return
		}

		lastPoductsID = nextPoductsID

		result <- rows
	}
}
