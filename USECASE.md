# PXE Pilot use case

Here comes a real life example where PXE Pilot can help. This is actually the use case
that leads to the creation of PXE Pilot.

## Initial situation

On a sub-network we have a PXE server to boot hosts over the network. It is composed of

* A DHCP server
* A TFTP server
* A HTTP (to serve additionnal files after PXE boot)

We have 3 differents boot configurations:

* __Ubuntu 16.06 Installer__ - This runs the Ubuntu installer in an automated way using a preseed configuration file. When the installation is done the host reboots
* __GRML 2017.05 Live__ - This runs a live distribution with an in-memory filesystem. This boot configuration is useful to run a server in a rescue mode to debug potential issues
* __Local__ - A PXE configuration that leads the host to boot on the local drive. This is useful when we want to force a host to boot on the local drive even if the network device has  a higher priority in the BIOS boot sequence

Here is the content of the TFTP root directory:

```
|-- pxelinux.cfg
|   |-- conf
|   |   |-- ubuntu-16.04
|   |   |-- grml-2017.05
|   |   `-- local
|-- pxelinux.0
|-- ldlinux.c32
|-- grml
|   `-- 2017.05
|       |-- initrd.img
|       `-- vmlinuz
`-- ubuntu
    `-- 1604
        |-- initrd.gz
        `-- linux
```

And here the contemt of the HTTP root directory:

```
.
|-- grml64-full_2017.05.iso    # Filesystem for the GRML disribution
|-- ubuntu-16.04.seed          # Preseed configuration for Ubuntu installer
`-- ubuntu-16.04.sh            # Script executed at the very end of the preseed
```

At this point, if we wish a specific host to boot with a specific configuration, we need
to create a configuration file named after the host's MAC address.

Let's say we wish the MAC address `23:34:a3:e9:09:cc` to boot with `grml-2017.05`
configuration. We can basically do

```
$ ln -s pxelinux.cfg/conf/grml-2017.05 pxelinux.cfg/01-23-34-a3-e9-09-cc
```

## Problem we wish to solve

First, let's consider the case where we wish to install Ubuntu 16.04 for a specific host.

__To achieve that, what can we do ?__

1. Deploy `ubuntu-16.04` configuration for this host
2. Start (or restart) the host

__Then what's happen ?__

1. The server PXE boot the Ubuntu installer
2. The installer loads the `ubuntu-16.04.seed` file containing preseed instructions
3. The installation, at the end, loads and executes the script `ubuntu-16.04.sh`
4. The server reboots
5. The server PXE boot the Ubuntu installer
6. ...

__Do you see the issue ?__

As long as the `ubuntu-16.04` configuration is present for this host, the host boots on
the installer and then reboot in an infinite loop.

## Looking for a solution

Once the operating system installation is complete, we would like to replace the host's PXE
configuration by the `local` configuration.

This way, after the installation completes, the server reboots on the freshly installed O/S.

If we are able to do that from the script `ubuntu-16.04.sh`, it's victory. That said, how
to change the configuration on the PXE server remotely ?

If there is an API exposing an opreration to switch from a PXE configuration to another
one, we are done.

## Finally

__Here comes PXE Pilot :)__

From the script `ubuntu-16.04.sh`, we just need to send a `PUT` request

```
curl -i -X PUT http://pxe-server:3478/configurations/local/deploy -d '
{
    "hosts": [
        {
            "macAddress": "23:34:a3:e9:09:cc"
        }
    ]
}'
```
