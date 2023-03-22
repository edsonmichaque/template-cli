#!/usr/bin/env bash

sed -i "s|edsonmichaque/tempate-cli|${1}|g" .gorelease.yml

find . -type f -name '*.go' -print0 | xargs -0 sed -i "s|edsonmichaque/tempate-cli|${1}|g"
find . -type f -name '*.mod' -print0 | xargs -0 sed -i "s|edsonmichaque/template-cli|${1}|g"
find . -type f -name '*.go' -print0 | xargs -0 sed -i "s|tempate|${2}|g"

mv cmd/template "cmd/${2}"
