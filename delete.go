package main

import (
	"log"
	"sync"

	"github.com/garyburd/redigo/redis"
)

func deleteKeys(queue chan Task, wg *sync.WaitGroup) {
	sourceConn := sourceConnection(config.Source)
	for task := range queue {
		for _, key := range task.list {
			if config.DryRun {
				log.Printf("Would have deleted %s", key)
				continue
			}

			if _, err := redis.String(sourceConn.Do("del", key)); err != nil {
				log.Printf("Deleted %s \n", key)
			} else {
				log.Printf("Could not deleted %s: %s\n", key, err)
			}
		}
	}

	wg.Done()
}
