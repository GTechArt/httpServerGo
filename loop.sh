#!/bin/bash

for i in $(seq 1 42)
do
	curl -s http://localhost:8080/app/ -o /dev/null
done
