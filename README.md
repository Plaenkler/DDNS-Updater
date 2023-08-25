# 🌐 DDNS Updater

[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![Release](https://img.shields.io/badge/Calver-YY.WW.REVISION-22bfda.svg)](https://calver.org/)
[![Linters](https://github.com/plaenkler/ddns-updater/actions/workflows/linters.yml/badge.svg)](https://github.com/plaenkler/ddns-updater/actions/workflows/linters.yml)
[![Go Report](https://goreportcard.com/badge/github.com/plaenkler/ddns-updater)](https://goreportcard.com/report/github.com/plaenkler/ddns-updater)
[![CQL](https://github.com/plaenkler/ddns-updater/workflows/CodeQL/badge.svg)](https://github.com/plaenkler/ddns-updater)
[![Support me](https://img.shields.io/badge/Support%20me%20%E2%98%95-orange.svg)](https://www.buymeacoffee.com/Plaenkler)

DDNS Updater provides a straightforward way to update dynamic DNS records without messing with the command-line or a file.
The user-friendly interface allows for straightforward secure setup and management.

> **Note:** If your DynDNS provider is not listed open an issue and I will integrate it.

<table>
  <tr>
    <td><img src="https://user-images.githubusercontent.com/60503970/219900612-b4d7d3c4-7e0a-4dca-bc73-63c4822c5133.png" width="480"/></td>
    <td><img src="https://user-images.githubusercontent.com/60503970/219900611-dfaa9c4b-13ac-4fc4-b7ca-1cdae47961a9.png" width="480"/></td>
  </tr>
</table>

## 🎯 Noteworthy features

- [x] Simple & User friendly UI
- [x] Secure authentication with TOTP
- [x] Available as Docker Container
- [x] Scheduled update service
- [x] Supports multiple IP resolvers
- [ ] Deploy as Windows Service

## 🏷️ Supported providers

`Strato` `DDNSS` `Dynu` `Aliyun` `NoIP` `DD24` `INWX`

> **Note:** The crossed out providers will be implemented in future releases. In addition, the implementation of an individual update link with user-specific input and return values is planned.

## 📜 Installation guide

### 🐋 Deploy with Docker

It is recommended to use [Compose](https://docs.docker.com/compose/) as it is very convenient. The following examples show simple deployment options:

#### Bridge network

```yaml
---
version: '3.9'

services:
  ddns:
    image: plaenkler/ddns-updater:latest
    container_name: ddns
    restart: always
    networks:
      - web
    ports:
      - 80:80
    volumes:
      - ./ddns:/app/data
    environment:
      - DDNS_INTERVAL=600
      - DDNS_TOTP=false
      - DDNS_PORT=80

networks:
  web:
    external: false
```

#### Macvlan network

```yaml
---
version: '3.9'

services:
  ddns:
    image: plaenkler/ddns-updater:latest
    container_name: ddns
    restart: always
    networks:
      web:
        ipv4_address: 10.10.10.2
    volumes:
      - ./ddns:/app/data
    environment:
      - DDNS_INTERVAL=600
      - DDNS_TOTP=false
      - DDNS_PORT=80

networks:
  web:
    name: web
    driver: macvlan
    driver_opts:
      parent: eth0
    ipam:
      config:
        - subnet: "10.10.10.0/24"
          ip_range: "10.10.10.0/24"
          gateway: "10.10.10.1"
```

> **Note:** DDNS Updater can also be deployed behind a proxy like [Traefik](https://doc.traefik.io/traefik/) or [NGINX](https://www.nginx.com/).

### 📄 Build from source

From the root of the source tree, run:

```text
go build -o ddns-updater.exe cmd/main.go
```

> **Note:** Make sure that [CGO](https://gist.github.com/Plaenkler/0c319b89fbc884a928612b7fdef97fbd) is operational!


### ⚙️ Configuration

Depending on personal preferences, there are several ways to configure DDNS Updater. Users can select from three different methods:

**1. User Interface (Frontend)**

The easiest way to configure DDNS Updater is to use the integrated web interface. This can be accessed via the browser at `http://host-ip`.
Changes to the interval take effect immediately. The program must be restarted for the other settings to take effect.

**2. Configuration File**

A config.yaml file is provided to store all settings. In the absence of this file, the program generates one. Users have the option to directly modify settings within this file. It is important to note that changes made here will only take effect upon restarting the program. Default settings within the file are as follows:

```yaml
Interval: 600
TOTP: false
Port: 80
Resolver: ""
```

**3. Environment Variables**

An alternative to the configuration file are environment variables. During the program start these are read they are superordinate to the configuration file.

```text
DDNS_INTERVAL=600
DDNS_TOTP=false
DDNS_PORT=80
DDNS_RESOLVER=ipv4.example.com
```