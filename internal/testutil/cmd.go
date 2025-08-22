package testutil

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

// ExecuteCommand executes a cobra command with the given args
// and returns the output and error
func ExecuteCommand(t *testing.T, cmd *cobra.Command, args []string) (string, error) {
	t.Helper()
	
	// Create buffer to capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	
	// Execute command
	err := cmd.Execute()
	
	return buf.String(), err
}

// ExecuteCommandSeparateOutputs executes a cobra command with the given args
// and returns stdout, stderr and error separately
func ExecuteCommandSeparateOutputs(t *testing.T, cmd *cobra.Command, args []string) (stdout, stderr string, err error) {
	t.Helper()
	
	// Create buffers to capture output
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd.SetOut(outBuf)
	cmd.SetErr(errBuf)
	cmd.SetArgs(args)
	
	// Execute command
	err = cmd.Execute()
	
	return outBuf.String(), errBuf.String(), err
}