#!/bin/bash

ver=$(cd twemoji; git describe --tags --abbrev=0)
output=twemoji.go

go run script/gen.go "${ver:1}" > $output
gofmt -w $output
