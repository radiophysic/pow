# Proof of Work

## Requirements

see [PRD](docs/PRD.md)

## How to run

0. Get [docker-compose](https://docs.docker.com/compose/install/)

1. Build docker image 

```shell
docker-compose up -d
```

2. Check containers' output:

```shell
docker-compose logs
```

## Out of scope

This is PoC implementation of **Proof of Work** (PoW) concept.

### Blockchain & storage

In real projects the mined blocks by PoW are stored and distributed over the network.

I decided to skip because of PoC:

- Connection encryption (TLS/SSL)
- Graceful shutdown
- Retries strategy / Reconnects
- Penalty for wrong noncens + hash
- Unit tests for client (I managed to write tests for service handlers)

## Decisions

### PoW function - Cuckoo Cycles

see [Cuckoo.pdf](https://github.com/tromp/cuckoo/blob/master/doc/cuckoo.pdf)

The Cuckoo Cycles' algorithm represents memory bandwidth bound concept.
It should resist to ASIC issue. 

### Client-Server message format

`encoding/gob` was selected because it's part of standard library 
