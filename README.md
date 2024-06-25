# Zia

### A simple reverse proxy

version: 0.3.2

Zia is a reverse proxy written in Go language. It was created to provide a simple and fast way to access old Docker containers.

- 0.3.2 fix
  -ssl=false start without frontend ssl
  -stdout=true redirect stdout and stderr to console instead logfile
- 0.3.1 simplified the execution without configuration file just params.

#### Features

- Reverse proxy
- SSL/TLS frontend support with Let's Encrypt certificates
- Round-robin load balancing
- SSL/TLS targets with self-signed or CA signed
- https->https targets, https->http targets, http->http targets, (http|https)->mixed targets
- Configurable timeout for proxy connections
- Access logging on logfiles or stdout

#### Requirements

- Go installed
- IPv4 network

#### Installing Zia

1. Build the code with `go build zia.go` or `make zia`
2. Install `make install`
3. If below 1024 ports required, in Linux run `sudo setcap CAP_NET_BIND_SERVICE=+eip zia` to allow Zia to bind to low ports
4. If needer there's a testTargetServer for testing zia

#### Parameters

- `-domain`: The domain name for which the reverse proxy will be set up. This is required for obtaining a Let's Encrypt SSL/TLS certificate.
- `-port`: The port on which the reverse proxy will listen. Default is 443 for HTTPS, 80 for HTTP.
- `-ssl`: Enable or disable SSL/TLS. If set to `true` (default), Zia will attempt to obtain a Let's Encrypt certificate for the specified domain. If set to `false`, Zia will run in HTTP mode without SSL/TLS.
- `-targets`: A comma-separated list of target URLs to which the reverse proxy will forward requests. These can be HTTP or HTTPS URLs.
- `-timeout`: The read timeout for proxy connections, in seconds. If not provided or set to 0, there is no timeout.

#### Run Zia

Run the binary with the desired options:

This will start Zia as a reverse proxy on port 443, with SSL/TLS enabled using a Let's Encrypt certificate for the domain `example.com`. It will load balance requests between the targets `https://127.0.0.1:8020` and `https://127.0.0.1:8021` using a round-robin algorithm.

##### zia for example.com port 8080 with 5 second timeout and https from let's encrypt

```bash
zia -domain example.com -targets https://127.0.0.1:8020,https://10.0.0.10:8021 -timeout 5
```

##### zia for example.com without ssl but targets has own signed certificate, stdout logs

```bash
zia -domain example.com -ssl=false -stdout=false -targets https://127.0.0.1:8020,https://10.0.0.10:8021
```

#### zia for example.com no ssl

```bash
zia -domain example.com -port 80 -ssl=false -targets http://127.0.0.1:8020,http://10.0.0.10:8021
```

### zia for example.com without ssl single target

```bash
zia -domain example.com -port 8080 -ssl=false -targets http://127.0.0.1:8020
```

#### Logging

Zia logs access requests to a file in the `/var/log/zia/<domain>` directory, with the format `acces_<port>.log`.

#### Notes

- Zia requires a valid domain name to obtain a Let's Encrypt certificate. It cannot be an IP address.
- If the `-ssl` flag is not provided, Zia will run in HTTP mode without SSL/TLS.
- The `-timeout` option sets a read timeout for proxy connections. If not provided or set to 0, there is no timeout.

#### Running Zia as a systemd service

Zia can be run as a systemd service for easier management and automatic startup on boot. A sample systemd unit file is included in the `assets/zia.service` file.

To install and run Zia as a systemd service:

1. Modify and copy the `zia@test.com.service` <<<change domain name>>>, file to the appropriate systemd directory (e.g., `/etc/systemd/system/`).
2. Modify the `ExecStart` line in the `zia@test.com.service` file to point to the correct path of the `zia` binary and the parameters, to match your request.
3. Run the following command to enable and start the service:

```bash
sudo systemctl enable --now zia
```
