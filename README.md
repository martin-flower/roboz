# roboz cleaner

example restful implementation of developer [exercise](exercise.md)

implemented in three containers:

1. **database** - implemented using postgres
1. **web** - restful http server implemented in go & fiber
1. **smoke** - post deploy end-to-end smoke tests (container runs and exits) - see [README.md](smoke/README.md)

## configuration

create a file `database/.env` containing for example

```
POSTGRES_PASSWORD=5dHEzcFukeMRacMWrmbd
POSTGRES_USER=robozuser
POSTGRES_DB=robozdb
```

add this file to .gitignore

## ports
- database port 6432
- http server port 5000

on **macos**, port 5000 is used for airplay receiver - free up this port by using system preferences to deactivate airplay receiver (only applies to recent models - option visible in system preferences)

## starting

* in root directory, `docker compose up`

## stopping

* in root directory, `docker compose down`

## landing page

[localhost:5000](http://localhost:5000) documents the rest endpoints

## hardware requirements

* test data scope requires a machine with at least 16GB memory

## comments

* the core algorithm has four alternative implementations - 
  * [intersection](web/service/clean/intersection/README.md) 
  * [intmap](web/service/clean/intmap/clean.go) (the one to use) 
  * [simplest](web/service/clean/simplest/clean.go)
  * [sortedset](web/service/clean/sortedset/clean.go)

* note that **none of the algorithms performs adequately at scale** - intmap is the one that goes the farthest

### TODO

* database - create application users - see [init.sql](database/init.sql)
* in order to move from just being an exercise to more of a proof-of-concept, it would be best to shield the executions from the rest server. The rest server should accept a small number of pending requests, it should acknowledge the request immediately, and then queue the request for the executor. The person sending the request should either receive a acknowledgement with a pending reference number, or should receive message that the service is busy (come back later).
* I also looked into some simple [performance](performance/README.md) tests, but given that the executor is so slow, these have little use until we can decouple the rest server from the executor.