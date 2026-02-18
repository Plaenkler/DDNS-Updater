# 🌐 DDNS Updater

[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![Release](https://img.shields.io/badge/Calver-YY.WW.REVISION-22bfda.svg)](https://calver.org/)
[![Linters](https://github.com/plaenkler/ddns-updater/actions/workflows/linters.yml/badge.svg)](https://github.com/plaenkler/ddns-updater/actions/workflows/linters.yml)
[![Go Report](https://goreportcard.com/badge/github.com/plaenkler/ddns-updater)](https://goreportcard.com/report/github.com/plaenkler/ddns-updater)
[![CQL](https://github.com/plaenkler/ddns-updater/workflows/CodeQL/badge.svg)](https://github.com/plaenkler/ddns-updater)
[![Support me](https://img.shields.io/badge/Support%20me%20%E2%98%95-orange.svg)](https://www.buymeacoffee.com/Plaenkler)

DDNS Updater provides a straightforward way to update dynamic DNS records without messing with the command-line or a file.
The user-friendly interface allows for straightforward secure setup and management.

<table>
  <tr>
    <td><img src="https://user-images.githubusercontent.com/60503970/219900612-b4d7d3c4-7e0a-4dca-bc73-63c4822c5133.png" width="480" alt="DDNS Updater Dashboard"/></td>
    <td><img src="https://user-images.githubusercontent.com/60503970/219900611-dfaa9c4b-13ac-4fc4-b7ca-1cdae47961a9.png" width="480" alt="DDNS Updater add Job"/></td>
  </tr>
</table>

## 🎯 Noteworthy features

- [x] Simple & User friendly UI
- [x] Secure authentication with TOTP
- [x] Encryption of sensitive data
- [x] Scheduled update service
- [x] Supports multiple IP resolvers
- [x] Deploy as Windows Service
- [x] Available as Docker Container
- [x] Custom update URL with check

## 🏷️ Supported providers

`Strato` `DDNSS` `Dynu` `Aliyun` `NoIP` `DD24` `INWX` `Infomaniak` `Hetzner` `IONOS`

> **Note:** If your DynDNS provider is not listed open an issue and I will integrate it.

### Custom provider

Select the provider **Custom** from the provider list.
Then enter the user-defined URL in the **URL** field.
Use the placeholder `<ipv4>` at the point where the IPv4 address is to be inserted.
In the **“Check”** field, enter a string that will be used for checking to ensure that the update server's response is successful.

## 📜 Installation guide

### 🐋 Deploy with Docker

It is recommended to use [Compose](https://docs.docker.com/compose/) as it is very convenient. The following examples show simple deployment options:

#### Bridge network

```yaml
---

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

DDNS Updater requires the following software to be installed on your system:

- [Go](https://golang.org/dl/)
- [CGO](https://gist.github.com/Plaenkler/0c319b89fbc884a928612b7fdef97fbd)

**Windows Service**

After cloning the repository, navigate to the root of the source tree and run the following command to build the executable file:

```text
go build -o ddns-updater.exe cmd/svc/main.go
```

The program can then be stored in any file path and set up as a service using the CMD:

```text
sc create DDNS-Updater binPath=[INSTALL_DIR]\ddns-updater.exe type=share start=auto DisplayName=DDNS-Updater
```

**Standalone application**

To build the standalone application, run the following command:

```text
go build -o ddns-updater cmd/def/main.go
```

### ⚙️ Configuration

Depending on personal preferences, there are several ways to configure DDNS Updater. Users can select from three different methods:

**1. User Interface (Frontend)**

The easiest way to configure DDNS Updater is to use the integrated web interface. This can be accessed via the browser at `http://host-ip`.
Changes to the interval take effect immediately. The program must be restarted for the other settings to take effect.

**2. Configuration File**

A config.yaml file is provided to store all settings. In the absence of this file, the program generates one. Users have the option to directly modify settings within this file. It is important to note that changes made here will only take effect upon restarting the program. Default settings within the file are as follows:

```yaml
# How often the IP address is checked in seconds
Interval: 600
# Enable TOTP authentication
TOTP: false
# Port for the web interface
Port: 80
# Custom IP resolver returns IPv4 address in plain text
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
