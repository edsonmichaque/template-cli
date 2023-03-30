#!/usr/bin/env sh

if [ ! -d "e2e/bats" ]; then
    git clone https://github.com/bats-core/bats e2e/bats
fi

if [ ! -d "e2e/test_helper/bats-assert" ]; then
    git clone https://github.com/bats-core/bats-assert e2e/test_helper/bats-assert
fi

if [ ! -d "e2e/test_helper/bats-support" ]; then
    git clone https://github.com/bats-core/bats-support e2e/test_helper/bats-support
fi
