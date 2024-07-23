import socket
import threading
import time
import struct

TELEMETRY_PACKETS = [
    {"name": "TimestampsOnly", "port": 10000, "size": 16},
    {"name": "17Bytes", "port": 10001, "size": 17},
    {"name": "26Bytes", "port": 10002, "size": 26},
    {"name": "30Bytes", "port": 10003, "size": 30},
    {"name": "38Bytes", "port": 10004, "size": 38},
    {"name": "OneKbyte", "port": 10005, "size": 1024},
]

SERVER_ADDRESS = 'localhost' 

def create_packet(size):
    timestamp = int(time.time() * 1000)
    # UdpSendTimestamp (8 bytes), ShmSendTimestamp (8 bytes zeroed), followed by zeroed remaining bytes
    packet = struct.pack('>Q', timestamp) + b'\x00' * (size - 8)
    return packet

def send_packet(port, size):
    with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as sock:
        server_address = (SERVER_ADDRESS, port)
        while True:
            packet = create_packet(size)
            sock.sendto(packet, server_address)

def run_test():
    while True:
        threads = []
        for packet in TELEMETRY_PACKETS:
            t = threading.Thread(target=send_packet, args=(packet["port"], packet["size"]))
            t.start()
            threads.append(t)

        for t in threads:
            t.join()

if __name__ == "__main__":
    run_test()
