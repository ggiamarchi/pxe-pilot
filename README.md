# PXE Pilot

[![Build Status](https://api.travis-ci.org/ggiamarchi/pxe-pilot.png?branch=master)](https://travis-ci.org/ggiamarchi/pxe-pilot)

PXE Pilot helps you to manage PXE configurations for your hosts through a simple API. His role is very simple, it allows you to:

- Know which configuration is associated with each host
- To switch from a configuration to another one for a specific host

The REST API can be called either directly or using the CLI.

# Why this software ?

Some other software like Foreman or Cobbler can be used to manage PXE configurations in a more
sofisticated way. For small deployments where only very basics features are needed, installing
such a solution can appear overkill. In this situation, PXE Pilot could be what you need.

See also [the use case](USECASE.md) that leads me to create this project.

# How does it work ?

PXE pilot knows your hosts (with their MAC addresses) and your PXE configurations (see the
"configuration" section)

It's important to understand that PXE Pilot does not help you to write PXE configurations. PXE Pilot
basically manages symlinks at the filesystem level to ensure each host uses the right configuration.

When a host have several MAC addresses, the first address in the list points to the desired configuration
and all the others point to the first one. This is to ensure the host will boot on the right configuration
whatever the network interface is used to boot.

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

# API Documentation

## Read configurations

```
GET /v1/configurations
```

###### Response

```json
{
    "configurations": [
        {
            "name": "ubuntu-16.04"
        },
        {
            "name": "local"
        },
        {
            "name": "grml-2017.05"
        }
    ]
}
```

###### Response codes

Code   | Name        | Description
-------|-------------|---------------------------------------------------
`20O`  | `Ok`        | Server configurations have been retrieved

## Show configurations

```
GET /v1/configurations/<name>
```

###### Response

```json
{
  "name": "local",
  "content": "default local\n\nlabel local\n    localboot 0\n"
}
```

###### Response codes

Code   | Name        | Description
-------|-------------|---------------------------------------------------
`20O`  | `Ok`        | Configuration detail have been retrieved

## Read hosts

```
GET /v1/hosts
```

###### Response

```json
[
    {
        "name": "h1",
        "macAddresses": [
            "00:00:00:00:00:01"
        ],
        "configuration": null
    },
    {
        "name": "h2",
        "macAddresses": [
            "00:00:00:00:00:02"
        ],
        "configuration": {
            "name": "ubuntu-16.04"
        }
    },
    {
        "name": "h3",
        "macAddresses": [
            "00:00:00:00:00:03",
            "00:00:00:00:00:33"
        ],
        "configuration": {
            "name": "local"
        }
    }
]
```

###### Response codes

Code   | Name          | Description
-------|---------------|---------------------------------------------------
`20O`  | `Ok`          | Host list had been retrieved


## Reboot a host

```
PATCH /v1/hosts/<name>/reboot
```

###### Response codes

Code   | Name          | Description
-------|---------------|---------------------------------------------------
`204`  | `No Content`  | Host had been successfully rebooted
`404`  | `Not Found`   | Host does not exist
`409`  | `Conflict`    | Reboot did not succeed for any reason


## Deploy a configuration for host(s)

```
PUT /v1/configurations/<configuration_name>/deploy
```

###### Body

```json
{
    "hosts": [
        {
            "name": "h1"
        },
        {
            "name": "h2"
        },
        {
            "mac_address": "83:06:0a:00:cf:03"
        }
    ]
}
```

###### Error response (code 4xx)

```json
{
    "message": "Configuration not found"
}
```

###### Parameters

Name           | In    | Type     | Required | Description
---------------|-------|----------|----------|-----------------------------------
`hosts`        | body  | Host[]   | Yes      | Hosts for whom to deploy configuration

###### Host (object)

Attribute      | Type     | Required | Description
---------------|----------|----------|---------------------------------------------
`name`         | string   | No       | Host name
`mac_address`  | string   | No       | Host MAC address


###### Response codes

Code   | Name          | Description
-------|---------------|---------------------------------------------------
`204`  | `No Content`  | Configurations had been deployed
`404`  | `Not found`   | Either the configuation or a host is not found
`400`  | `Bad request` | Malformed body


## Refresh hosts information

This API populate the ARP table for all subnets in the PXE Pilot configuration

```
PATCH /v1/refresh
```

###### Response codes

Code   | Name          | Description
-------|---------------|---------------------------------------------------
`204`  | `No Content`  | Refresh operation completed without any issue


# License

Everything in this repository is published under the MIT license.
