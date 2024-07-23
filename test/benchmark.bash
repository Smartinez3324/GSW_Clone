#!/bin/bash

create_packet() {
    local size=$1
    local port=$2
    local timestamp=$(date +%s%6N)

    # Create the packet with UdpSendTimestamp and zeroed ShmSendTimestamp
    # First 8 bytes: UdpSendTimestamp
    # Next 8 bytes: zeroed ShmSendTimestamp
    # Remaining bytes: zeroed
    local packet=$(printf '%016x' $timestamp)
    packet+="0000000000000000"
    packet+=$(head -c $((size - 16)) < /dev/zero | xxd -p | tr -d '\n')

    # Convert packet to binary
    local packet_bin=$(echo $packet | xxd -r -p)

    # Send
    echo -n -e $packet_bin | nc -u -w1 localhost $port
    echo "Sent packet to port $port with size $size bytes"
}

declare -A telemetry_packets=(
    ["10000"]=16
    ["10001"]=17
    ["10002"]=26
    ["10003"]=30
    ["10004"]=38
    ["10005"]=1024
)

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
