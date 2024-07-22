#!/bin/bash
# data/test/good.yaml will be the config used

# Port 10000
# Default: 42069 -> 0xA455
# Unsigned: 4294967295 -> 0xFFFFFFFF
# SixteenBit: 65535 -> 0x7FFF

# Port 10001
# BigEndian: 1337 -> 0x539
# LittleEndian: 1337 (reversed) -> 0x395

while true; do
    echo -ne "\x00\x00\xA4\x55\xFF\xFF\xFF\xFF\x7F\xFF" | nc -u -w0 localhost 10000
    echo -ne "\x00\x00\x05\x39\x39\x05\x00\x00" | nc -u -w0 localhost 10001

done
