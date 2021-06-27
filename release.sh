#!/bin/bash
go mod tidy
go test -v -race ./...
bash ./misc/scripts/bumpversion.sh
goreleaser --rm-dist #--snapshot --skip-publish
