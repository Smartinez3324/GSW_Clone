#!/bin/bash
# data/test/good.yaml will be the config used

# Port 10000
# Default: Funny numbers
# Unsigned: 4294967295
# SixteenBit: 65535

# Port 10001
# BigEndian: 1337
# LittleEndian: 1337 (reversed)

while true; do
    echo -ne "\x0A\x04\x05\x05
              \xFF\xFF\xFF\xFF
              \xFF\xFF

              " | nc -u -w0 localhost 10000
    echo -ne "\x00\x05\x03\x09
              \x09\x03\x05\x00
              " | nc -u -w0 localhost 10001

done
