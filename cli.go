package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ggiamarchi/pxe-pilot/api"
	"github.com/ggiamarchi/pxe-pilot/common/http"
	"github.com/ggiamarchi/pxe-pilot/logger"
	"github.com/ggiamarchi/pxe-pilot/model"

	cli "github.com/jawher/mow.cli"
	"github.com/olekukonko/tablewriter"
)

func setupCLI() {

	app := cli.App("pxe-pilot", "PXE Pilot")

	serverURL := app.StringOpt("s server", "http://localhost:3478", "Server URL for PXE Pilot client")
	debug := app.BoolOpt("d debug", false, "Show client logs on stdout")

	app.Command("server", "Run PXE Pilot server", func(cmd *cli.Cmd) {

		var configFile = cmd.StringOpt("c config", "/etc/pxe-pilot/pxe-pilot.yml", "PXE Pilot YAML configuration file")

		cmd.Action = func() {
			logger.Init(false)
			api.Run(*configFile)
		}
	})

	app.Command("bootloaders", "Bootloaders configuration commands", func(cmd *cli.Cmd) {
		cmd.Command("list", "List available bootloaders", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				logger.Init(!*debug)
				var bootloaders = &[]*model.Bootloader{}
				statusCode, err := http.Request("GET", *serverURL, "/v1/bootloaders", nil, bootloaders)
				if err != nil || statusCode != 200 {
					panic(err)
				}

				// Print data table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "File", "Config path"})
				for _, b := range *bootloaders {
					table.Append([]string{b.Name, b.File, b.ConfigPath})
				}
				table.Render()
			}
		})
	})

	app.Command("config", "PXE configuration commands", func(cmd *cli.Cmd) {
		cmd.Command("show", "Show PXE configurations", func(cmd *cli.Cmd) {

			cmd.Spec = "NAME"

			var (
				name = cmd.StringArg("NAME", "", "Configuration to show")
			)

			cmd.Action = func() {
				logger.Init(!*debug)

				var configuration = &model.ConfigurationDetails{}
				statusCode, err := http.Request("GET", *serverURL, fmt.Sprintf("/v1/configurations/%s", *name), nil, configuration)
				if err != nil || statusCode != 200 {
					panic(err)
				}

				fmt.Println(configuration.Content)
			}
		})
		cmd.Command("list", "List available PXE configurations", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				logger.Init(!*debug)
				var configurations = &[]*model.Configuration{}
				statusCode, err := http.Request("GET", *serverURL, "/v1/configurations", nil, configurations)
				if err != nil || statusCode != 200 {
					panic(err)
				}

				// Print data table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "Bootloader name", "Bootloader file", "Bootloader config path"})
				for _, c := range *configurations {
					table.Append([]string{c.Name, c.Bootloader.Name, c.Bootloader.File, c.Bootloader.ConfigPath})
				}
				table.Render()
			}
		})
		cmd.Command("deploy", "Deploy a configuration for a host", func(cmd *cli.Cmd) {

			cmd.Spec = "[-n] CONFIG HOSTNAMES..."

			var (
				now = cmd.BoolOpt("n now", false, "Trigger a server reboot when the configuration is set")

				config    = cmd.StringArg("CONFIG", "", "Configuration to deploy")
				hostnames = cmd.StringsArg("HOSTNAMES", []string{}, "Hosts for whom to deploy a configuration")
			)

			cmd.Action = func() {

				logger.Init(!*debug)

				hosts := make([]*model.HostQuery, len(*hostnames))

				for i, h := range *hostnames {
					hosts[i] = &model.HostQuery{
						Name:   h,
						Reboot: *now,
					}
				}

				hostsQuery := &model.HostsQuery{
					Hosts: hosts,
				}

				resp := &model.HostsResponse{}

				statusCode, err := http.Request("PUT", *serverURL, "/v1/configurations/"+*config+"/deploy", hostsQuery, resp)

				if err != nil || statusCode != 200 {
					cli.Exit(1)
				}

				// Print data table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetAutoWrapText(false)
				table.SetHeader([]string{"Name", "Configuration", "Rebooted"})

				for _, h := range resp.Hosts {
					table.Append([]string{h.Name, *config, h.Rebooted})
				}

				table.Render()
			}
		})
	})

	app.Command("host", "Host commands", func(cmd *cli.Cmd) {
		cmd.Command("list", "List hosts", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				logger.Init(!*debug)
				var hosts = &[]*model.Host{}
				statusCode, err := http.Request("GET", *serverURL, "/v1/hosts", nil, hosts)

				if err != nil {
					os.Stdout.WriteString("Error : " + err.Error())
				}

				if err != nil || statusCode != 200 {
					os.Stdout.WriteString("Error...")
					cli.Exit(1)
				}

				// Print data table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "Configuration", "MAC", "IPMI MAC", "IPMI HOST", "Power State"})
				table.SetAutoWrapText(false)

				for _, h := range *hosts {
					var configuration string
					if h.Configuration != nil {
						configuration = h.Configuration.Name
					}

					var ipmi *model.IPMI
					if h.IPMI != nil {
						ipmi = h.IPMI
					} else {
						ipmi = &model.IPMI{}
					}

					var macAddresses bytes.Buffer

					for i := 0; i < len(h.MACAddresses); i++ {
						if i != 0 {
							macAddresses.WriteString(" | ")
						}
						macAddresses.WriteString(h.MACAddresses[i])
					}

					table.Append([]string{h.Name, configuration, macAddresses.String(), ipmi.MACAddress, ipmi.Hostname, ipmi.Status})
				}
				table.Render()
			}
		})
		cmd.Command("reboot", "(re)boot a host", func(cmd *cli.Cmd) {
			cmd.Spec = "HOSTNAME"

			var (
				hostname = cmd.StringArg("HOSTNAME", "", "Host to reboot or reboot if powered off")
			)

			cmd.Action = func() {

				logger.Init(!*debug)

				statusCode, err := http.Request("PATCH", *serverURL, "/v1/hosts/"+*hostname+"/reboot", nil, nil)

				// Print data table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetAutoWrapText(false)
				table.SetHeader([]string{"Name", "Reboot"})

				if err != nil || statusCode != 204 {
					table.Append([]string{*hostname, "ERROR"})
					table.Render()
					cli.Exit(1)
				} else {
					table.Append([]string{*hostname, "OK"})
					table.Render()
				}
			}
		})
		cmd.Command("refresh", "Refresh hosts information", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				logger.Init(!*debug)
				statusCode, err := http.Request("PATCH", *serverURL, "/v1/refresh", nil, nil)

				// Print data table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Refresh"})
				table.SetAutoWrapText(false)
				if err != nil {
					table.Append([]string{"ERROR : " + err.Error()})
				}
				if err != nil || statusCode != 204 {
					table.Append([]string{"ERROR"})
					cli.Exit(1)
				}
				table.Append([]string{"OK"})
				table.Render()
			}
		})
	})

	err := app.Run(os.Args)
	if err != nil {
		logger.Error("%s", err)
	}
}
