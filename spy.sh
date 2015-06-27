#!/bin/sh
(
    ulimit -n 2048
    cd static
    spy -c green --inc "**/*.elm" go generate .
)
