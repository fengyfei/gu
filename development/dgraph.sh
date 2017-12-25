#!/bin/sh
docker exec -it -d dgraph dgraph server --bindall=true --memory_mb 2048 --zero localhost:5080
docker exec -it -d dgraph dgraph-ratel
