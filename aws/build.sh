#!/bin/sh

GOOS=linux go build main.go
zip demo.zip ./main