package main

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
)

func Test_MigrateAllKeysWithAPrefix(t *testing.T) {
	ClearRedis()

	config = Config{
		Source:  sourceServer.url,
		Dest:    destServer.url,
		Workers: 1,
		Batch:   10,
		Prefix:  "bar",
	}

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("bar:%d", i)
		sourceServer.conn.Do("SET", key, i)
	}

	sourceServer.conn.Do("SET", "baz:foo", "yolo")

	RunAction(migrateKeys)

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("bar:%d", i)
		exists, _ := redis.Bool(destServer.conn.Do("EXISTS", key))

		if !exists {
			t.Errorf("Could not find a key %s that should have been migrated", key)
		}
	}

	exists, _ := redis.Bool(destServer.conn.Do("EXISTS", "baz:foo"))

	if exists {
		t.Errorf("Found a key %s that should not have been migrated", "baz:foo")
	}
}

func Test_MigrateAllKeysWithTTLs(t *testing.T) {
	ClearRedis()

	config = Config{
		Source:  sourceServer.url,
		Dest:    destServer.url,
		Workers: 1,
		Batch:   10,
		Prefix:  "bar",
	}

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("bar:%d", i)
		sourceServer.conn.Do("SET", key, i, "EX", 600)
	}

	RunAction(migrateKeys)

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("bar:%d", i)
		exists, _ := redis.Bool(destServer.conn.Do("EXISTS", key))

		if !exists {
			t.Errorf("Could not find a key %s that should have been migrated", key)
		}

		ttl, _ := redis.Int64(destServer.conn.Do("PTTL", key))

		if ttl < 1 || ttl > 600000 {
			t.Errorf("Could not find a TTL for key %s that should have been migrated", key)
		}
	}
}

func Test_DoesNothingInDryRunModeForMigrate(t *testing.T) {
	ClearRedis()

	config = Config{
		Source:  sourceServer.url,
		Workers: 1,
		Batch:   10,
		Prefix:  "bar",
		DryRun:  true,
		Dest:    destServer.url,
	}

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("bar:%d", i)
		sourceServer.conn.Do("SET", key, i)
	}

	RunAction(migrateKeys)

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("bar:%d", i)
		exists, _ := redis.Bool(destServer.conn.Do("EXISTS", key))

		if exists {
			t.Errorf("In DryRun mode, but found a key %s that was actually migrated", key)
		}
	}
}
