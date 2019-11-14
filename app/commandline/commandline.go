package commandline

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type commandExec struct {
	executable    string
	arguments     []string
	commandOutput string
}

func New(executable string, arguments []string) commandExec {
	ce := commandExec{executable, arguments, ""}
	return ce
}

func (c commandExec) Run() {
	cmd := exec.Command(c.executable, c.arguments...)
	cmdStr := string(strings.Join(cmd.Args, " "))
	fmt.Printf("Executing [ %s ]\n", cmdStr)

	output, err := withCombinedOutput(cmd)
	if err != nil {
		log.Fatalf("command [ %s ] failed with %s\n", cmdStr, err)
	}

	c.saveOutput(output)
}

func (c commandExec) Output() string {
	return c.commandOutput
}

func (c *commandExec) saveOutput(output string) {
	c.commandOutput = output
}

func withCombinedOutput(cmd *exec.Cmd) (string, error) {
	out, err := cmd.CombinedOutput() //TODO: get stderr and stdout in separate strings
	outputStr := string(out)

	if err != nil {
		fmt.Printf("Output: \n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return "", err
	}

	fmt.Printf("Output: \n%s\n", outputStr)
	return outputStr, nil
}

func withSeparatedOutput(cmdStr string, cmd *exec.Cmd) error {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	errStr := string(stderr.Bytes())

	if err != nil {
		log.Fatalf("command [%s] failed with %s\n", cmdStr, err)
		fmt.Printf("Error output: %s\n", errStr)
	}

	outStr := string(stdout.Bytes())
	fmt.Printf("Output: \n%s\n%s\n", outStr, errStr)
	return nil
}
