#!/bin/bash

while true; do
    echo -ne "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B" | nc -u -w0 localhost 10000
done
