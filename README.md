# ChatTCP Capture

A network packet capture tool that provides HTTP APIs for network interface discovery and real-time packet capture using Server-Sent Events (SSE).

## Features

- **Network Interface Discovery**: Query available network interfaces with IPv4 and IPv6 addresses
- **Real-time Packet Capture**: Capture network packets using SSE (Server-Sent Events) for real-time streaming
- **Flexible Filtering**: Filter packets by protocol (TCP/UDP), IP address, and port
- **Automatic Cleanup**: Automatically stops capture when client disconnects
- **RESTful API**: Simple HTTP-based API interface

## Requirements

- Go 1.24.0 or higher
- Linux, macOS, or Windows with libpcap support
- Root/Administrator privileges (required for packet capture)
- **Dependencies:**
  - `github.com/gopacket/gopacket v1.5.0` - Packet capture and decoding library
  - `github.com/gin-gonic/gin v1.11.0` - HTTP web framework

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd chattcp-capture
```

2. Install dependencies:
```bash
go mod download
```

3. Build the server:
```bash
go build -o chattcp-capture cmd/server/main.go
```

## Quick Start

### Run with default port (8080):
```bash
go run cmd/server/main.go
```

### Run with custom port:
```bash
go run cmd/server/main.go -port 9090
```

### Run the compiled binary:
```bash
./chattcp-capture -port 8080
```

The server will start and listen on the specified port. You should see:
```
Server starting on :8080
```

## API Documentation

### 1. List Network Interfaces

Get a list of all available network interfaces.

**Endpoint:** `GET /api/interfaces`

**Response:**
```json
{
  "data": [
    {
      "Name": "en0",
      "IpV4Address": "192.168.1.100",
      "IpV6Address": "fe80::1"
    }
  ]
}
```

**Example:**
```bash
curl http://localhost:8080/api/interfaces
```

### 2. Start Packet Capture (SSE)

Start capturing packets and receive them via Server-Sent Events.

**Endpoint:** `GET /api/capture`

**Query Parameters:**
- `i` (required): Network interface name (e.g., "en0", "eth0")
- `p` (optional): Protocol filter - "tcp" or "udp", or port number for any port filter
- `h` (optional): Filter by any IP address (source or destination)
- `h.src` (optional): Filter by source IP address
- `h.dst` (optional): Filter by destination IP address
- `p.src` (optional): Filter by source port
- `p.dst` (optional): Filter by destination port

**Note:** The `p` parameter can be either:
- A protocol string ("tcp" or "udp") to filter by protocol
- A port number to filter by any port (source or destination)
- If both protocol and port filtering are needed, use `p` for protocol and `p.src`/`p.dst` for port

**Response:** Server-Sent Events stream

**Event Types:**
- `packet`: Contains captured packet data
- `error`: Error occurred during capture
- `close`: Capture stopped

**Example:**
```bash
# Capture TCP packets on port 8080
curl "http://localhost:8080/api/capture?i=en0&p=tcp&p.src=8080"

# Capture packets from specific source IP
curl "http://localhost:8080/api/capture?i=en0&h.src=192.168.1.100"

# Capture UDP packets on destination port 53
curl "http://localhost:8080/api/capture?i=en0&p=udp&p.dst=53"
```

## Packet Data Structure

Each captured packet is returned as JSON with the following structure:

```json
{
  "timestamp": 1234567890123,
  "packet_size": 1500,
  "src_ip": "192.168.1.100",
  "dst_ip": "192.168.1.1",
  "proto": "tcp",
  "tcp": {
    "src": 8080,
    "dst": 443,
    "seq": 123456,
    "ack": 789012,
    "data_offset": 5,
    "FIN": false,
    "SYN": true,
    "RST": false,
    "PSH": false,
    "ACK": false,
    "URG": false,
    "ECE": false,
    "CWR": false,
    "NS": false,
    "window": 65535,
    "checksum": 12345,
    "urgent": 0,
    "payload": [],
    "options": []
  },
  "udp": null
}
```

For UDP packets, the `tcp` field will be `null` and `udp` will contain:
```json
{
  "src": 5353,
  "dst": 53,
  "length": 100,
  "checksum": 12345,
  "payload": []
}
```

## Command Line Options

- `-port`: Specify the server port (default: 8080)
  ```bash
  ./chattcp-capture -port 9090
  ```

## Notes

1. **Permissions**: Packet capture requires root/administrator privileges on most systems. Run with appropriate permissions:
   ```bash
   sudo ./chattcp-capture
   ```

2. **Network Interface Names**: 
   - Linux: Usually `eth0`, `wlan0`, etc.
   - macOS: Usually `en0`, `en1`, etc.
   - Windows: Use interface names from `ipconfig` or `ipconfig /all`

3. **Performance**: For high-traffic networks, consider using more specific filters to reduce the number of captured packets.

4. **Connection Management**: The capture automatically stops when the client disconnects. Make sure to properly close SSE connections to free resources.

5. **CORS**: The server includes CORS headers allowing cross-origin requests. Adjust CORS settings in `cmd/server/main.go` if needed for production.

## Development

### Project Structure

```
chattcp-capture/
├── cmd/
│   └── server/
│       └── main.go          # HTTP server entry point
├── api/
│   └── api.go               # API handlers
├── capture.go               # Core capture logic
├── network_interfaces.go    # Network interface discovery
├── types.go                 # Type definitions
├── conver.go                # Packet conversion utilities
└── go.mod                   # Go module definition
```

### Running Tests

```bash
go test ./...
```

## License

This project is licensed under the MIT License.

### Third-party Libraries

This project uses the following third-party libraries:

- **gopacket** (BSD-3-Clause License) - https://github.com/gopacket/gopacket v1.5.0

## Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork the repository** and create your feature branch (`git checkout -b feature/amazing-feature`)
2. **Make your changes** following the existing code style
3. **Add tests** for new functionality if applicable
4. **Ensure all tests pass** (`go test ./...`)
5. **Commit your changes** with clear commit messages
6. **Push to your branch** (`git push origin feature/amazing-feature`)
7. **Open a Pull Request** with a detailed description of your changes

### Code Style

- Follow Go standard formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and concise

### Reporting Issues

When reporting issues, please include:
- Operating system and version
- Go version
- Steps to reproduce the issue
- Expected behavior vs actual behavior
- Any relevant error messages or logs
