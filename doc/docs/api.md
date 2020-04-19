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
`reboot`       | bool     | No       | Whether the host should be rebooted automatically or not


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
