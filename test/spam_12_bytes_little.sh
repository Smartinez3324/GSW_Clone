#!/bin/bash

while true; do
    echo -ne "\x03\x02\x01\x00\x07\x06\x05\x04\x0B\x0A\x09\x08" | nc -u -w0 localhost 10000
done
