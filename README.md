# Zia

### A simple reverse proxy written in Go

version: <span style="color: orange">0.3.4</span>

Zia is a reverse proxy written in Go language. It was created to provide a simple and fast way to access old Docker containers. Meanwhile I realized I need also a simple reverse proxy for multiple http servers. So I decided to create it.
If you have a real domain - zia will get the ssl certificate for it from Let's Encrypt.

#### Implementation
- 0.3.5
  - changed middleware for allowed ip addresses to use memory cache
  - implemented CIDR ip address validation for allowed ip
  - config can now have a CIDR ip address in allowed ip list or a single ip address
- 0.3.4
  - implemented html error page to hide internal ip and real error message
  - simplified default page renderer
- 0.3.3
  - implemented config file
  - implemented acces list
  - default page to show for ip addresses that are not in acces list
  - html renderer for default page 
  - test target server with name and port shown - for testing proxy config
- 0.3.2 fix
  - ssl=false start without frontend ssl
  - stdout=true redirect stdout and stderr to console instead logfile
- 0.3.1 simplified the execution without configuration file just params.

#### Features

- Reverse proxy
- SSL/TLS frontend support with Let's Encrypt certificates
- Round-robin load balancing
- SSL/TLS targets with self-signed or CA signed
- https->https targets, https->http targets, http->http targets, (http|https)->mixed targets
- Configurable timeout for proxy connections
- Access logging on logfiles or stdout
- Access list of IP for targets. Html domain page will be shown instead of targets if client ip is not defined in acces list

#### Requirements

- Go installed
- IPv4 network

#### Installing Zia

1. Build the code with `go build zia.go` or `make zia`
2. Install `make install`
3. If below 1024 ports required, in Linux run `sudo setcap CAP_NET_BIND_SERVICE=+eip zia` to allow Zia to bind to low ports
4. For testing purpose there's a testTargetServer for testing zia with params

#### Parameters

- `-dev`: Developer mode - force logs to stdout, and disable html renderer cache
- `-stdout`: Show logs to stdout - good for debuging in real mode.
- `-ver`: Show app version and exists.
- `-configfile`: Configuration file

#### Run Zia

Run the binary with the desired options:

Zia will load balance requests between the defined targets using a round-robin algorithm.

##### start in developer mode - stdout logs and disable renderer cache so you can edit html file and see changes

```bash
zia -config config/test.json -dev
```

##### start in production mode but with stdout logs

```bash
zia -config config/test.json -stdout
```

#### show version number

```bash
zia -version
```

#### Config file structure
```json
{
  "port": 8080,
  "ssl": true,
  "cert": "cert/localdev.pem",
  "key":  "cert/localdev.key",
  "domain_name": "test.local",
  "no_access_page": "index.html",
  "targets": ["http://127.0.0.1:8081", "http://192.168.200.110:8082"],
  "allowed_ips": [
    { "ip": "192.168.200.110", "name": "local dev host" },
    { "ip": "10.0.10.0/24", "name": "vpn connection" }
  ]
}
```
- <span style="color: lightblue">**port**</span>: (<span style="color: orange;">optional</span>) - tcp port to listen for connections, default is 443
- <span style="color: lightblue">**ssl**</span>: (<span style="color:orange">optional</span>) - enable ssl
- <span style="color: lightblue">**cert**</span>: (<span style="color:orange">optional</span>) - certificate file to load
- <span style="color: lightblue">**key**</span>: (<span style="color:orange">optional</span>) - certificate key to load
- <span style="color: lightblue">**domain_name**</span>: <<span style="color:red">required</span>> - domain name to use for the proxy and Let's encrypt request for ssl

- <span style="color: lightblue">**targets**</span>: <<span style="color:red">required</span>> - list of target hosts to proxy to, minimum one

- <span style="color: lightblue">**no_acces_page**</span>: (<span style="color:orange">optional</span>) - html page to show if allowed_ips is set and the client IP is not in the allowed_ips list
- <span style="color: lightblue">**allowed_ips**</span>: (<span style="color:orange">optional</span>) - List of allowed IP that can acces targets, otherwise they will be rendered the no_access_page

#### <span style="color: darkred">Important!!!</span>

You have to set the ***NET_BIND_SERVICE*** capability on the binary to be able to bind to ports below 1024. 
```bash
sudo setcap CAP_NET_BIND_SERVICE=+eip zia
```

#### Logging

Zia logs access requests to a file in the `/var/log/zia/<domain>` directory, with the format `acces_<port>.log`.

#### Example

in 2 separate terminals run:

| **terminal 1**                            | **terminal 2**                             |
|:--------------------------------------|:---------------------------------------|
| test_target_server                    | test_target_server -name two -port 8082|

then set config <span style="color: lightblue">**targets**</span> one to **http://127.0.0.1:8082** and **http://192.168.0.100:8081**
[ *change **192.168.0.100** with your internal IP* ]

now you have 2 http targets, run zia to test the config.


#### Running Zia as a systemd service

Zia can be run as a systemd service for easier management and automatic startup on boot. A sample systemd unit file is included in the `assets/zia.service` file.

To install and run Zia as a systemd service:

1. Modify and copy the `zia.service` to `zia@yourdomain` file to the appropriate systemd directory (e.g., `/etc/systemd/system/`).
2. Modify the `ExecStart` line in the `zia@yourdomain` file to point to the correct path of the `zia` binary and the parameters, to match your request.
3. Run the following command to enable and start the service:

```bash
sudo systemctl enable --now zia@yourdomain
```
