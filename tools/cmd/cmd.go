package cmd

import (
	"io"
	"os/exec"

	"github.com/easy-oj/common/logs"
)

type cmd struct {
	c *exec.Cmd
}

func New(name string, args ...string) *cmd {
	return &cmd{
		c: exec.Command(name, args...),
	}
}

func (c *cmd) WorkDir(dir string) *cmd {
	c.c.Dir = dir
	return c
}

func (c *cmd) Stdin(stdin io.Reader) *cmd {
	c.c.Stdin = stdin
	return c
}

func (c *cmd) Stdout(stdout io.Writer) *cmd {
	c.c.Stdout = stdout
	return c
}

func (c *cmd) Stderr(stderr io.Writer) *cmd {
	c.c.Stderr = stderr
	return c
}

func (c *cmd) Run() error {
	return c.c.Run()
}

func (c *cmd) JustRun() {
	if err := c.Run(); err != nil {
		logs.Warn("[CMD] run cmd with ignored error: %s", err.Error())
	}
}
