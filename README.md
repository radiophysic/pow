# Proof of Work

## Requirements

see [PRD](docs/PRD.md)

## How to run

TBD 

```shell
docker-compose up -d
```

## Out of scope

This is PoC implementation of **Proof of Work** (PoW) concept.

### Blockchain & storage

In real projects the mined blocks by PoW are stored and distributed over the network.

Also I decided to skip because of PoC:

- Connection encryption (TLS/SSL)
- Graceful shutdown
- Retries strategy / Reconnects
- Penalty for wrong noncens + hash
- fine logger
- config reader

## Decisions

### PoW function - Cuckoo Cycles

see [Cuckoo.pdf](https://github.com/tromp/cuckoo/blob/master/doc/cuckoo.pdf)

The Cuckoo Cycles' algorithm represents memory bandwidth bound concept.
It should resist to ASIC issue. 

### Client-Server message format

`encoding/gob` was selected because it's part of standard library 

## Work in progress

- separate main func at client/server files
- add cli arg (flag)
- clean code
- write tests
- write dockerfiles
- complete readme


