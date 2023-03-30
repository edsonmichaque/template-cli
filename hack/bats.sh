#!/usr/bin/env sh

if [ ! -d "test/bats" ]; then
    git clone https://github.com/bats-core/bats test/bats
fi

if [ ! -d "test/test_helper/bats-assert" ]; then
    git clone https://github.com/bats-core/bats-assert test/test_helper/bats-assert
fi

if [ ! -d "test/test_helper/bats-support" ]; then
    git clone https://github.com/bats-core/bats-support test/test_helper/bats-support
fi
