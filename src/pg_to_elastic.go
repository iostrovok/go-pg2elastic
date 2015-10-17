package main

import (
	_ "github.com/lib/pq"

	"dbihelper"
	"dbstruct"
	"iterator"
	"manager"
	"save2elastic"
)

// limitReadOneTime is count of reading from DB for one time
// If you have problem with postgresql performance you need to decrease it.
const limitReadOneTime = 100

func main() {

	dbmap := dbihelper.GetDBI()

	// Start manager
	toElasticCh := make(chan dbstruct.Row, 1)
	rowCh := make(chan dbstruct.Row, 1)
	manager.Start(toElasticCh, rowCh)

	// Start saving to elasticserach
	finishElasticSaveCh := save2elastic.Start(toElasticCh)

	// Start reading from DB
	rowsCh := iterator.ReadAll(dbmap, limitReadOneTime)
	readDBtoMAnager(rowCh, rowsCh)
	close(rowCh)

	// wait while elasticsearch are getting data
	<-finishElasticSaveCh
}

func readDBtoMAnager(rowCh chan dbstruct.Row, rowsCh chan []dbstruct.Row) {
	for {
		select {
		case rows, ok := <-rowsCh:
			if !ok {
				return
			}
			for _, oneRow := range rows {
				rowCh <- oneRow
			}
		}
	}
}
