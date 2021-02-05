package testutil

import (
	"fmt"
	"strings"

	"github.com/iac-io/myiac/internal/commandline"
)

type fakeRunner struct {
	cmd             string
	args            []string
	CmdLines        []string
	output          string
	cmdLineToOutput map[string]string
	currentCmdLine  string
}

func (fk *fakeRunner) SetupWithoutOutput(cmd string, args []string) {
	fk.cmd = cmd
	fk.args = args
}

func (fk *fakeRunner) Run() commandline.CommandOutput {
	currentCmdLine := fk.cmd + " " + strings.Join(fk.args, " ")
	fk.CmdLines = append(fk.CmdLines, currentCmdLine)

	// every time we call Run(), check the cmdLine and return the
	// corresponding output
	if currentOutput, ok := fk.cmdLineToOutput[fk.currentCmdLine]; ok {
		fk.output = currentOutput
		return commandline.CommandOutput{Output: currentOutput}
	}

	// default output
	return commandline.CommandOutput{Output: fk.output}
}

func (fk fakeRunner) RunVoid() {
}

func (fk fakeRunner) Output() string {
	return fk.output
}

func (fk fakeRunner) Setup(cmd string, args []string) {
}

func (fk fakeRunner) IgnoreError(ignoreError bool) {
}

func (fk fakeRunner) SetSuppressOutput(suppressOutput bool) {
}

func (fk fakeRunner) GetCmdLines() []string {
	return fk.CmdLines
}

func (fk *fakeRunner) SetupCmdLine(cmdLine string) {
	commandParts := strings.Split(cmdLine, " ")
	fk.CmdLines = append(fk.CmdLines, cmdLine)
	fk.cmd = commandParts[0]
	fk.args = commandParts[1:]
	fk.currentCmdLine = cmdLine
}

func (fk *fakeRunner) OutputForCmdLine(cmdLine string) (string, error) {
	if currentOutput, ok := fk.cmdLineToOutput[cmdLine]; ok {
		return currentOutput, nil
	} else {
		return "", fmt.Errorf("error: could not find output for cmdLine %s", cmdLine)
	}
}

func (fk *fakeRunner) FakeCommand(cmdLine string, output string) {
	if fk.cmdLineToOutput == nil {
		fk.cmdLineToOutput = make(map[string]string)
	}
	fk.cmdLineToOutput[cmdLine] = output
}

func FakeCommandRunner(output string) *fakeRunner {
	fakeRunner := new(fakeRunner)
	fakeRunner.output = output
	return fakeRunner
}
