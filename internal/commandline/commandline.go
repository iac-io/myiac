package commandline

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// CommandRunner Implicit interface for commandline package, need access to those methods here
type CommandRunner interface {
	RunVoid()
	Output() string
	Setup(cmd string, args []string)
	SetupWithoutOutput(cmd string, args []string)
	IgnoreError(ignoreError bool)
	Run() CommandOutput
}

type commandExec struct {
	executable       string
	arguments        []string
	commandOutput    string
	workingDir       string
	IsSuppressOutput bool
	ignoreError      bool
}

type CommandOutput struct {
	Output string
}

func NewEmpty() *commandExec {
	ce := &commandExec{"", make([]string, 0), "", "",
		false, false}
	return ce
}

func NewCommandLine(commandLine string) *commandExec {
	commandParts := strings.Split(commandLine, " ")
	executable := commandParts[0]
	args := commandParts[1:]
	ce := &commandExec{executable, args, "", "",
		false, false}
	return ce
}

func New(executable string, arguments []string) *commandExec {
	ce := &commandExec{executable, arguments, "", "",
		false, false}
	return ce
}

func NewWithWorkingDir(executable string, arguments []string, workingDir string) *commandExec {
	ce := &commandExec{executable, arguments, "", workingDir,
		false, false}
	return ce
}

func (c *commandExec) Setup(executable string, arguments []string) {
	c.executable = executable
	c.arguments = arguments
}

func (c *commandExec) SetupWithoutOutput(executable string, arguments []string)  {
	c.executable = executable
	c.arguments = arguments
	c.IsSuppressOutput = true
}

func (c *commandExec) SetWorkingDir(workingDir string) {
	c.workingDir = workingDir
}

func (c *commandExec) Run() CommandOutput {
	cmd := exec.Command(c.executable, c.arguments...)

	if c.workingDir != "" {
		cmd.Dir = c.workingDir
		fmt.Printf("Working dir is: %s\n", c.workingDir)
	}

	cmdStr := string(strings.Join(cmd.Args, " "))
	fmt.Printf("Executing [ %s ]\n", cmdStr)

	// output, err := withCombinedOutput(cmd, c.IsSuppressOutput)
	output, err := withProgress(cmd, c.IsSuppressOutput, c.ignoreError)

	if err != nil && !c.ignoreError {
		log.Fatalf("command [ %s ] failed with %s\n", cmdStr, err)
	}

	if c.ignoreError {
		log.Printf("Ignoring error for command [ %s ] with %v\n", cmdStr, err)
	} else {
		c.saveOutput(output)
	}

	outputResult := CommandOutput{Output:c.commandOutput}

	return outputResult
}

func (c *commandExec) RunVoid() {
	// Important: for this delegation to work properly and save the output, we need to
	// pass in a pointer, which is what ultimately gets modified in the 'saveOutput' method
	c.Run()
}

func (c commandExec) Output() string {
	return c.commandOutput
}

func (c *commandExec) IgnoreError(ignoreError bool) {
	c.ignoreError = ignoreError
}

func (c *commandExec) saveOutput(output string) {
	c.commandOutput = output
}

func (c *commandExec) SuppressOutput(suppressOutput bool) {
	c.IsSuppressOutput = suppressOutput
}

func withCombinedOutput(cmd *exec.Cmd, suppressOutput bool) (string, error) {
	out, err := cmd.CombinedOutput() //TODO: get stderr and stdout in separate strings
	outputStr := string(out)

	if !suppressOutput {
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

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

func withProgress(cmd *exec.Cmd, suppressOutput bool, ignoreError bool) (string, error) {

	var stdout, stderr []byte
	var errStdout, errStderr error

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	// cmd.Wait() should be called only after we finish reading
	// from stdoutIn and stderrIn.
	// wg ensures that we finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
		wg.Done()
	}()

	stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)

	wg.Wait()

	err = cmd.Wait()

	if err != nil && !ignoreError {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	if err != nil && ignoreError {
		log.Printf("cmd.Run() error ignored due to flag 'ignoreError' set %s\n", err)
	}

	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}

	outputStr, errorStr := string(stdout), string(stderr)

	// Not sure if this is the way, but there are valid data on stdout and stderr
	combinedOutputStr := outputStr + "\n" + errorStr
	if !suppressOutput {
		fmt.Printf("\nOutput:\n%s\n", combinedOutputStr)
	}

	return combinedOutputStr, nil
}
