# Prerequisites

Using PXE Pilot suppose you already have a PXE server (DHCP + TFTP) up and running. PXE Pilot is DHCP
and TFTP agnostic. It only work at the filesystem level into the TFTP root directory.

# Running PXE Pilot

## Configuration

PXE Pilot needs to know three things:

- The host list to manage
- The absolute path to the TFTP root
- The directory containing PXE configurations

All those information are described in the YAML file `/etc/pxe-pilot/pxe-pilot.yml`.

Optionnaly, IPMI MAC address (or IP address) and credentials can be specified. When IPMI is available,
PXE Pilot client shows power state for each host.

__Example:__

```yaml
---

hosts:
  - name: h1
    mac_addresses: ["00:00:00:00:00:01"]
  - name: h2
    mac_addresses: ["00:00:00:00:00:02"]
    ipmi:
      mac_address: "00:00:00:00:00:a2"
      username: "user"
      password: "pass"
      interface: "lanplus"
      subnets: "10.0.0.0/24"
  - name: h3
    mac_addresses: ["00:00:00:00:00:03", "00:00:00:00:00:33"]

tftp:
  root: "/var/tftp"

configuration:
  directory: /var/tftp/pxelinux.cfg/conf

server:
  port: 3478
```

## Running PXE Pilot server

Basically run

```
$ pxe-pilot server
```

## Querying PXE Pilot using the CLI

```
$ pxe-pilot --help

Usage: pxe-pilot [OPTIONS] COMMAND [arg...]

PXE Pilot

Options:
  -s, --server="http://localhost:3478"   Server URL for PXE Pilot client
  -d, --debug=false                      Show client logs on stdout

Commands:
  server       Run PXE Pilot server
  config       PXE configuration commands
  host         Host commands

Run 'pxe-pilot COMMAND --help' for more information on a command.
```

Th following examples assume PXE Pilot server is listening on `localhost:3478`. If not,
use the `--server` option to address your PXE Pilot server.

### List available configurations

```
$ pxe-pilot config list

+--------------+
|     NAME     |
+--------------+
| local        |
| ubuntu-14.04 |
| ubuntu-16.04 |
+--------------+
```

### List hosts

```
$ pxe-pilot host list

+------+---------------+---------------------------------------+-------------------+-----------+-------------+
| NAME | CONFIGURATION |             MAC ADDRESSES             |    IPMI MAC       | IPMI HOST | POWER STATE |
+------+---------------+---------------------------------------+-------------------+-----------+-------------+
| h1   | local         | 00:00:00:00:00:01                     |                   |           |             |
| h2   |               | 00:00:00:00:00:02                     | 00:00:00:00:00:a2 | 1.2.3.4   | On          |
| h3   | local         | 00:00:00:00:00:03 | 00:00:00:00:00:33 |                   |           |             |
+------+---------------+---------------------------------------+-------------------+-----------+-------------+
```

### Deploy configuration for host(s)

Deploy `ubuntu-16.04` configuration for hosts `h2`and `h3`.

```
$ pxe-pilot config deploy ubuntu-16.04 h2 h3

+------+---------------+
| NAME | CONFIGURATION |
+------+---------------+
| h2   | ubuntu-16.04  |
| h3   | ubuntu-16.04  |
+------+---------------+
```
