//go:build e2e
// +build e2e

package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.Println("start")

	exit := m.Run()

	log.Println("stop")

	os.Exit(exit)
}

func TestA(t *testing.T) {
	t.Log("a")
}

func TestB(t *testing.T) {
	t.Log("b")
}

func TestBar(t *testing.T) {
	type testcase struct {
		args   []string
		stdin  io.Reader
		stdout io.Writer
		stderr io.Writer
	}

	tt := map[string]testcase{
		"foo": {
			args:   []string{"bar"},
			stdin:  os.Stdin,
			stdout: os.Stdout,
			stderr: os.Stderr,
		},
	}

	for tn, tc := range tt {
		t.Run(tn, func(t *testing.T) {
			process := &Process{
				Command: "../../bin/template",
				Args:    tc.args,
			}

			err := process.Run()
			require.NoError(t, err)
		})
	}

}

type Process struct {
	Command string
	Stdout  io.Writer
	Stdin   io.Reader
	Stderr  io.Writer
	Args    []string
	Env     []string
}

func (m *Process) Run() error {
	cmd := exec.Command(m.Command, m.Args...)

	cmd.Args = m.Args
	cmd.Stdin = m.Stdin
	cmd.Stdout = m.Stdout
	cmd.Stderr = m.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
