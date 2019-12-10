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
	workingDir string
	SupressOutput bool
}

func NewEmpty() *commandExec {
	ce := &commandExec{"", make([]string,0), "", "", false}
	return ce
}

func New(executable string, arguments []string) *commandExec {
	ce := &commandExec{executable, arguments, "", "", false}
	return ce
}

func NewWithWorkingDir(executable string, arguments []string, workingDir string) *commandExec {
	ce := &commandExec{executable, arguments, "", workingDir, false}
	return ce
}

func (c *commandExec) Setup(executable string, arguments[]string) {
	c.executable = executable
	c.arguments = arguments 
}

func (c *commandExec) SetWorkingDir(workingDir string) {
	c.workingDir = workingDir
}

func (c commandExec) Run() commandExec {
	cmd := exec.Command(c.executable, c.arguments...)

	if (c.workingDir != "") {
		cmd.Dir = c.workingDir
		fmt.Printf("Working dir is: %s\n", c.workingDir)
	}

	cmdStr := string(strings.Join(cmd.Args, " "))
	fmt.Printf("Executing [ %s ]\n", cmdStr)

	output, err := withCombinedOutput(cmd, c.SupressOutput)
	if err != nil {
		log.Fatalf("command [ %s ] failed with %s\n", cmdStr, err)
	}

	c.saveOutput(output)
	return c
}

func (c commandExec) RunVoid() {
	c.Run()
}

func (c commandExec) Output() string {
	return c.commandOutput
}

func (c *commandExec) saveOutput(output string) {
	c.commandOutput = output
}

func withCombinedOutput(cmd *exec.Cmd, suppressOutput bool) (string, error) {
	out, err := cmd.CombinedOutput() //TODO: get stderr and stdout in separate strings
	outputStr := string(out)

	if (!suppressOutput) {
		fmt.Printf("Output: \n%s\n", outputStr)
	}
	
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return "", err
	}

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
