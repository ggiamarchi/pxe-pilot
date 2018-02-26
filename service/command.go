package service

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ggiamarchi/pxe-pilot/logger"
)

func ExecCommand(command string, args ...interface{}) (string, string, error) {

	fmtCommand := fmt.Sprintf(command, args...)

	splitCommand := strings.Split(fmtCommand, " ")

	logger.Info("Executing command :: %s :: with args :: %v => %s", command, args, fmtCommand)

	cmdName := splitCommand[0]
	cmdArgs := splitCommand[1:len(splitCommand)]

	cmd := exec.Command(cmdName, cmdArgs...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}
