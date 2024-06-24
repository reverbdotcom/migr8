package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

var keysProcessed uint64 = 0
var startedAt time.Time

func keyProcessed() {
	// there is no mutex here, but I don't care as this is just information and does not need
	// to be accurate
	keysProcessed += 1
	var duration time.Duration = time.Now().Sub(startedAt)
	kps := float64(keysProcessed) / float64(duration.Seconds())
	log.Printf("\r%v keys processd in %v KPS", keysProcessed, kps)
}

func scanKeys(queue chan Task, wg *sync.WaitGroup) {
	cursor := 0
	conn := sourceConnection(config.Source)

	key_search := fmt.Sprintf("%s*", config.Prefix)
	log.Println("Starting Scan with keys", key_search)

	for {
		// we scan with our cursor offset, starting at 0
		reply, _ := redis.Values(conn.Do("scan", cursor, "match", key_search, "count", config.Batch))

		var tmp_keys []string
		// this func name is confusing...it actually just converts array returns to Go values
		redis.Scan(reply, &cursor, &tmp_keys)

		// put this thing in the queue
		queue <- Task{list: tmp_keys}
		// check if we need to stop...
		if cursor == 0 {
			log.Println("Finished!")

			// close the channel
			close(queue)
			wg.Done()
			break
		}
	}
}
