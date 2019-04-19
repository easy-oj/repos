package cmd

import (
	"io"

	"github.com/easy-oj/common/settings"
)

func Tar(args ...string) *cmd {
	return New(settings.Path.Tar, args...)
}

func TarCreate(dir string, out io.Writer) *cmd {
	return Tar("-c", "-f", "-", ".").WorkDir(dir).Stdout(out)
}

func TarExtract(dir string, in io.Reader) *cmd {
	return Tar("-x", "-f", "-", "-C", dir).Stdin(in)
}
