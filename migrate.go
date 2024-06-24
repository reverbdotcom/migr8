package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/garyburd/redigo/redis"
)

func dumpKeyAndTTL(key string, sourceConn redis.Conn) (string, int64, error) {
	var err error
	var dumpedKey string
	var ttl int64

	if dumpedKey, err = redis.String(sourceConn.Do("dump", key)); err != nil {
		return dumpedKey, ttl, err
	}

	if ttl, err = redis.Int64(sourceConn.Do("pttl", key)); err != nil {
		return dumpedKey, ttl, err
	}

	return dumpedKey, ttl, err
}

func dumpAndRestore(sourceConn redis.Conn, destConn redis.Conn, key string) {
	dumpedKey, dumpedKeyTTL, err := dumpKeyAndTTL(key, sourceConn)

	if err != nil {
		log.Println(err)
		return
	}

	// when doing pttl, -1 means no expiration
	// when doing restore, 0 means no expiration
	if dumpedKeyTTL == -1 {
		dumpedKeyTTL = 0
	}

	if config.DryRun {
		log.Printf("Would have restored %s with ttl %d", key, dumpedKeyTTL)
		return
	}
	_, err = destConn.Do("restore", key, dumpedKeyTTL, dumpedKey)

	if err != nil {
		log.Printf("error: %s\n", err)
		return
	}

	keyProcessed()
}

func migrateKeys(queue chan Task, wg *sync.WaitGroup) {
	sourceConn := sourceConnection(config.Source)
	destConn := destConnection(config.Dest)

	for task := range queue {
		for _, key := range task.list {
			dumpAndRestore(sourceConn, destConn, key)
		}
	}

	wg.Done()
}

func shouldClearAllKeys(dest string) bool {
	fmt.Println("Are you sure you want to delete all keys at", dest, "? . Please type Y or N.")

	var response string
	if _, err := fmt.Scanln(&response); err == nil {
		return response == "Y"
	}

	return false
}

func clearDestination(dest string) {
	if shouldClearAllKeys(dest) {
		log.Println("Deleting all keys of destination")
		destConn := destConnection(dest)

		if _, err := destConn.Do("flushall"); err != nil {
			log.Printf("error in flushing: %s\n", err)
		}
	} else {
		log.Println("Skipping key deletion")
	}
}
