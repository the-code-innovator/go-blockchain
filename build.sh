#!/bin/bash
go clean
rm -rf ./tmp
mkdir -p ./tmp/blocks
go build