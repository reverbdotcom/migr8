migr8
---

Redis Migration Utility written in Go

## Build
migr8 uses [gb](http://getgb.io) to vendor dependencies.

To install it run, `go get github.com/constabulary/gb/...`

Tests require that `redis-server` is somewhere in your $PATH.

`make` To run tests and create a binary

## Usage
```
NAME:
   migr8 - It's time to move some redis

USAGE:
   migr8 [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   migrate	Migrate one redis to a new redis
   delete	Delete all keys with the given prefix
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dry-run, --dr			Run in dry-run mode
   --source, -s "127.0.0.1:6379"	The redis server to pull data from
   --dest, -d "127.0.0.1:6379"		The destination redis server
   --workers, -w "2"			The count of workers to spin up
   --batch, -b "10"			The batch size
   --prefix, -p 			The key prefix to act on
   --clear-dest, -c			Clear the destination of all it's keys and values
   --help, -h				show help
   --version, -v			print the version
```

#### Cross Compile for Linux:
*Note:* You will need the Go cross compile tools. If you're using homebrew: `brew install go --cross-compile-common`

`make linux` will build a linux binary in bin/
