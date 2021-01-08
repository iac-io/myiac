package testutil

import (
	"strings"

	"github.com/iac-io/myiac/internal/commandline"
)

type fakeRunner struct {
	cmd      string
	args     []string
	CmdLines []string
	output   string
}

func (fk *fakeRunner) SetupWithoutOutput(cmd string, args []string) {
	fk.cmd = cmd
	fk.args = args
}

func (fk *fakeRunner) Run() commandline.CommandOutput {
	currentCmdLine := fk.cmd + " " + strings.Join(fk.args, " ")
	fk.CmdLines = append(fk.CmdLines, currentCmdLine)
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

func (fk fakeRunner) GetCmdLines() []string {
	return fk.CmdLines
}

func (fk *fakeRunner) SetupCmdLine(cmdLine string) {
	commandParts := strings.Split(cmdLine, " ")
	fk.cmd = commandParts[0]
	fk.args = commandParts[1:]
}

func FakeKubernetesRunner(output string) *fakeRunner {
	fakeRunner := new(fakeRunner)
	fakeRunner.output = output
	return fakeRunner
}
