package save2elastic

/*
	Here:
	1) save to elasticsearch
	2) report about finish work
*/

import (
	"fmt"
	"log"

	"dbstruct"
)

func Start(toElasticCh chan dbstruct.Row) chan bool {
	allDoneCh := make(chan bool, 1)
	go putToElastic(toElasticCh, allDoneCh)
	return allDoneCh
}

func putToElastic(toElasticCh chan dbstruct.Row, allDoneCh chan bool) {
	for {
		select {
		case row, ok := <-toElasticCh:
			if !ok {
				log.Println(`"Send to elastic" finished`)
				allDoneCh <- true
				return
			}
			/*
				Here we have to define the saving into elasticsearch
			*/
			fmt.Printf("Send to elastic: %d _ %d\n", row.PoductsID, row.OffersID)
		}

	}
}
