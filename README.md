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

# License

Everything in this repository is published under the MIT license.
