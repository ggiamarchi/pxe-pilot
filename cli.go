package main

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
)

func setupCLI() {
	app := cli.App("pxepilot", "PXE Pilot")

	app.Command("server", "Run PXE Pilot server", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			// TODO
		}
	})

	app.Command("config", "PXE configuration commands", func(cmd *cli.Cmd) {
		cmd.Command("list", "List available PXE configurations", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				os.Stdout.WriteString("config list : not implemented\n")
			}
		})
		cmd.Command("deploy", "Deploy a configuration for a machine", func(cmd *cli.Cmd) {

			cmd.Spec = "MACHINE CONFIG"

			var (
				machine = cmd.StringArg("MACHINE", "", "Machine for whom to deploy a configuration")
				config  = cmd.StringArg("CONFIG", "", "Configuration to deploy")
			)

			cmd.Action = func() {
				os.Stdout.WriteString(fmt.Sprintf("config deploy : not implemented\n - %s - %s", *machine, *config))
			}
		})
	})

	app.Command("machine", "Machine commands", func(cmd *cli.Cmd) {
		cmd.Command("list", "List machines", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				os.Stdout.WriteString("machine list : not implemented\n")
			}
		})
	})

	app.Run(os.Args)
}
