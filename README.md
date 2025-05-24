<div align="center">

  <h1>macinfo client</h1>

  [![Go Version](https://img.shields.io/badge/go-1.20%2B-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

  <sub>Created by <a href="https://github.com/jgengo">Jordane Gengo (Titus)</a></sub><br>
  <sub>From <a href="https://hive.fi">Hive Helsinki</a></sub><br>
  <sub>Highly inspired by <a href="#">macreport</a> (clem)</sub>
</div>

---

## Overview

**macinfo client** is a lightweight Go agent for macOS that periodically collects detailed system information using [osquery](https://osquery.io/) and securely syncs it to a remote API. It is designed for fleet management, asset inventory, and compliance monitoring in organizations, but is also useful for individuals who want deep insight into their Mac.

---

## Features

- Collects rich system data: hostname, active user, UUID, USB devices, uptime, last reboot, temperature sensors, OS version
- Uses [osquery](https://osquery.io/) for robust, extensible data collection
- Sends data as JSON to a configurable API endpoint
- Handles API token management and updates
- Sentry integration for error monitoring
- Configurable sync interval
- Runs as a background service (launchd compatible)
- Simple YAML configuration

---

## Requirements

- macOS 10.13+
- [osquery](https://osquery.io/) installed and running (default socket: `/var/osquery/osquery.em`)
- Go 1.20+ (for development/building)

---

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/jgengo/macinfo-client.git
   cd macinfo-client
   ```
2. **Build the binary:**
   ```bash
   make build
   # or
   go build -o macinfo ./cmd/macinfo
   ```
3. **Copy the config file:**
   ```bash
   cp configs/config.yml /etc/macinfo.yml
   # Or to another location of your choice
   ```

---

## Configuration

The client is configured via a YAML file. Example:

```yaml
osquery_sock: /var/osquery/osquery.em
api_url: http://localhost:3000
api_token: your_api_token_here
sync_interval: 1 # minutes
sentry_dsn: "your_sentry_dsn"
```

**Field explanations:**
- `osquery_sock`: Path to the osquery extension manager socket
- `api_url`: Base URL of the API server to sync data to
- `api_token`: API token for authentication (will be updated if server issues a new one)
- `sync_interval`: How often (in minutes) to sync data
- `sentry_dsn`: (Optional) Sentry DSN for error tracking

You can specify a custom config path at runtime:
```bash
./macinfo -cfg /opt/path/you/want/macinfo.yml
```

---

## Usage

Run the client manually:
```bash
./macinfo
```
Or with a custom config:
```bash
./macinfo -cfg /path/to/macinfo.yml
```

### Running as a Service (launchd)

A sample launchd plist is provided for running MacInfo Client as a background service:

1. Copy the plist:
   ```bash
   sudo cp configs/fi.hive.macinfo.plist /Library/LaunchDaemons/fi.hive.macinfo.plist
   ```
2. Edit the `ProgramArguments` in the plist if your binary/config path differs.
3. Load the service:
   ```bash
   sudo launchctl load /Library/LaunchDaemons/fi.hive.macinfo.plist
   sudo launchctl start fi.hive.macinfo
   ```

---

## Example Output

The client collects and sends data like:
```json
{
  "token": "...",
  "hostname": "my-macbook",
  "active_user": "johndoe",
  "uuid": "C02XXXXXXX",
  "usb_devices": [
    { "vendor": "Apple Inc.", "model": "Apple Keyboard" }
  ],
  "uptime": 123456,
  "last_reboot": ["Sun Apr 19 11:41"],
  "sensors": [
    { "name": "CPU Proximity", "celsius": 54.0 }
  ],
  "os_version": { "version": "13.6.0", "build": "22G120" }
}
```

---

## Development

- Standard Go project layout
- Main entrypoint: `cmd/macinfo/main.go`
- Core logic: `internal/`
- Config: `configs/config.yml`
- Build with `make build` or `go build`

---

## Credits & Inspiration

- Created by [Jordane Gengo (Titus)](https://github.com/jgengo) @ [Hive Helsinki](https://hive.fi)
- Highly inspired by macreport (clem)

---

## License

MIT License. See [LICENSE](LICENSE).

---

## Links

- [GitHub Repository](https://github.com/jgengo/macinfo-client)
- [osquery Documentation](https://osquery.io/docs/)

