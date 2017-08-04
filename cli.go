package main

import (
	"bytes"
	"fmt"
	"os"

	"dev.splitted-desktop.com/horizon/pxe-pilot/api"
	"dev.splitted-desktop.com/horizon/pxe-pilot/common/http"
	"dev.splitted-desktop.com/horizon/pxe-pilot/logger"
	"dev.splitted-desktop.com/horizon/pxe-pilot/model"

	cli "github.com/jawher/mow.cli"
	"github.com/olekukonko/tablewriter"
)

func setupCLI() {

	app := cli.App("pxepilot", "PXE Pilot")

	serverURL := app.StringOpt("s server", "http://localhost:3478", "Server URL for PXE Pilot client")
	debug := app.BoolOpt("d debug", false, "Show client logs on stdout")

	app.Command("server", "Run PXE Pilot server", func(cmd *cli.Cmd) {

		var configFile = cmd.StringOpt("c config", "/etc/pxepilot/pxepilot.yml", "PXE Pilot YAML configuration file")

		cmd.Action = func() {
			logger.Init(false)
			api.Run(*configFile)
		}
	})

	app.Command("config", "PXE configuration commands", func(cmd *cli.Cmd) {
		cmd.Command("list", "List available PXE configurations", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				logger.Init(!*debug)
				var configurations = &[]*model.Configuration{}
				statusCode, err := http.Request("GET", *serverURL, "/configurations", nil, configurations)
				if err != nil || statusCode != 200 {
					panic(err)
				}

				// Print data table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name"})
				for _, c := range *configurations {
					table.Append([]string{c.Name})
				}
				table.Render()
			}
		})
		cmd.Command("deploy", "Deploy a configuration for a host", func(cmd *cli.Cmd) {

			cmd.Spec = "HOST CONFIG"

			var (
				host   = cmd.StringArg("HOST", "", "Host for whom to deploy a configuration")
				config = cmd.StringArg("CONFIG", "", "Configuration to deploy")
			)

			cmd.Action = func() {
				logger.Init(!*debug)
				os.Stdout.WriteString(fmt.Sprintf("config deploy : not implemented\n - %s - %s", *host, *config))
			}
		})
	})

	app.Command("host", "Host commands", func(cmd *cli.Cmd) {
		cmd.Command("list", "List hosts", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				logger.Init(!*debug)
				var hosts = &[]*model.Host{}
				statusCode, err := http.Request("GET", *serverURL, "/hosts", nil, hosts)
				if err != nil || statusCode != 200 {
					panic(err)
				}

				// Print data table
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "Configuration", "MAC Addresses"})
				table.SetAutoWrapText(false)

				for _, h := range *hosts {
					var configuration string
					if h.Configuration != nil {
						configuration = h.Configuration.Name
					}

					var macAddresses bytes.Buffer

					for i := 0; i < len(h.MACAddresses); i++ {
						if i != 0 {
							macAddresses.WriteString(" | ")
						}
						macAddresses.WriteString(h.MACAddresses[i])
					}

					table.Append([]string{h.Name, configuration, macAddresses.String()})
				}
				table.Render()
			}
		})
	})

	app.Run(os.Args)
}
