package manager

import (
	// rand is using for elasticsearch condition imitation
	"math/rand"

	"log"
	"sync"
	"time"

	"dbstruct"
)

const (
	checkTTL = 1 * time.Second
)

var sleepLock = &sync.RWMutex{}

func Start(toElasticCh, rowCh chan dbstruct.Row) {
	stopSleepCh := make(chan bool, 1)
	go fromDBtoElastic(toElasticCh, rowCh, stopSleepCh)
	go SleepTime(stopSleepCh)
}

func fromDBtoElastic(toElasticCh, rowCh chan dbstruct.Row, stopSleepCh chan bool) {
	for row := range rowCh {
		sleepLock.RLock()
		toElasticCh <- row
		sleepLock.RUnlock()
	}
	close(toElasticCh)

	// Stop sleep goroutine
	stopSleepCh <- true
}

/*
	If we need to make a break in our work we have to define that here.
	For example: we can check state of the elasticsearch
	cluster (with web or other interface) and stop work if
	there are "yellow/red" nodes.
*/
func SleepTime(stopSleepCh chan bool) {

	// state == true sleepLock is LOCK here
	// state == false sleepLock is UNLOCK here
	state := false
	tiker := time.Tick(checkTTL)

	for {
		select {
		case <-tiker:
			elastciNow := is_elastic_claster_good()
			if !elastciNow && !state {
				sleepLock.Lock()
				state = true
				log.Println("waiting good conditions of elasticsearch")
			} else if elastciNow && state {
				log.Println("elasticsearch can get records")
				sleepLock.Unlock()
				state = false
			}
		case <-stopSleepCh:
			log.Println("Finish sleep")
			return
		}
	}

}

// imitation of check condition of  elasticsearch cluster
func is_elastic_claster_good() bool {
	if rand.Intn(100) > 50 {
		return true
	}
	return false
}
