# Zia

### just a simple reverse proxy

version: 0.0.4 beta

Zia is a reverse proxy written in go language. I needed a simple and fast reverse proxy to acces some old docker containers.

##### ToDo

- https mode with valid certificate - let's encrypt
- some statistics
- color mode cli
- limit connections

##### Features

- reverse proxy
- simple json config file
- https mode with self-signed certificate

###### requirements

- golang installed
- ipv4 network

#### Installing Zia

just build the code, ``` make zia ``` and run ` sudo make install `
rename config_sample.json to config.sample and edit the config

Config format:
work in progress ... to finish documentation
