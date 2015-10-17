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
const limitReadOneTime = 10

func main() {

	dbmap := dbihelper.GetDBI()

	// Start manager
	toElasticCh := make(chan dbstruct.Row, 1)
	rowsCh := make(chan []dbstruct.Row, 1)
	manager.Start(toElasticCh, rowsCh)

	// Start saving to elasticserach
	finishElasticSaveCh := save2elastic.Start(toElasticCh)

	// Start reading from DB
	iterator.ReadAll(dbmap, rowsCh, limitReadOneTime)

	// wait while elasticsearch are getting data
	<-finishElasticSaveCh
}
