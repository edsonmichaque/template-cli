#!/usr/bin/env bash

sed -i "s|edsonmichaque/template-cli|${1}|g" .goreleaser.yaml

find . -type f -name '*.go' -print0 | xargs -0 sed -i "s|edsonmichaque/template-cli|${1}|g"
find . -type f -name '*.mod' -print0 | xargs -0 sed -i "s|edsonmichaque/template-cli|${1}|g"
find . -type f -name '*.go' -print0 | xargs -0 sed -i "s|template|${2}|g"
find . -type f -name '*.go' -print0 | xargs -0 sed -i "s/TEMPLATE/$(echo "$2" | tr a-z A-Z)/g"

mv cmd/template "cmd/${2}"
