package iterator

import (
	"log"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"

	"dbstruct"
)

const (
	Limit = 1000
)

func ReadAll(dbmap *gorp.DbMap, limit int) chan []dbstruct.Row {
	result := make(chan []dbstruct.Row, 1)
	go _readAll(dbmap, limit, result)
	return result
}

func _readAll(dbmap *gorp.DbMap, limit int, result chan []dbstruct.Row) {
	lastPoductsID := 0
	lastOffersID := 0

	for {
		var rows []dbstruct.Row

		_, err := dbmap.Select(&rows, dbstruct.SqlLine, lastPoductsID, lastPoductsID, lastOffersID, limit)
		if err != nil {
			log.Println(err)
			close(result)
			return
		}

		if len(rows) == 0 {
			close(result)
			return
		}

		lastPoductsID = rows[len(rows)-1].PoductsID
		lastOffersID = rows[len(rows)-1].OffersID

		result <- rows
	}

	close(result)
}
