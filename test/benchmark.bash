#!/bin/bash

# Function to create a packet
create_packet() {
    local size=$1
    local port=$2
    local timestamp=$(date +%s%6N) # Current timestamp in microseconds

    local packet=$(printf '%016x' $timestamp)
    packet+="0000000000000000"
    packet+=$(head -c $((size - 16)) < /dev/zero | tr '\0' '\x01' | xxd -p | tr -d '\n')

    local packet_bin=$(echo $packet | xxd -r -p)

    printf "$packet_bin" | nc -u -w1 localhost $port
    echo "Sent packet to port $port with size $size bytes"
}

# Array of telemetry packet configurations
declare -A telemetry_packets=(
    ["10000"]=16
    ["10001"]=17
    ["10002"]=26
    ["10003"]=30
    ["10004"]=38
    ["10005"]=1024
)

# Function to send packets
send_packets() {
    while true; do
        for port in "${!telemetry_packets[@]}"; do
            create_packet "${telemetry_packets[$port]}" "$port" &
        done
        wait
    done
}

# Run the send_packets function
send_packets
