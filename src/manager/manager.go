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
	checkTTL = 2 * time.Second
)

var sleepLock = &sync.RWMutex{}

func Start(toElasticCh chan dbstruct.Row, rowsCh chan []dbstruct.Row) {
	stopSleepCh := make(chan bool, 1)
	go fromDBtoElastic(toElasticCh, rowsCh, stopSleepCh)
	go SleepTime(stopSleepCh)
}

func fromDBtoElastic(toElasticCh chan dbstruct.Row, rowsCh chan []dbstruct.Row, stopSleepCh chan bool) {

	for rows := range rowsCh {
		for _, oneRow := range rows {
			sleepLock.RLock()
			toElasticCh <- oneRow
			sleepLock.RUnlock()
		}
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
			elastciNow := is_elastic_cluster_good()
			if !elastciNow && !state {
				sleepLock.Lock()
				state = true
				log.Printf("[%s] waiting good conditions of elasticsearch\n", time.Now().Format("2006-01-02T15:04:05.999999999"))
			} else if elastciNow && state {
				log.Printf("[%s] elasticsearch can get records\n", time.Now().Format("2006-01-02T15:04:05.999999999"))
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
func is_elastic_cluster_good() bool {
	/*
		If we need to check state of elastic cluster
		we can send "http://elastic:9200/_cluster/health" or
		"http://elastic:9200/_cat/indices" and parse the answer.
	*/
	if rand.Intn(100) > 10 {
		return true
	}
	return false
}
